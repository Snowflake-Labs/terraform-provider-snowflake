package resources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The number that represents the percentage threshold at which to suspend all warehouses.",
	},
	"suspend_immediate_trigger": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The number that represents the percentage threshold at which to immediately suspend all warehouses.",
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
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	check := checkAccountAgainstWarehouses(d, name)

	if check != nil {
		return check
	}

	cb := snowflake.NewResourceMonitorBuilder(name).Create()
	// Set optionals
	if v, ok := d.GetOk("notify_users"); ok {
		cb.SetStringList("notify_users", expandStringList(v.(*schema.Set).List()))
	}
	if v, ok := d.GetOk("credit_quota"); ok {
		cb.SetInt("credit_quota", v.(int))
	}
	if v, ok := d.GetOk("frequency"); ok {
		cb.SetString("frequency", v.(string))
	}
	if v, ok := d.GetOk("start_timestamp"); ok {
		cb.SetString("start_timestamp", v.(string))
	}
	if v, ok := d.GetOk("end_timestamp"); ok {
		cb.SetString("end_timestamp", v.(string))
	}
	// Set triggers
	if v, ok := d.GetOk("suspend_trigger"); ok {
		cb.SuspendAt(v.(int))
	}
	if v, ok := d.GetOk("suspend_immediate_trigger"); ok {
		cb.SuspendImmediatelyAt(v.(int))
	}
	nTrigs := expandIntList(d.Get("notify_triggers").(*schema.Set).List())
	for _, t := range nTrigs {
		cb.NotifyAt(t)
	}

	stmt := cb.Statement()
	if err := snowflake.Exec(db, stmt); err != nil {
		return fmt.Errorf("error creating resource monitor %v err = %w", name, err)
	}

	d.SetId(name)

	if d.Get("set_for_account").(bool) {
		if err := snowflake.Exec(db, cb.SetOnAccount()); err != nil {
			return fmt.Errorf("error setting resource monitor %v on account err = %w", name, err)
		}
	}

	if v, ok := d.GetOk("warehouses"); ok {
		for _, w := range v.(*schema.Set).List() {
			if err := snowflake.Exec(db, cb.SetOnWarehouse(w.(string))); err != nil {
				return fmt.Errorf("error setting resource monitor %v on warehouse %v err = %w", name, w.(string), err)
			}
		}
	}

	if err := ReadResourceMonitor(d, meta); err != nil {
		return err
	}

	return nil
}

// ReadResourceMonitor implements schema.ReadFunc.
func ReadResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stmt := snowflake.NewResourceMonitorBuilder(d.Id()).Show()

	row := snowflake.QueryRow(db, stmt)

	rm, err := snowflake.ScanResourceMonitor(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] resource monitor (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	// Set string values
	nullStrings := map[string]sql.NullString{
		"name":            rm.Name,
		"frequency":       rm.Frequency,
		"start_timestamp": rm.StartTime,
		"end_timestamp":   rm.EndTime,
	}
	if err := setDataFromNullStrings(d, nullStrings); err != nil {
		return err
	}

	if len(rm.NotifyUsers.String) > 0 {
		if err := d.Set("notify_users", strings.Split(rm.NotifyUsers.String, ", ")); err != nil {
			return err
		}
	}

	// Snowflake returns credit_quota as a float, but only accepts input as an int
	if rm.CreditQuota.Valid {
		cqf, err := strconv.ParseFloat(rm.CreditQuota.String, 64)
		if err != nil {
			return err
		}
		if err := d.Set("credit_quota", int(cqf)); err != nil {
			return err
		}
	}

	// Triggers
	sTrig, err := extractTriggerInts(rm.SuspendAt)
	if err != nil {
		return err
	}
	if err := d.Set("suspend_trigger", sTrig); err != nil {
		return err
	}
	siTrig, err := extractTriggerInts(rm.SuspendImmediatelyAt)
	if err != nil {
		return err
	}
	if err := d.Set("suspend_immediate_trigger", siTrig); err != nil {
		return err
	}
	nTrigs, err := extractTriggerInts(rm.NotifyAt)
	if err != nil {
		return err
	}
	if err := d.Set("notify_triggers", nTrigs); err != nil {
		return err
	}

	// Account level
	if err := d.Set("set_for_account", rm.Level.Valid && rm.Level.String == "ACCOUNT"); err != nil {
		return err
	}

	return err
}

// setDataFromNullString blanks the value if v is null, otherwise sets the value to the value of v.
func setDataFromNullStrings(data *schema.ResourceData, ns map[string]sql.NullString) error {
	for k, v := range ns {
		var err error
		if v.Valid {
			err = data.Set(k, v.String) // lintignore:R001
		} else {
			err = data.Set(k, "") // lintignore:R001
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// extractTriggerInts converts the triggers in the DB (stored as a comma
// separated string with trailing %s) into a slice of ints.
func extractTriggerInts(s sql.NullString) ([]int, error) {
	// Check if this is NULL
	if !s.Valid {
		return []int{}, nil
	}
	ints := strings.Split(s.String, ",")
	out := make([]int, 0, len(ints))
	for _, i := range ints {
		myInt, err := strconv.Atoi(i[:len(i)-1])
		if err != nil {
			return out, fmt.Errorf("failed to convert %v to integer err = %w", i, err)
		}
		out = append(out, myInt)
	}
	return out, nil
}

// UpdateResourceMonitor implements schema.UpdateFunc.
func UpdateResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	check := checkAccountAgainstWarehouses(d, id)

	if check != nil {
		return check
	}

	ub := snowflake.NewResourceMonitorBuilder(id).Alter()
	var runSetStatement bool

	if d.HasChange("notify_users") {
		runSetStatement = true
		ub.SetStringList(`NOTIFY_USERS`, expandStringList(d.Get("notify_users").(*schema.Set).List()))
	}

	if d.HasChange("credit_quota") {
		runSetStatement = true
		ub.SetInt(`CREDIT_QUOTA`, d.Get("credit_quota").(int))
	}

	if d.HasChange("frequency") {
		runSetStatement = true
		ub.SetString(`FREQUENCY`, d.Get("frequency").(string))
	}

	if d.HasChange("start_timestamp") {
		runSetStatement = true
		ub.SetString(`START_TIMESTAMP`, d.Get("start_timestamp").(string))
	}

	if d.HasChange("end_timestamp") {
		runSetStatement = true
		ub.SetString(`END_TIMESTAMP`, d.Get("end_timestamp").(string))
	}

	// Set triggers
	if d.HasChange("suspend_trigger") {
		runSetStatement = true
		ub.SuspendAt(d.Get("suspend_trigger").(int))
	}
	if d.HasChange("suspend_immediate_trigger") {
		runSetStatement = true
		ub.SuspendImmediatelyAt(d.Get("suspend_immediate_trigger").(int))
	}
	nTrigs := expandIntList(d.Get("notify_triggers").(*schema.Set).List())
	for _, t := range nTrigs {
		runSetStatement = true
		ub.NotifyAt(t)
	}

	if runSetStatement {
		if err := snowflake.Exec(db, ub.Statement()); err != nil {
			return fmt.Errorf("error updating resource monitor %v\n%w", id, err)
		}
	}

	// Remove from account
	if d.HasChange("set_for_account") && !d.Get("set_for_account").(bool) {
		if err := snowflake.Exec(db, ub.UnsetOnAccount()); err != nil {
			return fmt.Errorf("error unsetting resource monitor %v on account err = %w", id, err)
		}
	}

	// Remove from all old warehouses
	if d.HasChange("warehouses") {
		oldV, v := d.GetChange("warehouses")
		res := intersectionAAndNotB(oldV.(*schema.Set).List(), v.(*schema.Set).List())
		for _, w := range res {
			if err := snowflake.Exec(db, ub.UnsetOnWarehouse(w)); err != nil {
				return fmt.Errorf("error setting resource monitor %v on warehouse %v err = %w", id, w, err)
			}
		}
	}

	// Add to account
	if d.HasChange("set_for_account") && d.Get("set_for_account").(bool) {
		if err := snowflake.Exec(db, ub.SetOnAccount()); err != nil {
			return fmt.Errorf("error setting resource monitor %v on account err = %w", id, err)
		}
	}

	// Add to all new warehouses
	if d.HasChange("warehouses") {
		oldV, v := d.GetChange("warehouses")
		res := intersectionAAndNotB(v.(*schema.Set).List(), oldV.(*schema.Set).List())
		for _, w := range res {
			if err := snowflake.Exec(db, ub.SetOnWarehouse(w)); err != nil {
				return fmt.Errorf("error setting resource monitor %v on warehouse %v err = %w", id, w, err)
			}
		}
	}

	return ReadResourceMonitor(d, meta)
}

// DeleteResourceMonitor implements schema.DeleteFunc.
func DeleteResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	stmt := snowflake.NewResourceMonitorBuilder(d.Id()).Drop()
	if err := snowflake.Exec(db, stmt); err != nil {
		return fmt.Errorf("error deleting resource monitor %v err = %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}
