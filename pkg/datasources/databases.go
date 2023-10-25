package datasources

import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databasesSchema = map[string]*schema.Schema{
	"terse": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optionally returns only the columns `created_on` and `name` in the results",
	},
	"history": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optionally includes dropped databases that have not yet been purged The output also includes an additional `dropped_on` column",
	},
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Optionally filters the databases by a pattern",
	},
	"starts_with": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Optionally filters the databases by a pattern",
	},
	"databases": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Snowflake databases",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"owner": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"is_default": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"is_current": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"origin": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"retention_time": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"created_on": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"options": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"replication_configuration": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"accounts": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"ignore_edition_check": {
								Type:     schema.TypeBool,
								Computed: true,
							},
						},
					},
				},
			},
		},
	},
}

// Databases the Snowflake current account resource.
func Databases() *schema.Resource {
	return &schema.Resource{
		Read:   ReadDatabases,
		Schema: databasesSchema,
	}
}

// ReadDatabases read the current snowflake account information.
func ReadDatabases(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	opts := sdk.ShowDatabasesOptions{}
	if terse, ok := d.GetOk("terse"); ok {
		opts.Terse = sdk.Bool(terse.(bool))
	}
	if history, ok := d.GetOk("history"); ok {
		opts.History = sdk.Bool(history.(bool))
	}
	if pattern, ok := d.GetOk("pattern"); ok {
		opts.Like = &sdk.Like{
			Pattern: sdk.String(pattern.(string)),
		}
	}
	if startsWith, ok := d.GetOk("starts_with"); ok {
		opts.StartsWith = sdk.String(startsWith.(string))
	}
	databases, err := client.Databases.Show(ctx, &opts)
	if err != nil {
		return err
	}
	d.SetId("databases_read")
	flattenedDatabases := []map[string]interface{}{}
	for _, database := range databases {
		flattenedDatabase := map[string]interface{}{}
		flattenedDatabase["name"] = database.Name
		flattenedDatabase["comment"] = database.Comment
		flattenedDatabase["owner"] = database.Owner
		flattenedDatabase["is_default"] = database.IsDefault
		flattenedDatabase["is_current"] = database.IsCurrent
		flattenedDatabase["origin"] = database.Origin
		flattenedDatabase["created_on"] = database.CreatedOn.String()
		flattenedDatabase["options"] = database.Options
		flattenedDatabase["retention_time"] = database.RetentionTime
		flattenedDatabases = append(flattenedDatabases, flattenedDatabase)
	}
	err = d.Set("databases", flattenedDatabases)
	if err != nil {
		return err
	}
	return nil
}
