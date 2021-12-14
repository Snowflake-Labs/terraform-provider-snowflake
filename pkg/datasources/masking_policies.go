package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var maskingPoliciesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the maskingPolicies from.",
	},
	"masking_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The maskingPolicies in the schema",
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
				"kind": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func MaskingPolicies() *schema.Resource {
	return &schema.Resource{
		Read:   ReadMaskingPolicies,
		Schema: maskingPoliciesSchema,
	}
}

func ReadMaskingPolicies(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentMaskingPolicies, err := snowflake.ListMaskingPolicies(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] masking policies in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse masking policies in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	maskingPolicies := []map[string]interface{}{}

	for _, maskingPolicy := range currentMaskingPolicies {
		maskingPolicyMap := map[string]interface{}{}

		maskingPolicyMap["name"] = maskingPolicy.Name.String
		maskingPolicyMap["database"] = maskingPolicy.DatabaseName.String
		maskingPolicyMap["schema"] = maskingPolicy.SchemaName.String
		maskingPolicyMap["comment"] = maskingPolicy.Comment.String
		maskingPolicyMap["kind"] = maskingPolicy.Kind.String

		maskingPolicies = append(maskingPolicies, maskingPolicyMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("masking_policies", maskingPolicies)
}
