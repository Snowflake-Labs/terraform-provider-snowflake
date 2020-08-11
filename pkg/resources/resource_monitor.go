package resources

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
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
		Type:        schema.TypeFloat,
		Optional:    true,
		Computed:    true,
		Description: "The amount of credits allocated monthly to the resource monitor, round up to 2 decimal places.",
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
}

// ResourceMonitor returns a pointer to the resource representing a resource monitor
func ResourceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: CreateResourceMonitor,
		Read:   ReadResourceMonitor,
		// Update: UpdateResourceMonitor, @TODO implement updates
		Delete: DeleteResourceMonitor,
		Exists: ResourceMonitorExists,

		Schema: resourceMonitorSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateResourceMonitor implents schema.CreateFunc
func CreateResourceMonitor(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	cb := snowflake.ResourceMonitor(name).Create()
	// Set optionals
	if v, ok := data.GetOk("credit_quota"); ok {
		cb.SetFloat("credit_quota", v.(float64))
	}
	if v, ok := data.GetOk("frequency"); ok {
		cb.SetString("frequency", v.(string))
	}
	if v, ok := data.GetOk("start_timestamp"); ok {
		cb.SetString("start_timestamp", v.(string))
	}
	if v, ok := data.GetOk("end_timestamp"); ok {
		cb.SetString("end_timestamp", v.(string))
	}
	// Set triggers
	sTrigs := expandIntList(data.Get("suspend_triggers").(*schema.Set).List())
	for _, t := range sTrigs {
		cb.SuspendAt(t)
	}
	siTrigs := expandIntList(data.Get("suspend_immediate_triggers").(*schema.Set).List())
	for _, t := range siTrigs {
		cb.SuspendImmediatelyAt(t)
	}
	nTrigs := expandIntList(data.Get("notify_triggers").(*schema.Set).List())
	for _, t := range nTrigs {
		cb.NotifyAt(t)
	}

	stmt := cb.Statement()

	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating resource monitor %v", name)
	}

	data.SetId(name)

	return ReadResourceMonitor(data, meta)
}

// ReadResourceMonitor implements schema.ReadFunc
func ReadResourceMonitor(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stmt := snowflake.ResourceMonitor(data.Id()).Show()

	row := snowflake.QueryRow(db, stmt)

	rm, err := snowflake.ScanResourceMonitor(row)
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
	err = setDataFromNullStrings(data, nullStrings)
	if err != nil {
		return err
	}

	// Credit quota is a float
	if rm.CreditQuota.Valid {
		err = data.Set("credit_quota", rm.CreditQuota.Float64)
	} else {
		err = data.Set("credit_quota", 0.0) // not sure if this is the right approach
	}
	if err != nil {
		return err
	}

	// Triggers
	sTrigs, err := extractTriggerInts(rm.SuspendAt)
	if err != nil {
		return err
	}
	err = data.Set("suspend_triggers", sTrigs)
	if err != nil {
		return err
	}
	siTrigs, err := extractTriggerInts(rm.SuspendImmediatelyAt)
	if err != nil {
		return err
	}
	err = data.Set("suspend_immediate_triggers", siTrigs)
	if err != nil {
		return err
	}
	nTrigs, err := extractTriggerInts(rm.NotifyAt)
	if err != nil {
		return err
	}
	err = data.Set("notify_triggers", nTrigs)

	return err
}

// setDataFromNullString blanks the value if v is null, otherwise sets the value to the value of v
func setDataFromNullStrings(data *schema.ResourceData, ns map[string]sql.NullString) error {
	for k, v := range ns {
		var err error
		if v.Valid {
			err = data.Set(k, v.String)
		} else {
			err = data.Set(k, "")
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
func DeleteResourceMonitor(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	stmt := snowflake.ResourceMonitor(data.Id()).Drop()

	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error deleting resource monitor %v", data.Id())
	}

	data.SetId("")
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
