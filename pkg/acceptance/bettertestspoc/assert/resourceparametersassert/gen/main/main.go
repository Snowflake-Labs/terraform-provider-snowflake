//go:build exclude

package main

import (
	objectparametersassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

func main() {
	gencommons.NewGenerator(
		objectparametersassertgen.GetAllSnowflakeObjectParameters,
		gen.ModelFromSnowflakeObjectParameters,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ objectparametersassertgen.SnowflakeObjectParameters, model gen.ResourceParametersAssertionsModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_parameters_snowflake" + "_gen.go"
}
