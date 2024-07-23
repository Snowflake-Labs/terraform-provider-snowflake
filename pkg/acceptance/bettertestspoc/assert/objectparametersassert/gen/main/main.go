//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

func main() {
	gencommons.NewGenerator(
		gen.GetAllSnowflakeObjectParameters,
		gen.ModelFromSnowflakeObjectParameters,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ gen.SnowflakeObjectParameters, model gen.SnowflakeObjectParametersAssertionsModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_parameters_snowflake" + "_gen.go"
}
