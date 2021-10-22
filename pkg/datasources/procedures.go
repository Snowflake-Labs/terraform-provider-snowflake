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

var proceduresSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the procedures from.",
	},
	"procedures": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The procedures in the schema",
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

func Procedures() *schema.Resource {
	return &schema.Resource{
		Read:   ReadProcedures,
		Schema: proceduresSchema,
	}
}

func ReadProcedures(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentProcedures, err := snowflake.ListProcedures(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] procedures in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse procedures in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	procedures := []map[string]interface{}{}

	for _, procedure := range currentProcedures {
		procedureMap := map[string]interface{}{}

		procedureSignatureMap, err := parseArguments(procedure.Arguments.String)
		if err != nil {
			return err
		}

		procedureMap["name"] = procedure.Name.String
		procedureMap["database"] = procedure.DatabaseName.String
		procedureMap["schema"] = procedure.SchemaName.String
		procedureMap["comment"] = procedure.Comment.String
		procedureMap["argument_types"] = procedureSignatureMap["argumentTypes"].([]string)
		procedureMap["return_type"] = procedureSignatureMap["returnType"].(string)

		procedures = append(procedures, procedureMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("procedures", procedures)
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
