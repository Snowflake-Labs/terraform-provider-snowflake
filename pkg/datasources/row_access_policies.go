package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var rowAccessPoliciesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the row access policyfrom.",
	},
	"row_access_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The row access policy in the schema",
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

func RowAccessPolicies() *schema.Resource {
	return &schema.Resource{
		Read:   ReadRowAccessPolicies,
		Schema: rowAccessPoliciesSchema,
	}
}

func ReadRowAccessPolicies(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentRowAccessPolicies, err := snowflake.ListRowAccessPolicies(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] row access policy in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse row access policy in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	rowAccessPolicies := []map[string]interface{}{}

	for _, rowAccessPolicy := range currentRowAccessPolicies {
		rowAccessPolicyMap := map[string]interface{}{}

		rowAccessPolicyMap["name"] = rowAccessPolicy.Name.String
		rowAccessPolicyMap["database"] = rowAccessPolicy.DatabaseName.String
		rowAccessPolicyMap["schema"] = rowAccessPolicy.SchemaName.String
		rowAccessPolicyMap["comment"] = rowAccessPolicy.Comment.String

		rowAccessPolicies = append(rowAccessPolicies, rowAccessPolicyMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("row_access_policies", rowAccessPolicies)
}
