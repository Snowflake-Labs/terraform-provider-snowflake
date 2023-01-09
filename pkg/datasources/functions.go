package datasources

import (
	"database/sql"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var functionsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the functions from.",
	},
	"functions": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The functions in the schema",
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
				"argument_types": {
					Type:     schema.TypeList,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
					Computed: true,
				},
				"return_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Functions() *schema.Resource {
	return &schema.Resource{
		Read:   ReadFunctions,
		Schema: functionsSchema,
	}
}

// todo: fix this. ListUserFunctions isn't using the right struct right now and also the signature of this doesn't support all the features it could for example, database and schema should be optional, and you could also list by account.
func ReadFunctions(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	d.SetId("functions")
	currentFunctions, err := snowflake.ListUserFunctions(databaseName, schemaName, db)
	if err != nil {
		log.Printf("[DEBUG] error listing functions: %v", err)
		return nil
	}

	functions := []map[string]interface{}{}

	for _, function := range currentFunctions {
		functionMap := map[string]interface{}{}

		functionSignatureMap, err := parseArguments(function.Arguments.String)
		if err != nil {
			return err
		}

		functionMap["name"] = function.Name.String
		functionMap["database"] = databaseName
		functionMap["schema"] = schemaName
		functionMap["comment"] = function.Description.String
		functionMap["argument_types"] = functionSignatureMap["argumentTypes"].([]string)
		functionMap["return_type"] = functionSignatureMap["returnType"].(string)

		functions = append(functions, functionMap)
	}

	return d.Set("functions", functions)
}
