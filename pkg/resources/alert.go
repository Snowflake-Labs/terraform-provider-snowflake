package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	alertIDDelimiter = '|'
)

var alertSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the alert; must be unique for the database and schema in which the alert is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the alert.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the alert.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the alert.",
	},
	"warehouse": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The warehouse the alert will use.",
		ForceNew:    true,
	},
	"alert_schedule": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "The schedule for periodically running an alert.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cron": {
					Type:          schema.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"alert_schedule.interval"},
					Description:   "Specifies the cron expression for the alert. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"expression": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the cron expression for the alert. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
							},
							"time_zone": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the time zone for alert refresh.",
							},
						},
					},
				},
				"interval": {
					Type:          schema.TypeInt,
					Optional:      true,
					ConflictsWith: []string{"alert_schedule.cron"},
					Description:   "Specifies the interval in minutes for the alert schedule. The interval must be greater than 0 and less than 1440 (24 hours).",
				},
			},
		},
	},
	"condition": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The SQL statement that represents the condition for the alert. (SELECT, SHOW, CALL)",
		ForceNew:         false,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"action": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The SQL statement that should be executed if the condition returns one or more rows.",
		ForceNew:         false,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies if an alert should be 'started' (enabled) after creation or should remain 'suspended' (default).",
	},
}

type alertID struct {
	DatabaseName string
	SchemaName   string
	AlertName    string
}

// String() takes in a alertID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|AlertName.
func (aId *alertID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = alertIDDelimiter
	dataIdentifiers := [][]string{{aId.DatabaseName, aId.SchemaName, aId.AlertName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strAlertID := strings.TrimSpace(buf.String())
	return strAlertID, nil
}

// alertIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|AlertName
// and returns a alertID object.
func alertIDFromString(stringID string) (*alertID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = pipeIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per alert")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	alertResult := &alertID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		AlertName:    lines[0][2],
	}
	return alertResult, nil
}

// Alert returns a pointer to the resource representing an alert.
func Alert() *schema.Resource {
	return &schema.Resource{
		Create: CreateAlert,
		Read:   ReadAlert,
		Update: UpdateAlert,
		Delete: DeleteAlert,

		Schema: alertSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadAlert implements schema.ReadFunc.
func ReadAlert(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	alertID, err := alertIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := alertID.DatabaseName
	SchemaName := alertID.SchemaName
	name := alertID.AlertName

	builder := snowflake.NewAlertBuilder(name, database, SchemaName)
	qry := builder.Show()
	row := snowflake.QueryRow(db, qry)
	alert, err := snowflake.ScanAlert(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] alert (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("enabled", alert.IsEnabled()); err != nil {
		return err
	}

	if err := d.Set("name", alert.Name); err != nil {
		return err
	}

	if err := d.Set("database", alert.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("schema", alert.SchemaName); err != nil {
		return err
	}

	if err := d.Set("warehouse", alert.Warehouse); err != nil {
		return err
	}

	alertSchedule := alert.Schedule
	if alertSchedule != "" {
		if strings.Contains(alertSchedule, "MINUTE") {
			interval, err := strconv.Atoi(strings.TrimSuffix(alertSchedule, " MINUTE"))
			if err != nil {
				return err
			}
			err = d.Set("alert_schedule", []interface{}{
				map[string]interface{}{
					"interval": interval,
				},
			})
			if err != nil {
				return err
			}
		} else {
			repScheduleParts := strings.Split(alertSchedule, " ")
			timeZone := repScheduleParts[len(repScheduleParts)-1]
			expression := strings.TrimSuffix(strings.TrimPrefix(alertSchedule, "USING CRON "), " "+timeZone)
			err = d.Set("alert_schedule", []interface{}{
				map[string]interface{}{
					"cron": []interface{}{
						map[string]interface{}{
							"expression": expression,
							"time_zone":  timeZone,
						},
					},
				},
			})
			if err != nil {
				return err
			}
		}
	}

	if err := d.Set("comment", alert.Comment); err != nil {
		return err
	}

	if err := d.Set("condition", alert.Condition); err != nil {
		return err
	}

	if err := d.Set("action", alert.Action); err != nil {
		return err
	}
	return nil
}

// CreateAlert implements schema.CreateFunc.
func CreateAlert(d *schema.ResourceData, meta interface{}) error {
	var err error
	db := meta.(*sql.DB)

	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)

	builder := snowflake.NewAlertBuilder(name, database, schemaName)

	if v, ok := d.GetOk("alert_schedule"); ok {
		alertSchedule := v.([]interface{})[0].(map[string]interface{})
		if v, ok := alertSchedule["cron"]; ok {
			c := v.([]interface{})
			if len(c) > 0 {
				cron := c[0].(map[string]interface{})
				cronExpression := cron["expression"].(string)
				builder.WithAlertScheduleCronExpression(cronExpression)
				if v, ok := cron["time_zone"]; ok {
					timeZone := v.(string)
					builder.WithAlertScheduleTimeZone(timeZone)
				}
			}
		}
		if v, ok := alertSchedule["interval"]; ok {
			interval := v.(int)
			if interval > 0 {
				builder.WithAlertScheduleInterval(interval)
			}
		}
	}

	enabled := d.Get("enabled").(bool)

	warehouse := d.Get("warehouse").(string)
	builder.WithWarehouse(warehouse)

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	condition := d.Get("condition").(string)
	builder.WithCondition(condition)

	action := d.Get("action").(string)
	builder.WithAction(action)

	q := builder.Create()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error creating alert %v err = %w", name, err)
	}

	alertID := &alertID{
		DatabaseName: database,
		SchemaName:   schemaName,
		AlertName:    name,
	}
	dataIDInput, err := alertID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	if enabled {
		if err := snowflake.WaitResumeAlert(db, name, database, schemaName); err != nil {
			log.Printf("[WARN] failed to resume alert %s", name)
		}
	}

	return ReadAlert(d, meta)
}

// UpdateAlert implements schema.UpdateFunc.
func UpdateAlert(d *schema.ResourceData, meta interface{}) error {
	alertID, err := alertIDFromString(d.Id())
	if err != nil {
		return err
	}

	db := meta.(*sql.DB)
	database := alertID.DatabaseName
	schemaName := alertID.SchemaName
	name := alertID.AlertName
	builder := snowflake.NewAlertBuilder(name, database, schemaName)

	enabled := d.Get("enabled").(bool)
	if d.HasChanges("enabled", "warehouse", "alert_schedule", "condition", "action", "comment") {
		if err := snowflake.WaitSuspendAlert(db, name, database, schemaName); err != nil {
			log.Printf("[WARN] failed to suspend alert %s", name)
		}
	}

	if d.HasChange("warehouse") {
		var q string
		newWarehouse := d.Get("warehouse")
		q = builder.ChangeWarehouse(newWarehouse.(string))

		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating warehouse on alert %v", d.Id())
		}
	}

	if d.HasChange("alert_schedule") {
		_, n := d.GetChange("alert_schedule")
		alertSchedule := n.([]interface{})[0].(map[string]interface{})
		log.Printf("[DEBUG] alertSchedule: %v", alertSchedule)
		log.Printf("[DEBUG] alertSchedule[cron]: %v", alertSchedule["cron"])
		c := alertSchedule["cron"].([]interface{})
		if len(c) > 0 {
			if len(c) > 0 {
				cron := c[0].(map[string]interface{})
				cronExpression := cron["expression"].(string)

				timeZone := ""
				if v, ok := cron["time_zone"]; ok {
					timeZone = v.(string)
				}
				stmt := builder.ChangeAlertCronSchedule(cronExpression, timeZone)
				if err := snowflake.Exec(db, stmt); err != nil {
					return fmt.Errorf("error updating alert cron schedule %v err = %w", name, err)
				}
			}
		} else {
			log.Printf("[DEBUG] alertSchedule[interval]: %v", alertSchedule["interval"])
			interval := alertSchedule["interval"].(int)
			stmt := builder.ChangeAlertIntervalSchedule(interval)
			if err := snowflake.Exec(db, stmt); err != nil {
				return fmt.Errorf("error updating alert interval schedule %v err = %w", name, err)
			}
		}
	}

	if d.HasChange("comment") {
		var q string
		_, newVal := d.GetChange("comment")
		q = builder.ChangeComment(newVal.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating comment on alert %v", d.Id())
		}
	}

	if d.HasChange("condition") {
		newVal := d.Get("condition")
		q := builder.ChangeCondition(newVal.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating condition on alert %v", d.Id())
		}
	}

	if d.HasChange("action") {
		newVal := d.Get("action")
		q := builder.ChangeAction(newVal.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating action on alert %v", d.Id())
		}
	}

	if enabled {
		if err := snowflake.WaitResumeAlert(db, name, database, schemaName); err != nil {
			log.Printf("[WARN] failed to resume alert %s", name)
		}
	} else {
		if err := snowflake.WaitSuspendAlert(db, name, database, schemaName); err != nil {
			log.Printf("[WARN] failed to suspend alert %s", name)
		}
	}
	return ReadAlert(d, meta)
}

// DeleteAlert implements schema.DeleteFunc.
func DeleteAlert(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	alterID, err := alertIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := alterID.DatabaseName
	schemaName := alterID.SchemaName
	name := alterID.AlertName

	q := snowflake.NewAlertBuilder(name, database, schemaName).Drop()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting alert %v err = %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}
