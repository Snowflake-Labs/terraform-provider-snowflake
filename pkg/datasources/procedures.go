package datasources

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		ReadContext: ReadContextProcedures,
		Schema:      proceduresSchema,
	}
}

func ReadContextProcedures(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	req := sdk.NewShowProcedureRequest()
	if databaseName != "" {
		req.WithIn(&sdk.In{Database: sdk.NewAccountObjectIdentifier(databaseName)})
	}
	if schemaName != "" {
		req.WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)})
	}
	procedures, err := client.Procedures.Show(ctx, req)
	if err != nil {
		id := fmt.Sprintf(`%v|%v`, databaseName, schemaName)

		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Unable to parse procedures in schema (%s)", id),
				Detail:   "See our document on design decisions for procedures: <LINK (coming soon)>",
			},
		}
	}
	proceduresList := []map[string]interface{}{}

	for _, procedure := range procedures {
		procedureMap := map[string]interface{}{}
		procedureMap["name"] = procedure.Name
		procedureMap["database"] = procedure.CatalogName
		procedureMap["schema"] = procedure.SchemaName
		procedureMap["comment"] = procedure.Description
		procedureSignatureMap, err := parseArguments(procedure.Arguments)
		if err != nil {
			return diag.FromErr(err)
		}
		procedureMap["argument_types"] = procedureSignatureMap["argumentTypes"].([]string)
		procedureMap["return_type"] = procedureSignatureMap["returnType"].(string)
		proceduresList = append(proceduresList, procedureMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	if err := d.Set("procedures", proceduresList); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func parseArguments(arguments string) (map[string]interface{}, error) {
	r := regexp.MustCompile(`(?P<callable_name>[^(]+)\((?P<argument_signature>[^)]*)\) RETURN (?P<return_type>.*)`)
	matches := r.FindStringSubmatch(arguments)
	if len(matches) == 0 {
		return nil, fmt.Errorf(`could not parse arguments: %v`, arguments)
	}
	callableSignatureMap := make(map[string]interface{})

	argumentTypes := strings.Split(matches[2], ", ")

	callableSignatureMap["callableName"] = matches[1]
	callableSignatureMap["argumentTypes"] = argumentTypes
	callableSignatureMap["returnType"] = matches[3]

	return callableSignatureMap, nil
}
