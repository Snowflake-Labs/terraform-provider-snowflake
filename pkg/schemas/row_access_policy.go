package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var RowAccessPolicyDescribeSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"signature": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"return_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"body": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func RowAccessPolicyDescriptionToSchema(description sdk.RowAccessPolicyDescription) map[string]any {
	return map[string]any{
		"name":        description.Name,
		"signature":   description.Signature,
		"return_type": description.ReturnType,
		"body":        description.Body,
	}
}

func RowAccessPolicyArgumentsToSchema(args []sdk.RowAccessPolicyArgument) []map[string]any {
	schema := make([]map[string]any, len(args))
	for i, v := range args {
		schema[i] = map[string]any{
			"name": v.Name,
			"type": v.Type,
		}
	}
	return schema
}
