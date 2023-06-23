package datasources

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var alertsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The database from which to return the alerts from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The schema from which to return the alerts from.",
	},
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the command output by object name.",
	},
	"alerts": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Lists alerts for the current/specified database or schema, or across the entire account.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the alert.",
				},
				"database_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Database in which the alert is stored.",
				},
				"schema_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Schema in which the alert is stored.",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Comment for the alert.",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Role that owns the alert (i.e. has the OWNERSHIP privilege on the alert)",
				},
				"condition": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The text of the SQL statement that serves as the condition when the alert should be triggered.",
				},
				"action": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The text of the SQL statement that should be executed when the alert is triggered.",
				},
			},
		},
	},
}

// Alerts Snowflake Roles resource.
func Alerts() *schema.Resource {
	return &schema.Resource{
		Read:   ReadAlerts,
		Schema: alertsSchema,
	}
}

// ReadAlerts Reads the database metadata information.
func ReadAlerts(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	d.SetId("alerts_read")
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	alertPattern := d.Get("pattern").(string)

	listAlerts, err := snowflake.ListAlerts(databaseName, schemaName, alertPattern, db)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("[DEBUG] no alerts found in account (%s)", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Println("[DEBUG] failed to list alerts")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] list alerts: %v", listAlerts)

	alerts := []map[string]interface{}{}
	for _, alert := range listAlerts {
		alertMap := map[string]interface{}{}
		alertMap["name"] = alert.Name
		alertMap["comment"] = alert.Comment
		alertMap["owner"] = alert.Owner
		alerts = append(alerts, alertMap)
	}

	if err := d.Set("alerts", alerts); err != nil {
		return err
	}
	return nil
}
