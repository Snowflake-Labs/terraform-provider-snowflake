package datasources

import (
	"database/sql"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
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
	databaseName := d.Get("database").(string)

	log.Printf("[DEBUG] database name %s", databaseName)

	currentSchemas, err := snowflake.ListSchemas(databaseName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] schemas in database (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse schemas in database (%s)", d.Id())
		d.SetId("")
		return nil
	}

	schemas := []map[string]interface{}{}

	for _, schema := range currentSchemas {
		schemaMap := map[string]interface{}{}

		schemaMap["name"] = schema.Name.String
		schemaMap["database"] = schema.DatabaseName.String
		schemaMap["comment"] = schema.Comment.String

		schemas = append(schemas, schemaMap)
	}

	d.SetId(databaseName)
	return d.Set("schemas", schemas)
}
