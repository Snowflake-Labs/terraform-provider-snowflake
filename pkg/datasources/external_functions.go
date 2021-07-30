package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var externalFunctionsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the external functions from.",
	},
	"external_functions": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The external functions in the schema",
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
				"language": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func ExternalFunctions() *schema.Resource {
	return &schema.Resource{
		Read:   ReadExternalFunctions,
		Schema: externalFunctionsSchema,
	}
}

func ReadExternalFunctions(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentExternalFunctions, err := snowflake.ListExternalFunctions(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] external functions in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse external functions in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	externalFunctions := []map[string]interface{}{}

	for _, externalFunction := range currentExternalFunctions {
		externalFunctionMap := map[string]interface{}{}

		externalFunctionMap["name"] = externalFunction.ExternalFunctionName.String
		externalFunctionMap["database"] = externalFunction.DatabaseName.String
		externalFunctionMap["schema"] = externalFunction.SchemaName.String
		externalFunctionMap["comment"] = externalFunction.Comment.String
		externalFunctionMap["language"] = externalFunction.Language.String

		externalFunctions = append(externalFunctions, externalFunctionMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("external_functions", externalFunctions)
}
