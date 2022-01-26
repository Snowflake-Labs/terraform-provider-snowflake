package resources

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"
)

var validFrequencies = []string{"MONTHLY", "DAILY", "WEEKLY", "YEARLY", "NEVER"}

var resourceMonitorSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the resource monitor; must be unique for your account.",
		ForceNew:    true,
	},
	"credit_quota": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "The number of credits allocated monthly to the resource monitor.",
		ForceNew:    true,
	},
	"frequency": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "The frequency interval at which the credit usage resets to 0. If you set a frequency for a resource monitor, you must also set START_TIMESTAMP.",
		ValidateFunc: validation.StringInSlice(validFrequencies, false),
		ForceNew:     true,
	},
	"start_timestamp": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The date and time when the resource monitor starts monitoring credit usage for the assigned warehouses.",
		ForceNew:    true,
	},
	"end_timestamp": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The date and time when the resource monitor suspends the assigned warehouses.",
		ForceNew:    true,
	},
	"suspend_triggers": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Optional:    true,
		Description: "A list of percentage thresholds at which to suspend all warehouses.",
		ForceNew:    true,
	},
	"suspend_immediate_triggers": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Optional:    true,
		Description: "A list of percentage thresholds at which to immediately suspend all warehouses.",
		ForceNew:    true,
	},
	"notify_triggers": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Optional:    true,
		Description: "A list of percentage thresholds at which to send an alert to subscribed users.",
		ForceNew:    true,
	},
	"set_for_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the resource monitor should be applied globally to your Snowflake account.",
		Default:     false,
		ForceNew:    true,
	},
	"warehouses": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "A list of warehouses to apply the resource monitor to.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		ForceNew:    true,
	},
}

// ResourceMonitor returns a pointer to the resource representing a resource monitor
func ResourceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: CreateResourceMonitor,
		Read:   ReadResourceMonitor,
		// Update: UpdateResourceMonitor, @TODO implement updates
		Delete: DeleteResourceMonitor,

		Schema: resourceMonitorSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateResourceMonitor implents schema.CreateFunc
func CreateResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	cb := snowflake.ResourceMonitor(name).Create()
	// Set optionals
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
	sTrigs := expandIntList(d.Get("suspend_triggers").(*schema.Set).List())
	for _, t := range sTrigs {
		cb.SuspendAt(t)
	}
	siTrigs := expandIntList(d.Get("suspend_immediate_triggers").(*schema.Set).List())
	for _, t := range siTrigs {
		cb.SuspendImmediatelyAt(t)
	}
	nTrigs := expandIntList(d.Get("notify_triggers").(*schema.Set).List())
	for _, t := range nTrigs {
		cb.NotifyAt(t)
	}

	stmt := cb.Statement()

	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating resource monitor %v", name)
	}

	d.SetId(name)

	if d.Get("set_for_account").(bool) {
		if err := snowflake.Exec(db, cb.SetOnAccount()); err != nil {
			return errors.Wrapf(err, "error setting resource monitor %v on account", name)
		}
	}

	if v, ok := d.GetOk("warehouses"); ok {
		for _, w := range v.(*schema.Set).List() {
			if err := snowflake.Exec(db, cb.SetOnWarehouse(w.(string))); err != nil {
				return errors.Wrapf(err, "error setting resource monitor %v on warehouse %v", name, w.(string))
			}
		}
	}

	if err := ReadResourceMonitor(d, meta); err != nil {
		return err
	}

	return nil
}

// ReadResourceMonitor implements schema.ReadFunc
func ReadResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stmt := snowflake.ResourceMonitor(d.Id()).Show()

	row := snowflake.QueryRow(db, stmt)

	rm, err := snowflake.ScanResourceMonitor(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
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
	err = setDataFromNullStrings(d, nullStrings)
	if err != nil {
		return err
	}

	// Snowflake returns credit_quota as a float, but only accepts input as an int
	if rm.CreditQuota.Valid {
		cqf, err := strconv.ParseFloat(rm.CreditQuota.String, 64)
		if err != nil {
			return err
		}

		err = d.Set("credit_quota", int(cqf))
		if err != nil {
			return err
		}
	}

	// Triggers
	sTrigs, err := extractTriggerInts(rm.SuspendAt)
	if err != nil {
		return err
	}
	err = d.Set("suspend_triggers", sTrigs)
	if err != nil {
		return err
	}
	siTrigs, err := extractTriggerInts(rm.SuspendImmediatelyAt)
	if err != nil {
		return err
	}
	err = d.Set("suspend_immediate_triggers", siTrigs)
	if err != nil {
		return err
	}
	nTrigs, err := extractTriggerInts(rm.NotifyAt)
	if err != nil {
		return err
	}
	err = d.Set("notify_triggers", nTrigs)

	// Account level
	d.Set("set_for_account", rm.Level.Valid && rm.Level.String == "ACCOUNT")

	return err
}

// setDataFromNullString blanks the value if v is null, otherwise sets the value to the value of v
func setDataFromNullStrings(data *schema.ResourceData, ns map[string]sql.NullString) error {
	for k, v := range ns {
		var err error
		if v.Valid {
			err = data.Set(k, v.String) //lintignore:R001
		} else {
			err = data.Set(k, "") //lintignore:R001
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// extractTriggerInts converts the triggers in the DB (stored as a comma
// separated string with trailling %s) into a slice of ints
func extractTriggerInts(s sql.NullString) ([]int, error) {
	// Check if this is NULL
	if !s.Valid {
		return []int{}, nil
	}
	ints := strings.Split(s.String, ",")
	var out []int
	for _, i := range ints {
		myInt, err := strconv.Atoi(i[:len(i)-1])
		if err != nil {
			return out, errors.Wrapf(err, "failed to convert %v to integer", i)
		}
		out = append(out, myInt)
	}
	return out, nil
}

// DeleteResourceMonitor implements schema.DeleteFunc
func DeleteResourceMonitor(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	stmt := snowflake.ResourceMonitor(d.Id()).Drop()

	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error deleting resource monitor %v", d.Id())
	}

	d.SetId("")
	return nil
}

// ResourceMonitorExists implements schema.ExistsFunc
func ResourceMonitorExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)

	q := snowflake.ResourceMonitor(data.Id()).Show()

	rows, err := db.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
