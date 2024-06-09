package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ParameterSchema represents Snowflake parameter object.
// TODO: should be generated later based on the sdk.Parameter
var ParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"value": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"level": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"description": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ParameterToSchema(parameter *sdk.Parameter) map[string]any {
	parameterSchema := make(map[string]any)
	parameterSchema["key"] = parameter.Key
	parameterSchema["value"] = parameter.Value
	parameterSchema["default"] = parameter.Default
	parameterSchema["level"] = parameter.Level
	parameterSchema["description"] = parameter.Description
	return parameterSchema
}
