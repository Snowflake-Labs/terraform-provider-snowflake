package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceMonitorSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Identifier for the resource monitor; must be unique for your account.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	// TODO: This can be set to empty
	"notify_users": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Specifies the list of users (their identifiers) to receive email notifications on resource monitors.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"credit_quota": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "The number of credits allocated to the resource monitor per frequency interval. When total usage for all warehouses assigned to the monitor reaches this number for the current frequency interval, the resource monitor is considered to be at 100% of quota.",
	},
	"frequency": {
		Type:     schema.TypeString,
		Optional: true,
		// TODO: No default ? By default it's MONTHLY
		RequiredWith:     []string{"start_timestamp"},
		ValidateDiagFunc: sdkValidation(sdk.ToResourceMonitorFrequency),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToResourceMonitorFrequency), IgnoreChangeToCurrentSnowflakeValueInShow("frequency")),
		Description:      "The frequency interval at which the credit usage resets to 0. If you set a `frequency` for a resource monitor, you must also set `start_timestamp`. If you specify `NEVER` for the frequency, the credit usage for the warehouse does not reset.",
	},
	"start_timestamp": {
		// TODO: Skip checking if start_timestamp == IMMEDIATELY
		Type:         schema.TypeString,
		Optional:     true,
		RequiredWith: []string{"frequency"},
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return strings.ToUpper(oldValue) == "IMMEDIATELY"
		},
		// TODO: Not detecting external changes, because it's hard
		Description: "The date and time when the resource monitor starts monitoring credit usage for the assigned warehouses. If you set a `start_timestamp` for a resource monitor, you must also set `frequency`.",
	},
	"end_timestamp": {
		Type:     schema.TypeString,
		Optional: true,
		// TODO: Not detecting external changes, because it's hard
		Description: "The date and time when the resource monitor suspends the assigned warehouses.",
	},
	"trigger": {
		Type:     schema.TypeSet,
		Optional: true,
		// TODO: Throw error on 0 triggers alter
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"threshold": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"on_threshold_reached": {
					Type:             schema.TypeString,
					Required:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToResourceMonitorTriggerAction),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToResourceMonitorTriggerAction),
				},
			},
		},
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW RESOURCE MONITORS` for the given resource monitor.",
		Elem: &schema.Resource{
			Schema: schemas.ShowResourceMonitorSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func ResourceMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResourceMonitor,
		ReadContext:   ReadResourceMonitor,
		UpdateContext: UpdateResourceMonitor,
		DeleteContext: DeleteResourceMonitor,

		Schema: resourceMonitorSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "notify_users", "credit_quota", "frequency", "start_timestamp", "end_timestamp", "trigger"),
		),
	}
}

func CreateResourceMonitor(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	opts := new(sdk.CreateResourceMonitorOptions)
	with := new(sdk.ResourceMonitorWith)

	if v := d.Get("credit_quota").(int); v != IntDefault {
		with.CreditQuota = sdk.Pointer(v)
	}

	if v, ok := d.GetOk("notify_users"); ok {
		userIds := expandStringList(v.(*schema.Set).List())
		users := make([]sdk.NotifiedUser, len(userIds))
		for i, userId := range userIds {
			users[i] = sdk.NotifiedUser{
				Name: sdk.NewAccountObjectIdentifier(userId),
			}
		}
		opts.With.NotifyUsers = &sdk.NotifyUsers{Users: users}
	}

	if v, ok := d.GetOk("frequency"); ok {
		frequency, err := sdk.ToResourceMonitorFrequency(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.With.Frequency = frequency
	}

	if v, ok := d.GetOk("start_timestamp"); ok {
		opts.With.StartTimestamp = sdk.Pointer(v.(string))
	}

	if v, ok := d.GetOk("end_timestamp"); ok {
		opts.With.EndTimestamp = sdk.Pointer(v.(string))
	}

	if v, ok := d.GetOk("trigger"); ok {
		triggerDefinitions, err := extractTriggerDefinitions(v.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		opts.With.Triggers = triggerDefinitions
	}

	if !reflect.DeepEqual(*with, sdk.ResourceMonitorWith{}) {
		opts.With = with
	}

	err := client.ResourceMonitors.Create(ctx, id, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadResourceMonitor(ctx, d, meta)
}

func ReadResourceMonitor(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query resource monitor. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Resource Monitor: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set("notify_users", resourceMonitor.NotifyUsers); err != nil {
		return diag.FromErr(err)
	}

	if resourceMonitor.CreditQuota != nil {
		if err := d.Set("credit_quota", *resourceMonitor.CreditQuota); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("credit_quota", IntDefault); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("frequency", resourceMonitor.Frequency); err != nil {
		return diag.FromErr(err)
	}

	// TODO: Do i read timestamps ???

	var triggers []any

	if len(resourceMonitor.NotifyAt) > 0 {
		for _, notifyAt := range resourceMonitor.NotifyAt {
			triggers = append(triggers, map[string]any{
				"threshold":            notifyAt,
				"on_threshold_reached": sdk.TriggerActionNotify,
			})
		}
	}

	if resourceMonitor.SuspendAt != nil {
		triggers = append(triggers, map[string]any{
			"threshold":            resourceMonitor.SuspendAt,
			"on_threshold_reached": sdk.TriggerActionSuspend,
		})
	}

	if resourceMonitor.SuspendImmediateAt != nil {
		triggers = append(triggers, map[string]any{
			"threshold":            resourceMonitor.SuspendImmediateAt,
			"on_threshold_reached": sdk.TriggerActionSuspendImmediate,
		})
	}

	if err := d.Set("trigger", triggers); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.ResourceMonitorToSchema(resourceMonitor)}); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateResourceMonitor(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var runSetStatement bool
	opts := sdk.AlterResourceMonitorOptions{}
	set := sdk.ResourceMonitorSet{}

	if d.HasChange("credit_quota") {
		runSetStatement = true
		if v := d.Get("credit_quota").(int); v != IntDefault {
			set.CreditQuota = sdk.Pointer(d.Get("credit_quota").(int))
		} else {
			// TODO: Set to null
			//set.CreditQuota = sdk.Pointer(d.Get("credit_quota").(int))
		}
	}

	if d.HasChange("frequency") || d.HasChange("start_timestamp") {
		runSetStatement = true
		frequency, err := sdk.ToResourceMonitorFrequency(d.Get("frequency").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		set.Frequency = frequency
		set.StartTimestamp = sdk.Pointer(d.Get("start_timestamp").(string))
	}

	if d.HasChange("end_timestamp") {
		runSetStatement = true
		if v := d.Get("end_timestamp").(string); v != "" {
			set.EndTimestamp = sdk.Pointer(d.Get("end_timestamp").(string))
		} else {
			// TODO: Set to null
			//set.EndTimestamp = sdk.Pointer(d.Get("end_timestamp").(string))
		}
	}

	if d.HasChange("notify_users") {
		runSetStatement = true
		userIds := expandStringList(d.Get("notify_users").(*schema.Set).List())
		users := make([]sdk.NotifiedUser, len(userIds))
		for i, userId := range userIds {
			users[i] = sdk.NotifiedUser{
				Name: sdk.NewAccountObjectIdentifier(userId),
			}
		}
		set.NotifyUsers = &sdk.NotifyUsers{
			Users: users,
		}
	}

	if d.HasChange("trigger") {
		v := d.Get("trigger").(*schema.Set).List()
		if len(v) > 0 {
			triggerDefinitions, err := extractTriggerDefinitions(v)
			if err != nil {
				return diag.FromErr(err)
			}
			opts.Triggers = triggerDefinitions
		} else {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to update resource monitor.",
					Detail:   "Due to Snowflake limitations triggers cannot be completely removed form resource monitor after having at least 1 trigger. The only way it to re-create resource monitor without any triggers specified.",
				},
			}
		}
	}

	if runSetStatement {
		if set != (sdk.ResourceMonitorSet{}) {
			opts.Set = &set
		}
		if err := client.ResourceMonitors.Alter(ctx, id, &opts); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadResourceMonitor(ctx, d, meta)
}

func DeleteResourceMonitor(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.ResourceMonitors.Drop(ctx, id, &sdk.DropResourceMonitorOptions{IfExists: sdk.Bool(true)})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func extractTriggerDefinitions(triggers []any) ([]sdk.TriggerDefinition, error) {
	triggerDefinitions := make([]sdk.TriggerDefinition, len(triggers))
	for i, trigger := range triggers {
		triggerMap := trigger.(map[string]any)
		threshold := triggerMap["threshold"].(int)
		triggerAction, err := sdk.ToResourceMonitorTriggerAction(triggerMap["on_threshold_reached"].(string))
		if err != nil {
			return nil, err
		}
		triggerDefinitions[i] = sdk.TriggerDefinition{
			Threshold:     threshold,
			TriggerAction: *triggerAction,
		}
	}
	return triggerDefinitions, nil
}
