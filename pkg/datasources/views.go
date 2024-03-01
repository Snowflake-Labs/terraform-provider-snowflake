package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var viewsSchema = map[string]*schema.Schema{
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
	"views": {
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

func Views() *schema.Resource {
	return &schema.Resource{
		Read:   ReadViews,
		Schema: viewsSchema,
	}
}

func ReadViews(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	schemaId := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)
	extractedViews, err := client.Views.Show(ctx, sdk.NewShowViewRequest().WithIn(
		&sdk.In{Schema: schemaId},
	))
	if err != nil {
		log.Printf("[DEBUG] failed when searching views in schema (%s), err = %s", schemaId.FullyQualifiedName(), err.Error())
		d.SetId("")
		return nil
	}

	views := make([]map[string]any, len(extractedViews))

	for i, view := range extractedViews {
		views[i] = map[string]any{
			"name":     view.Name,
			"database": view.DatabaseName,
			"schema":   view.SchemaName,
			"comment":  view.Comment,
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(databaseName, schemaName))
	return d.Set("views", views)
}
