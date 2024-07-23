//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

func main() {
	gencommons.NewGenerator(
		gen.GetSdkObjectDetails,
		gen.ModelFromSdkObjectDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ gencommons.SdkObjectDetails, model gen.SnowflakeObjectAssertionsModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_snowflake" + "_gen.go"
}
