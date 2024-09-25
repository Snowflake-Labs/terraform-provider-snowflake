package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowTaskParametersSchema = make(map[string]*schema.Schema)
)

func init() {
	for _, param := range sdk.AllTaskParameters {
		ShowTaskParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func TaskParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	taskParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(userParameters, sdk.UserParameter(param.Key)) {
			taskParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return taskParametersValue
}
