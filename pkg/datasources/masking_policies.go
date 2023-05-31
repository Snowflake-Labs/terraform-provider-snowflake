package datasources

import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	maskingPolicies, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
		In: &sdk.In{
			Schema: sdk.NewSchemaIdentifier(databaseName, schemaName),
		},
	})
	if err != nil {
		return err
	}
	maskingPoliciesList := []map[string]interface{}{}
	for _, maskingPolicy := range maskingPolicies {
		maskingPolicyMap := map[string]interface{}{}
		maskingPolicyMap["name"] = maskingPolicy.Name
		maskingPolicyMap["database"] = maskingPolicy.DatabaseName
		maskingPolicyMap["schema"] = maskingPolicy.SchemaName
		maskingPolicyMap["comment"] = maskingPolicy.Comment
		maskingPolicyMap["kind"] = maskingPolicy.Kind
		maskingPoliciesList = append(maskingPoliciesList, maskingPolicyMap)
	}
	if err := d.Set("masking_policies", maskingPoliciesList); err != nil {
		return err
	}
	d.SetId(helpers.EncodeSnowflakeID(databaseName, schemaName))
	return nil
}
