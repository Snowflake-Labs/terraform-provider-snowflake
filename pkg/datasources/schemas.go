package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var schemasSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schemas": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The schemas in the database",
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
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Schemas() *schema.Resource {
	return &schema.Resource{
		Read:   ReadSchemas,
		Schema: schemasSchema,
	}
}

func ReadSchemas(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	databaseName := d.Get("database").(string)
	databaseID := sdk.NewAccountObjectIdentifier(databaseName)

	currentSchemas, err := client.Schemas.Show(ctx, &sdk.ShowSchemaOptions{
		In: &sdk.SchemaIn{
			Database: sdk.Bool(true),
			Name:     databaseID,
		},
	})
	if err != nil {
		log.Printf("[DEBUG] unable to show schemas in database (%s)", databaseName)
		d.SetId("")
		return nil
	}

	schemas := make([]map[string]any, len(currentSchemas))
	for i, cs := range currentSchemas {
		schemas[i] = map[string]any{
			"name":     cs.Name,
			"database": cs.DatabaseName,
			"comment":  cs.Comment,
		}
	}

	d.SetId(databaseName)
	return d.Set("schemas", schemas)
}
