package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeMaskingPolicySchema represents output of DESCRIBE query for the single masking policy.
var DescribeMaskingPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"signature": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
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

func MaskingPolicyDescriptionToSchema(details sdk.MaskingPolicyDetails) map[string]any {
	signatureElem := make([]map[string]any, len(details.Signature))
	for i, v := range details.Signature {
		signatureElem[i] = map[string]any{
			"name": v.Name,
			"type": string(v.Type),
		}
	}
	return map[string]any{
		"name":        details.Name,
		"signature":   signatureElem,
		"return_type": details.ReturnType,
		"body":        details.Body,
	}
}

// TODO(this pr) merge with row access policy arguments?
func MaskingPolicyArgumentsToSchema(args []sdk.TableColumnSignature) []map[string]any {
	schema := make([]map[string]any, len(args))
	for i, v := range args {
		schema[i] = map[string]any{
			"name": v.Name,
			"type": v.Type,
		}
	}
	return schema
}
