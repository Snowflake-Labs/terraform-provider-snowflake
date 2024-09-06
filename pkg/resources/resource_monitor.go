package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"reflect"
)

var resourceMonitorSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Identifier for the resource monitor; must be unique for your account.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
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
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("credit_quota"),
		Description:      "The number of credits allocated to the resource monitor per frequency interval. When total usage for all warehouses assigned to the monitor reaches this number for the current frequency interval, the resource monitor is considered to be at 100% of quota.",
	},
	// TODO: Describe that default it's MONTHLY, but after unsetting this field it will be the thing was previously set in the configuration; because no unset is available
	"frequency": {
		Type:             schema.TypeString,
		Optional:         true,
		RequiredWith:     []string{"start_timestamp"},
		ValidateDiagFunc: sdkValidation(sdk.ToResourceMonitorFrequency),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToResourceMonitorFrequency), IgnoreChangeToCurrentSnowflakeValueInShow("frequency")),
		Description:      "The frequency interval at which the credit usage resets to 0. If you set a `frequency` for a resource monitor, you must also set `start_timestamp`. If you specify `NEVER` for the frequency, the credit usage for the warehouse does not reset.",
	},
	// TODO: Describe that default it's MONTHLY, but after unsetting this field it will be the thing was previously set in the configuration; because no unset is available
	// TODO: Describe that it's advised for now to specify full dates of format 2024-10-04 00:00 otherwise diffs may occur
	"start_timestamp": {
		Type:             schema.TypeString,
		Optional:         true,
		RequiredWith:     []string{"frequency"},
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("start_time"),
		Description:      "The date and time when the resource monitor starts monitoring credit usage for the assigned warehouses. If you set a `start_timestamp` for a resource monitor, you must also set `frequency`.",
	},
	// TODO: Describe that it's advised for now to specify full dates of format 2024-10-04 00:00 otherwise diffs may occur
	"end_timestamp": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("end_time"),
		Description:      "The date and time when the resource monitor suspends the assigned warehouses.",
	},
	"trigger": {
		Type:     schema.TypeSet,
		Optional: true,
		// TODO: Throw error on CREATE with only triggers (SQL compilation error).
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
		ReadContext:   ReadResourceMonitor(true),
		UpdateContext: UpdateResourceMonitor,
		DeleteContext: DeleteResourceMonitor,

		Schema: resourceMonitorSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportResourceMonitor,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(resourceMonitorSchema, ShowOutputAttributeName, "notify_users", "credit_quota", "frequency", "start_timestamp", "end_timestamp", "trigger"),
		),
	}
}

func ImportResourceMonitor(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting resource monitor import")
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := d.Set("name", id.Name()); err != nil {
		return nil, err
	}
	if err := d.Set("credit_quota", resourceMonitor.CreditQuota); err != nil {
		return nil, err
	}
	if err := d.Set("frequency", resourceMonitor.Frequency); err != nil {
		return nil, err
	}
	if err := d.Set("start_timestamp", resourceMonitor.StartTime); err != nil {
		return nil, err
	}
	if err := d.Set("end_timestamp", resourceMonitor.EndTime); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateResourceMonitor(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	opts := new(sdk.CreateResourceMonitorOptions)
	with := new(sdk.ResourceMonitorWith)

	if v, ok := d.GetOk("credit_quota"); ok {
		with.CreditQuota = sdk.Pointer(v.(int))
	}

	if v, ok := d.GetOk("notify_users"); ok {
		userIds := expandStringList(v.(*schema.Set).List())
		users := make([]sdk.NotifiedUser, len(userIds))
		for i, userId := range userIds {
			users[i] = sdk.NotifiedUser{
				Name: sdk.NewAccountObjectIdentifier(userId),
			}
		}
		with.NotifyUsers = &sdk.NotifyUsers{Users: users}
	}

	if v, ok := d.GetOk("frequency"); ok {
		frequency, err := sdk.ToResourceMonitorFrequency(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		with.Frequency = frequency
	}

	if v, ok := d.GetOk("start_timestamp"); ok {
		with.StartTimestamp = sdk.Pointer(v.(string))
	}

	if v, ok := d.GetOk("end_timestamp"); ok {
		with.EndTimestamp = sdk.Pointer(v.(string))
	}

	if v, ok := d.GetOk("trigger"); ok {
		triggerDefinitions, err := extractTriggerDefinitions(v.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		with.Triggers = triggerDefinitions
	}

	if !reflect.DeepEqual(*with, sdk.ResourceMonitorWith{}) {
		opts.With = with
	}

	err := client.ResourceMonitors.Create(ctx, id, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadResourceMonitor(false)(ctx, d, meta)
}

func ReadResourceMonitor(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) { // TODO: Test for sdk.ErrObjectNotExistOrAuthorized
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

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"credit_quota", "credit_quota", resourceMonitor.CreditQuota, resourceMonitor.CreditQuota, nil},
				showMapping{"frequency", "frequency", string(resourceMonitor.Frequency), resourceMonitor.Frequency, nil},
				showMapping{"start_time", "start_timestamp", resourceMonitor.StartTime, resourceMonitor.StartTime, nil},
				showMapping{"end_time", "end_timestamp", resourceMonitor.EndTime, resourceMonitor.EndTime, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, warehouseSchema, []string{
			"credit_quota",
			"frequency",
			"start_timestamp",
			"end_timestamp",
		}); err != nil {
			return diag.FromErr(err)
		}

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
}

func UpdateResourceMonitor(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var runSetStatement bool
	var runUnsetStatement bool
	opts := sdk.AlterResourceMonitorOptions{}
	set := sdk.ResourceMonitorSet{}
	unset := sdk.ResourceMonitorUnset{}

	if d.HasChange("credit_quota") {
		runSetStatement = true
		if v, ok := d.GetOk("credit_quota"); ok {
			set.CreditQuota = sdk.Pointer(v.(int))
		} else {
			runUnsetStatement = true
			unset.CreditQuota = sdk.Bool(true)
		}
	}

	if (d.HasChange("frequency") || d.HasChange("start_timestamp")) &&
		(d.Get("frequency").(string) != "" && d.Get("start_timestamp").(string) != "") {
		runSetStatement = true
		frequency, err := sdk.ToResourceMonitorFrequency(d.Get("frequency").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		set.Frequency = frequency
		set.StartTimestamp = sdk.Pointer(d.Get("start_timestamp").(string))
	}

	if d.HasChange("end_timestamp") {
		if v, ok := d.GetOk("end_timestamp"); ok {
			runSetStatement = true
			set.EndTimestamp = sdk.Pointer(v.(string))
		} else {
			runUnsetStatement = true
			unset.EndTimestamp = sdk.Bool(true)
		}
	}

	if d.HasChange("notify_users") {
		userIds := expandStringList(d.Get("notify_users").(*schema.Set).List())
		if len(userIds) > 0 {
			runSetStatement = true
			users := make([]sdk.NotifiedUser, len(userIds))
			for i, userId := range userIds {
				users[i] = sdk.NotifiedUser{
					Name: sdk.NewAccountObjectIdentifier(userId),
				}
			}
			set.NotifyUsers = &sdk.NotifyUsers{
				Users: users,
			}
		} else {
			runUnsetStatement = true
			unset.NotifyUsers = sdk.Bool(true)
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

	if runUnsetStatement {
		if unset != (sdk.ResourceMonitorUnset{}) {
			opts.Unset = &unset
		}
		if err := client.ResourceMonitors.Alter(ctx, id, &opts); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadResourceMonitor(false)(ctx, d, meta)
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
