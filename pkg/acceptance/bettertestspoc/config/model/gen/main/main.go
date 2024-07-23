package main

import (
	resourceassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

func main() {
	gencommons.NewGenerator(
		resourceassertgen.GetResourceSchemaDetails,
		gen.ModelFromResourceSchemaDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ gencommons.ResourceSchemaDetails, model gen.ResourceConfigBuilderModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_model" + "_gen.go"
}
