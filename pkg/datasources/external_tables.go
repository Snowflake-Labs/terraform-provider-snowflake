package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var externalTablesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the external tables from.",
	},
	"external_tables": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The external tables in the schema",
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

func ExternalTables() *schema.Resource {
	return &schema.Resource{
		Read:   ReadExternalTables,
		Schema: externalTablesSchema,
	}
}

func ReadExternalTables(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentExternalTables, err := snowflake.ListExternalTables(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] external tables in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse external tables in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	externalTables := []map[string]interface{}{}

	for _, externalTable := range currentExternalTables {
		externalTableMap := map[string]interface{}{}

		externalTableMap["name"] = externalTable.ExternalTableName.String
		externalTableMap["database"] = externalTable.DatabaseName.String
		externalTableMap["schema"] = externalTable.SchemaName.String
		externalTableMap["comment"] = externalTable.Comment.String

		externalTables = append(externalTables, externalTableMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("external_tables", externalTables)
}
