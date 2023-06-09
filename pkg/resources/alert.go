package resources

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	client := sdk.NewClientFromDB(db)
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	ctx := context.Background()
	alert, err := client.Alerts.ShowByID(ctx, objectIdentifier)

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
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)

	ctx := context.Background()
	objectIdentifier := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	alertSchedule := getAlertSchedule(d)

	warehouseName := d.Get("warehouse").(string)
	warehouse := sdk.NewAccountObjectIdentifier(warehouseName)

	opts := &sdk.CreateAlertOptions{}

	if v, ok := d.GetOk("comment"); ok {
		opts.Comment = sdk.String(v.(string))
	}

	condition := d.Get("condition").(string)

	action := d.Get("action").(string)

	err := client.Alerts.Create(ctx, objectIdentifier, warehouse, alertSchedule, condition, action, opts)

	if err != nil {
		return err
	}

	enabled := d.Get("enabled").(bool)

	if enabled {
		opts := sdk.AlterAlertOptions{State: &sdk.Resume}
		err := client.Alerts.Alter(ctx, objectIdentifier, &opts)
		if err != nil {
			return err
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))

	return ReadAlert(d, meta)
}

func getAlertSchedule(d *schema.ResourceData) sdk.AlertSchedule {
	var alertSchedule sdk.AlertSchedule
	if v, ok := d.GetOk("alert_schedule"); ok {
		schedule := v.([]interface{})[0].(map[string]interface{})
		if v, ok := schedule["cron"]; ok {
			c := v.([]interface{})
			if len(c) > 0 {
				cron := c[0].(map[string]interface{})
				cronExpression := cron["expression"].(string)
				timeZone := cron["time_zone"].(string)
				alertSchedule = sdk.AlertScheduleCronExpression{Expression: cronExpression, TimeZone: timeZone}
			}
		}
		if v, ok := schedule["interval"]; ok {
			interval := v.(int)
			if interval > 0 {
				alertSchedule = sdk.AlertScheduleInterval(interval)
			}
		}
	}
	return alertSchedule
}

// UpdateAlert implements schema.UpdateFunc.
func UpdateAlert(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	ctx := context.Background()

	enabled := d.Get("enabled").(bool)
	if d.HasChanges("enabled", "warehouse", "alert_schedule", "condition", "action", "comment") {
		if err := snowflake.WaitSuspendAlert(ctx, client, objectIdentifier); err != nil {
			log.Printf("[WARN] failed to suspend alert %s", objectIdentifier.Name())
		}
	}

	if d.HasChange("warehouse") {
		alterOptions := &sdk.AlterAlertOptions{}
		warehouseName := d.Get("warehouse").(string)
		warehouse := sdk.NewAccountObjectIdentifier(warehouseName)
		alterOptions.Set = &sdk.AlertSet{
			Warehouse: &warehouse,
		}
		err := client.Alerts.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return fmt.Errorf("error updating warehouse on alert %v", objectIdentifier.Name())
		}
	}

	if d.HasChange("alert_schedule") {
		alertSchedule := getAlertSchedule(d).String()
		alterOptions := &sdk.AlterAlertOptions{}
		alterOptions.Set = &sdk.AlertSet{
			Schedule: &alertSchedule,
		}
		err := client.Alerts.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return fmt.Errorf("error updating schedule on alert %v", objectIdentifier.Name())
		}

	}

	if d.HasChange("comment") {
		newComment := d.Get("comment").(string)
		alterOptions := &sdk.AlterAlertOptions{}
		alterOptions.Set = &sdk.AlertSet{
			Comment: &newComment,
		}
		err := client.Alerts.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return fmt.Errorf("error updating schedule on comment %v", objectIdentifier.Name())
		}
	}

	if d.HasChange("condition") {
		condition := d.Get("condition").(string)
		alterOptions := &sdk.AlterAlertOptions{}
		alterOptions.ModifyCondition = &condition
		err := client.Alerts.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return fmt.Errorf("error updating schedule on condition %v", objectIdentifier.Name())
		}
	}

	if d.HasChange("action") {
		action := d.Get("action").(string)
		alterOptions := &sdk.AlterAlertOptions{}
		alterOptions.ModifyAction = &action
		err := client.Alerts.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return fmt.Errorf("error updating schedule on action %v", objectIdentifier.Name())
		}
	}

	if enabled {
		if err := snowflake.WaitResumeAlert(ctx, client, objectIdentifier); err != nil {
			log.Printf("[WARN] failed to resume alert %s", objectIdentifier.Name())
		}
	} else {
		if err := snowflake.WaitSuspendAlert(ctx, client, objectIdentifier); err != nil {
			log.Printf("[WARN] failed to suspend alert %s", objectIdentifier.Name())
		}
	}
	return ReadAlert(d, meta)
}

// DeleteAlert implements schema.DeleteFunc.
func DeleteAlert(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.Alerts.Drop(ctx, objectIdentifier)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
