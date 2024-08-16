package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validFrequencies = []string{"MONTHLY", "DAILY", "WEEKLY", "YEARLY", "NEVER"}

var resourceMonitorSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the resource monitor; must be unique for your account.",
		ForceNew:    true,
	},
	"notify_users": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Specifies the list of users to receive email notifications on resource monitors.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"credit_quota": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "The number of credits allocated monthly to the resource monitor.",
	},
	"frequency": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "The frequency interval at which the credit usage resets to 0. If you set a frequency for a resource monitor, you must also set START_TIMESTAMP.",
		ValidateFunc: validation.StringInSlice(validFrequencies, false),
	},
	"start_timestamp": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The date and time when the resource monitor starts monitoring credit usage for the assigned warehouses.",
	},
	"end_timestamp": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The date and time when the resource monitor suspends the assigned warehouses.",
	},
	"suspend_trigger": {
		Type:          schema.TypeInt,
		Optional:      true,
		Description:   "The number that represents the percentage threshold at which to suspend all warehouses.",
		ConflictsWith: []string{"suspend_triggers"},
	},
	"suspend_triggers": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeInt},
		Optional:      true,
		Description:   "A list of percentage thresholds at which to suspend all warehouses.",
		ConflictsWith: []string{"suspend_trigger"},
		Deprecated:    "Use suspend_trigger instead",
	},
	"suspend_immediate_trigger": {
		Type:          schema.TypeInt,
		Optional:      true,
		Description:   "The number that represents the percentage threshold at which to immediately suspend all warehouses.",
		ConflictsWith: []string{"suspend_immediate_triggers"},
	},
	"suspend_immediate_triggers": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeInt},
		Optional:      true,
		Description:   "A list of percentage thresholds at which to suspend all warehouses.",
		ConflictsWith: []string{"suspend_immediate_trigger"},
		Deprecated:    "Use suspend_immediate_trigger instead",
	},
	"notify_triggers": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Optional:    true,
		Description: "A list of percentage thresholds at which to send an alert to subscribed users.",
	},
	"set_for_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the resource monitor should be applied globally to your Snowflake account (defaults to false).",
		Default:     false,
	},
	"warehouses": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "A list of warehouses to apply the resource monitor to.",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// ResourceMonitor returns a pointer to the resource representing a resource monitor.
func ResourceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: CreateResourceMonitor,
		Read:   ReadResourceMonitor,
		Update: UpdateResourceMonitor,
		Delete: DeleteResourceMonitor,

		Schema: resourceMonitorSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func checkAccountAgainstWarehouses(d *schema.ResourceData, name string) error {
	account := d.Get("set_for_account").(bool)
	v := d.Get("warehouses")

	if len(v.(*schema.Set).List()) > 0 && account {
		return fmt.Errorf("error creating resource monitor %v on account err = set_for_account cannot be true and give warehouses", name)
	}
	return nil
}

// CreateResourceMonitor implements schema.CreateFunc.
func CreateResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)

	check := checkAccountAgainstWarehouses(d, name)

	if check != nil {
		return check
	}

	ctx := context.Background()
	objectIdentifier := sdk.NewAccountObjectIdentifier(name)

	// Set optionals.
	opts := &sdk.CreateResourceMonitorOptions{}
	if v, ok := d.GetOk("notify_users"); ok {
		userNames := expandStringList(v.(*schema.Set).List())
		users := []sdk.NotifiedUser{}
		for _, name := range userNames {
			users = append(users, sdk.NotifiedUser{Name: name})
		}
		if opts.With == nil {
			opts.With = &sdk.ResourceMonitorWith{}
		}
		opts.With.NotifyUsers = &sdk.NotifyUsers{Users: users}
	}

	if v, ok := d.GetOk("credit_quota"); ok {
		if opts.With == nil {
			opts.With = &sdk.ResourceMonitorWith{}
		}
		opts.With.CreditQuota = sdk.Int(v.(int))
	}
	if v, ok := d.GetOk("frequency"); ok {
		frequency, err := sdk.FrequencyFromString(v.(string))
		if err != nil {
			return err
		}
		if opts.With == nil {
			opts.With = &sdk.ResourceMonitorWith{}
		}
		opts.With.Frequency = frequency
	}

	if v, ok := d.GetOk("start_timestamp"); ok {
		if opts.With == nil {
			opts.With = &sdk.ResourceMonitorWith{}
		}
		opts.With.StartTimestamp = sdk.Pointer(v.(string))
	}
	if v, ok := d.GetOk("end_timestamp"); ok {
		if opts.With == nil {
			opts.With = &sdk.ResourceMonitorWith{}
		}
		opts.With.EndTimestamp = sdk.Pointer(v.(string))
	}

	triggers := collectResourceMonitorTriggers(d)
	if len(triggers) > 0 {
		if opts.With == nil {
			opts.With = &sdk.ResourceMonitorWith{}
		}
		opts.With.Triggers = triggers
	}

	err := client.ResourceMonitors.Create(ctx, objectIdentifier, opts)
	if err != nil {
		return fmt.Errorf("error creating resource monitor %v err = %w", name, err)
	}
	d.SetId(name)

	if d.Get("set_for_account").(bool) {
		accountOpts := sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				ResourceMonitor: objectIdentifier,
			},
		}
		if err := client.Accounts.Alter(ctx, &accountOpts); err != nil {
			return fmt.Errorf("error setting resource monitor %v on account err = %w", name, err)
		}
	}

	if v, ok := d.GetOk("warehouses"); ok {
		for _, w := range v.(*schema.Set).List() {
			warehouseOpts := sdk.AlterWarehouseOptions{
				Set: &sdk.WarehouseSet{
					ResourceMonitor: objectIdentifier,
				},
			}
			warehouseId := sdk.NewAccountObjectIdentifier(w.(string))
			if err := client.Warehouses.Alter(ctx, warehouseId, &warehouseOpts); err != nil {
				return fmt.Errorf("error setting resource monitor %v on warehouse %v err = %w", name, warehouseId.Name(), err)
			}
		}
	}

	return ReadResourceMonitor(d, meta)
}

// ReadResourceMonitor implements schema.ReadFunc.
func ReadResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	ctx := context.Background()
	resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, id)
	if err != nil {
		return err
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return err
	}

	if err := d.Set("name", resourceMonitor.Name); err != nil {
		return err
	}
	if err := d.Set("frequency", string(resourceMonitor.Frequency)); err != nil {
		return err
	}

	if err := d.Set("start_timestamp", resourceMonitor.StartTime); err != nil {
		return err
	}

	if err := d.Set("end_timestamp", resourceMonitor.EndTime); err != nil {
		return err
	}

	if len(resourceMonitor.NotifyUsers) > 0 {
		if err := d.Set("notify_users", resourceMonitor.NotifyUsers); err != nil {
			return err
		}
	}

	// Snowflake returns credit_quota as a float, but only accepts input as an int
	if err := d.Set("credit_quota", int(resourceMonitor.CreditQuota)); err != nil {
		return err
	}

	// Triggers
	if resourceMonitor.SuspendAt != nil {
		if err := d.Set("suspend_trigger", *resourceMonitor.SuspendAt); err != nil {
			return err
		}
	} else {
		if err := d.Set("suspend_trigger", nil); err != nil {
			return err
		}
	}

	if resourceMonitor.SuspendImmediateAt != nil {
		if err := d.Set("suspend_immediate_trigger", *resourceMonitor.SuspendImmediateAt); err != nil {
			return err
		}
	} else {
		if err := d.Set("suspend_immediate_trigger", nil); err != nil {
			return err
		}
	}

	if err := d.Set("notify_triggers", resourceMonitor.NotifyTriggers); err != nil {
		return err
	}

	// Account level
	if err := d.Set("set_for_account", resourceMonitor.Level == sdk.ResourceMonitorLevelAccount); err != nil {
		return err
	}

	return err
}

// UpdateResourceMonitor implements schema.UpdateFunc.
func UpdateResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)

	check := checkAccountAgainstWarehouses(d, name)

	if check != nil {
		return check
	}

	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	ctx := context.Background()
	var runSetStatement bool

	opts := sdk.AlterResourceMonitorOptions{}

	set := sdk.ResourceMonitorSet{}
	if d.HasChange("credit_quota") {
		runSetStatement = true
		set.CreditQuota = sdk.Pointer(d.Get("credit_quota").(int))
	}

	if d.HasChange("frequency") || d.HasChange("start_timestamp") {
		runSetStatement = true
		frequency, err := sdk.FrequencyFromString(d.Get("frequency").(string))
		if err != nil {
			return err
		}
		set.Frequency = frequency
		set.StartTimestamp = sdk.Pointer(d.Get("start_timestamp").(string))
	}

	if d.HasChange("end_timestamp") {
		runSetStatement = true
		set.EndTimestamp = sdk.Pointer(d.Get("end_timestamp").(string))
	}

	if d.HasChange("notify_users") {
		runSetStatement = true

		userNames := expandStringList(d.Get("notify_users").(*schema.Set).List())
		users := []sdk.NotifiedUser{}
		for _, name := range userNames {
			users = append(users, sdk.NotifiedUser{Name: name})
		}
		set.NotifyUsers = &sdk.NotifyUsers{
			Users: users,
		}
	}

	if set != (sdk.ResourceMonitorSet{}) {
		opts.Set = &set
	}

	// If ANY of the triggers changed, we collect all triggers and set them
	if d.HasChange("suspend_trigger") || d.HasChange("suspend_triggers") ||
		d.HasChange("suspend_immediate_trigger") || d.HasChange("suspend_immediate_triggers") ||
		d.HasChange("notify_triggers") {
		runSetStatement = true
		triggers := collectResourceMonitorTriggers(d)
		opts.Triggers = triggers
	}

	if runSetStatement {
		if err := client.ResourceMonitors.Alter(ctx, objectIdentifier, &opts); err != nil {
			return fmt.Errorf("error updating resource monitor %v\n%w", objectIdentifier.Name(), err)
		}
	}

	// Remove from account
	if d.HasChange("set_for_account") && !d.Get("set_for_account").(bool) {
		accountOpts := sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				ResourceMonitor: sdk.NewAccountObjectIdentifier("NULL"),
			},
		}
		if err := client.Accounts.Alter(ctx, &accountOpts); err != nil {
			return fmt.Errorf("error unsetting resource monitor %v on account err = %w", objectIdentifier.Name(), err)
		}
	}

	// Remove from all old warehouses
	if d.HasChange("warehouses") {
		oldV, v := d.GetChange("warehouses")
		res := ADiffB(oldV.(*schema.Set).List(), v.(*schema.Set).List())
		for _, w := range res {
			warehouseOpts := sdk.AlterWarehouseOptions{
				Unset: &sdk.WarehouseUnset{
					ResourceMonitor: sdk.Bool(true),
				},
			}
			warehouseId := sdk.NewAccountObjectIdentifier(w)
			if err := client.Warehouses.Alter(ctx, warehouseId, &warehouseOpts); err != nil {
				return fmt.Errorf("error unsetting resource monitor %v on warehouse %v err = %w", name, warehouseId.Name(), err)
			}
		}
	}

	// Add to account
	if d.HasChange("set_for_account") && d.Get("set_for_account").(bool) {
		accountOpts := sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				ResourceMonitor: objectIdentifier,
			},
		}
		if err := client.Accounts.Alter(ctx, &accountOpts); err != nil {
			return fmt.Errorf("error setting resource monitor %v on account err = %w", name, err)
		}
	}

	// Add to all new warehouses
	if d.HasChange("warehouses") {
		oldV, v := d.GetChange("warehouses")
		res := ADiffB(v.(*schema.Set).List(), oldV.(*schema.Set).List())
		for _, w := range res {
			warehouseOpts := sdk.AlterWarehouseOptions{
				Set: &sdk.WarehouseSet{
					ResourceMonitor: objectIdentifier,
				},
			}
			warehouseId := sdk.NewAccountObjectIdentifier(w)
			if err := client.Warehouses.Alter(ctx, warehouseId, &warehouseOpts); err != nil {
				return fmt.Errorf("error setting resource monitor %v on warehouse %v err = %w", name, warehouseId.Name(), err)
			}
		}
	}

	return ReadResourceMonitor(d, meta)
}

func collectResourceMonitorTriggers(d *schema.ResourceData) []sdk.TriggerDefinition {
	triggers := []sdk.TriggerDefinition{}
	var suspendTrigger *sdk.TriggerDefinition
	if v, ok := d.GetOk("suspend_trigger"); ok {
		suspendTrigger = &sdk.TriggerDefinition{
			Threshold:     v.(int),
			TriggerAction: sdk.TriggerActionSuspend,
		}
	}

	if v, ok := d.GetOk("suspend_triggers"); ok {
		siTrigs := expandIntList(v.(*schema.Set).List())
		for _, threshold := range siTrigs {
			if suspendTrigger == nil || suspendTrigger.Threshold > threshold {
				suspendTrigger = &sdk.TriggerDefinition{
					Threshold:     threshold,
					TriggerAction: sdk.TriggerActionSuspend,
				}
			}
		}
	}
	if suspendTrigger != nil {
		triggers = append(triggers, *suspendTrigger)
	}
	var suspendImmediateTrigger *sdk.TriggerDefinition

	if v, ok := d.GetOk("suspend_immediate_trigger"); ok {
		suspendImmediateTrigger = &sdk.TriggerDefinition{
			Threshold:     v.(int),
			TriggerAction: sdk.TriggerActionSuspendImmediate,
		}
	}

	if v, ok := d.GetOk("suspend_immediate_triggers"); ok {
		siTrigs := expandIntList(v.(*schema.Set).List())
		for _, threshold := range siTrigs {
			if suspendImmediateTrigger == nil || (suspendTrigger != nil && suspendTrigger.Threshold > threshold) {
				suspendImmediateTrigger = &sdk.TriggerDefinition{
					Threshold:     threshold,
					TriggerAction: sdk.TriggerActionSuspendImmediate,
				}
			}
		}
	}
	if suspendImmediateTrigger != nil {
		triggers = append(triggers, *suspendImmediateTrigger)
	}

	nTrigs := expandIntList(d.Get("notify_triggers").(*schema.Set).List())
	for _, t := range nTrigs {
		triggers = append(triggers, sdk.TriggerDefinition{
			Threshold:     t,
			TriggerAction: sdk.TriggerActionNotify,
		})
	}
	return triggers
}

// DeleteResourceMonitor implements schema.DeleteFunc.
func DeleteResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.ResourceMonitors.Drop(ctx, objectIdentifier, &sdk.DropResourceMonitorOptions{IfExists: sdk.Bool(true)})
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
