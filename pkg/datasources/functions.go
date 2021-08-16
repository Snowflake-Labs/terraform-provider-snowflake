package datasources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
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

func ReadFunctions(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentFunctions, err := snowflake.ListFunctions(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] functions in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse functions in schema (%s)", d.Id())
		d.SetId("")
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
		functionMap["database"] = function.DatabaseName.String
		functionMap["schema"] = function.SchemaName.String
		functionMap["comment"] = function.Comment.String
		functionMap["argument_types"] = functionSignatureMap["argumentTypes"].([]string)
		functionMap["return_type"] = functionSignatureMap["returnType"].(string)

		functions = append(functions, functionMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("functions", functions)
}

func parseArguments(arguments string) (map[string]interface{}, error) {
	r := regexp.MustCompile(`(?P<callable_name>[^(]+)\((?P<argument_signature>[^)]*)\) RETURN (?P<return_type>.*)`)
	matches := r.FindStringSubmatch(arguments)
	if len(matches) == 0 {
		return nil, errors.New(fmt.Sprintf(`Could not parse arguments: %v`, arguments))
	}
	callableSignatureMap := make(map[string]interface{})

	argumentTypes := strings.Split(matches[2], ", ")

	callableSignatureMap["callableName"] = matches[1]
	callableSignatureMap["argumentTypes"] = argumentTypes
	callableSignatureMap["returnType"] = matches[3]

	return callableSignatureMap, nil
}
