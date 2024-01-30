package datasources

import (
	"context"
	"database/sql"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	schemaId := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)
	extractedMaterializedViews, err := client.MaterializedViews.Show(ctx, sdk.NewShowMaterializedViewRequest().WithIn(
		&sdk.In{Schema: schemaId},
	))
	if err != nil {
		log.Printf("[DEBUG] failed when searching materialized views in schema (%s), err = %s", schemaId.FullyQualifiedName(), err.Error())
		d.SetId("")
		return nil
	}

	materializedViews := make([]map[string]any, len(extractedMaterializedViews))

	for i, materializedView := range extractedMaterializedViews {
		materializedViews[i] = map[string]any{
			"name":     materializedView.Name,
			"database": materializedView.DatabaseName,
			"schema":   materializedView.SchemaName,
			"comment":  materializedView.Comment,
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(databaseName, schemaName))
	return d.Set("materialized_views", materializedViews)
}
