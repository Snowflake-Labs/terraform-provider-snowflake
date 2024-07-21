//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func main() {
	gencommons.NewGenerator(
		getAllSnowflakeObjectParameters,
		gen.ModelFromSnowflakeObjectParameters,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getAllSnowflakeObjectParameters() []gen.SnowflakeObjectParameters {
	return allObjectsParameters
}

func getFilename(_ gen.SnowflakeObjectParameters, model gen.SnowflakeObjectParametersAssertionsModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_parameters_snowflake" + "_gen.go"
}

// TODO: use SDK definition after parameters rework (+ preprocessing here)
var allObjectsParameters = []gen.SnowflakeObjectParameters{
	{
		Name:   "User",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeUser,
		Parameters: []gen.SnowflakeParameter{
			{string(sdk.UserParameterEnableUnredactedQuerySyntaxError), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterNetworkPolicy), "string", "", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterPreventUnloadToInternalStages), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
}
