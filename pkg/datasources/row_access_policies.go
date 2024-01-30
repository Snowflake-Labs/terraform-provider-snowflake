package datasources

import (
	"context"
	"database/sql"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
		Description: "The schema from which to return the row access policy from.",
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
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	schemaId := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)
	extractedRowAccessPolicies, err := client.RowAccessPolicies.Show(ctx, sdk.NewShowRowAccessPolicyRequest().WithIn(
		&sdk.In{Schema: schemaId},
	))
	if err != nil {
		log.Printf("[DEBUG] failed when searching row access policies in schema (%s), err = %s", schemaId.FullyQualifiedName(), err.Error())
		d.SetId("")
		return nil
	}

	rowAccessPolicies := make([]map[string]any, len(extractedRowAccessPolicies))

	for i, rowAccessPolicy := range extractedRowAccessPolicies {
		rowAccessPolicies[i] = map[string]any{
			"name":     rowAccessPolicy.Name,
			"database": rowAccessPolicy.DatabaseName,
			"schema":   rowAccessPolicy.SchemaName,
			"comment":  rowAccessPolicy.Comment,
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(databaseName, schemaName))
	return d.Set("row_access_policies", rowAccessPolicies)
}
