package datasources

import (
	"context"
	"database/sql"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

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
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	d.SetId("alerts_read")
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	alertPattern := d.Get("pattern").(string)
	var like sdk.Like
	if alertPattern != "" {
		like = sdk.Like{Pattern: &alertPattern}
	}
	listAlerts, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
		In: &sdk.In{
			Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName),
		},
		Like: &like,
	})
	if err != nil {
		log.Printf("[DEBUG] failed to list alerts in schema (%s)", d.Id())
		d.SetId("")
		return err
	}
	log.Printf("[DEBUG] list alerts: %v", listAlerts)
	alerts := make([]map[string]any, 0, len(listAlerts))
	for _, alert := range listAlerts {
		alertMap := map[string]any{}
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
