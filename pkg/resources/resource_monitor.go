package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	"frequency": {
		Type:             schema.TypeString,
		Optional:         true,
		RequiredWith:     []string{"start_timestamp"},
		ValidateDiagFunc: sdkValidation(sdk.ToResourceMonitorFrequency),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToResourceMonitorFrequency), IgnoreChangeToCurrentSnowflakeValueInShow("frequency")),
		Description:      fmt.Sprintf("The frequency interval at which the credit usage resets to 0. Valid values are (case-insensitive): %s. If you set a `frequency` for a resource monitor, you must also set `start_timestamp`. If you specify `NEVER` for the frequency, the credit usage for the warehouse does not reset. After removing this field from the config, the previously set value will be preserved on the Snowflake side, not the default value. That's due to Snowflake limitation and the lack of unset functionality for this parameter.", possibleValuesListed(sdk.AllFrequencyValues)),
	},
	"start_timestamp": {
		Type:             schema.TypeString,
		Optional:         true,
		RequiredWith:     []string{"frequency"},
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("start_time"),
		Description:      "The date and time when the resource monitor starts monitoring credit usage for the assigned warehouses. If you set a `start_timestamp` for a resource monitor, you must also set `frequency`.  After removing this field from the config, the previously set value will be preserved on the Snowflake side, not the default value. That's due to Snowflake limitation and the lack of unset functionality for this parameter.",
	},
	"end_timestamp": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("end_time"),
		Description:      "The date and time when the resource monitor suspends the assigned warehouses.",
	},
	"notify_triggers": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Specifies a list of percentages of the credit quota. After reaching any of the values the users passed in the notify_users field will be notified (to receive the notification they should have notifications enabled). Values over 100 are supported.",
		Elem: &schema.Schema{
			Type:             schema.TypeInt,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		},
	},
	"suspend_trigger": {
		Type:             schema.TypeInt,
		Optional:         true,
		Description:      "Represents a numeric value specified as a percentage of the credit quota. Values over 100 are supported. After reaching this value, all assigned warehouses while allowing currently running queries to complete will be suspended. No new queries can be executed by the warehouses until the credit quota for the resource monitor is increased. In addition, this action sends a notification to all users who have enabled notifications for themselves.",
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
	},
	"suspend_immediate_trigger": {
		Type:             schema.TypeInt,
		Optional:         true,
		Description:      "Represents a numeric value specified as a percentage of the credit quota. Values over 100 are supported. After reaching this value, all assigned warehouses immediately cancel any currently running queries or statements. In addition, this action sends a notification to all users who have enabled notifications for themselves.",
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
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
			ComputedIfAnyAttributeChanged(resourceMonitorSchema, ShowOutputAttributeName, "notify_users", "credit_quota", "frequency", "start_timestamp", "end_timestamp", "notify_triggers", "suspend_trigger", "suspend_immediate_trigger"),
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
	if err := d.Set("notify_triggers", resourceMonitor.NotifyAt); err != nil {
		return nil, err
	}
	if err := d.Set("suspend_trigger", resourceMonitor.SuspendAt); err != nil {
		return nil, err
	}
	if err := d.Set("suspend_immediate_trigger", resourceMonitor.SuspendImmediateAt); err != nil {
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

	triggers := make([]sdk.TriggerDefinition, 0)
	if notifyTriggers, ok := d.GetOk("notify_triggers"); ok {
		for _, triggerThreshold := range notifyTriggers.(*schema.Set).List() {
			triggers = append(triggers, sdk.TriggerDefinition{
				Threshold:     triggerThreshold.(int),
				TriggerAction: sdk.TriggerActionNotify,
			})
		}
	}

	if suspendTriggerThreshold, ok := d.GetOk("suspend_trigger"); ok {
		triggers = append(triggers, sdk.TriggerDefinition{
			Threshold:     suspendTriggerThreshold.(int),
			TriggerAction: sdk.TriggerActionSuspend,
		})
	}

	if suspendImmediateTriggerThreshold, ok := d.GetOk("suspend_immediate_trigger"); ok {
		triggers = append(triggers, sdk.TriggerDefinition{
			Threshold:     suspendImmediateTriggerThreshold.(int),
			TriggerAction: sdk.TriggerActionSuspendImmediate,
		})
	}

	if len(triggers) > 0 {
		with.Triggers = triggers
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

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"credit_quota", "credit_quota", resourceMonitor.CreditQuota, resourceMonitor.CreditQuota, nil},
				showMapping{"frequency", "frequency", string(resourceMonitor.Frequency), resourceMonitor.Frequency, nil},
				showMapping{"start_time", "start_timestamp", resourceMonitor.StartTime, resourceMonitor.StartTime, nil},
				showMapping{"end_time", "end_timestamp", resourceMonitor.EndTime, resourceMonitor.EndTime, nil},
				showMapping{"notify_at", "notify_triggers", resourceMonitor.NotifyAt, resourceMonitor.NotifyAt, nil},
				showMapping{"suspend_at", "suspend_trigger", resourceMonitor.SuspendAt, resourceMonitor.SuspendAt, nil},
				showMapping{"suspend_immediately_at", "suspend_immediate_trigger", resourceMonitor.SuspendImmediateAt, resourceMonitor.SuspendImmediateAt, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, warehouseSchema, []string{
			"credit_quota",
			"frequency",
			"start_timestamp",
			"end_timestamp",
			"notify_triggers",
			"suspend_trigger",
			"suspend_immediate_trigger",
		}); err != nil {
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

	if d.HasChanges("notify_triggers", "suspend_trigger", "suspend_immediate_trigger") {
		triggers := make([]sdk.TriggerDefinition, 0)

		if notifyTriggers, ok := d.GetOk("notify_triggers"); ok {
			for _, triggerThreshold := range notifyTriggers.(*schema.Set).List() {
				triggers = append(triggers, sdk.TriggerDefinition{
					Threshold:     triggerThreshold.(int),
					TriggerAction: sdk.TriggerActionNotify,
				})
			}
		}

		if suspendTriggerThreshold, ok := d.GetOk("suspend_trigger"); ok {
			triggers = append(triggers, sdk.TriggerDefinition{
				Threshold:     suspendTriggerThreshold.(int),
				TriggerAction: sdk.TriggerActionSuspend,
			})
		}

		if suspendImmediateTriggerThreshold, ok := d.GetOk("suspend_immediate_trigger"); ok {
			triggers = append(triggers, sdk.TriggerDefinition{
				Threshold:     suspendImmediateTriggerThreshold.(int),
				TriggerAction: sdk.TriggerActionSuspendImmediate,
			})
		}

		if len(triggers) > 0 {
			opts.Triggers = triggers
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

	// This is to prevent SQL compilation errors from Snowflake, because you cannot only alter triggers.
	// It's going to set credit quota to the same value as before making it pass SQL compilation stage.
	if len(opts.Triggers) > 0 && !runSetStatement && !runUnsetStatement {
		if creditQuota, ok := d.GetOk("credit_quota"); ok {
			runSetStatement = true
			set.CreditQuota = sdk.Pointer(creditQuota.(int))
		} else {
			runUnsetStatement = true
			unset.CreditQuota = sdk.Bool(true)
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
