package datasources

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	databaseName := d.Get("database").(string)
	databaseID := sdk.NewAccountObjectIdentifier(databaseName)

	log.Printf("[DEBUG] database name %s", databaseName)

	currentSchemas, err := client.Schemas.Show(ctx, &sdk.ShowSchemaOptions{
		In: &sdk.SchemaIn{
			Database: sdk.Bool(true),
			Name:     databaseID,
		},
	})

	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] schemas in database (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse schemas in database (%s)", d.Id())
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
