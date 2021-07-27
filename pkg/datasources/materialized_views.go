package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var materializedViewsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the views from.",
	},
	"materialized_views": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The views in the schema",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"database": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"schema": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func MaterializedViews() *schema.Resource {
	return &schema.Resource{
		Read:   ReadMaterializedViews,
		Schema: materializedViewsSchema,
	}
}

func ReadMaterializedViews(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentViews, err := snowflake.ListMaterializedViews(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] materialized views in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] materialized unable to parse views in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	views := []map[string]interface{}{}

	for _, view := range currentViews {
		viewMap := map[string]interface{}{}

		viewMap["name"] = view.Name.String
		viewMap["database"] = view.DatabaseName.String
		viewMap["schema"] = view.SchemaName.String
		viewMap["comment"] = view.Comment.String

		views = append(views, viewMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("materialized_views", views)
}
