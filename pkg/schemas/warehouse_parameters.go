package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowWarehouseParametersSchema contains all Snowflake parameters for the warehouses.
// TODO [SNOW-1473425]: descriptions (take from .Description; tool to validate changes later)
// TODO [SNOW-1473425]: should be generated later based on sdk.WarehouseParameters
var ShowWarehouseParametersSchema = map[string]*schema.Schema{
	"max_concurrency_level": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ShowParameterSchema,
		},
	},
	"statement_queued_timeout_in_seconds": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ShowParameterSchema,
		},
	},
	"statement_timeout_in_seconds": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: ShowParameterSchema,
		},
	},
	//"max_concurrency_level":               ParameterSchema,
	//"statement_queued_timeout_in_seconds": ParameterSchema,
	//"statement_timeout_in_seconds":        ParameterSchema,
}

// TODO [SNOW-1473425]: validate all present?
func WarehouseParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	warehouseParameters := make(map[string]any)
	for _, param := range parameters {
		parameterSchema := ParameterToSchema(param)
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.ObjectParameterMaxConcurrencyLevel),
			string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds),
			string(sdk.ObjectParameterStatementTimeoutInSeconds):
			warehouseParameters[strings.ToLower(key)] = []map[string]any{parameterSchema}
		}
	}
	return warehouseParameters
}
