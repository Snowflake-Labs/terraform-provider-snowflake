package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowWarehouseParametersSchema contains all Snowflake parameters for the warehouses.
// TODO: descriptions (take from .Description; tool to validate changes later)
// TODO: should be generated later based on the existing Snowflake parameters for warehouses
var ShowWarehouseParametersSchema = map[string]*schema.Schema{
	"max_concurrency_level": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	"statement_queued_timeout_in_seconds": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
	"statement_timeout_in_seconds": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ParameterSchema,
		},
	},
}

// TODO: validate all present?
func WarehouseParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	warehouseParameters := make(map[string]any)
	for _, param := range parameters {
		parameterSchema := ParameterToSchema(param)
		switch strings.ToUpper(param.Key) {
		case string(sdk.ObjectParameterMaxConcurrencyLevel):
			warehouseParameters["max_concurrency_level"] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds):
			warehouseParameters["statement_queued_timeout_in_seconds"] = []map[string]any{parameterSchema}
		case string(sdk.ObjectParameterStatementTimeoutInSeconds):
			warehouseParameters["statement_timeout_in_seconds"] = []map[string]any{parameterSchema}
		}
	}
	return warehouseParameters
}
