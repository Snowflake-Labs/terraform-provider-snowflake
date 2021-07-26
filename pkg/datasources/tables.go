package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var tablesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the tables from.",
	},
	"tables": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The tables in the schema",
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

func Tables() *schema.Resource {
	return &schema.Resource{
		Read:   ReadTables,
		Schema: tablesSchema,
	}
}

func ReadTables(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentTables, err := snowflake.ListTables(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] tables in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse tables in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	tables := []map[string]interface{}{}

	for _, table := range currentTables {
		tableMap := map[string]interface{}{}

		tableMap["name"] = table.TableName.String
		tableMap["database"] = table.DatabaseName.String
		tableMap["schema"] = table.SchemaName.String
		tableMap["comment"] = table.Comment.String

		tables = append(tables, tableMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("tables", tables)
}
