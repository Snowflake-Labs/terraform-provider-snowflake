package main

import (
	resourcemodelgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model/gen"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

func main() {
	genhelpers.NewGenerator(
		gen.GetProviderSchemaDetails,
		func() func(genhelpers.ResourceSchemaDetails) resourcemodelgen.ResourceConfigBuilderModel {
			return func(resourceSchemaDetails genhelpers.ResourceSchemaDetails) resourcemodelgen.ResourceConfigBuilderModel {
				details := resourcemodelgen.ModelFromResourceSchemaDetails(resourceSchemaDetails)
				details.AdditionalStandardImports = slices.DeleteFunc(details.AdditionalStandardImports, func(dep string) bool { return dep == "encoding/json" })
				return details
			}
		}(),
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.ResourceSchemaDetails, model resourcemodelgen.ResourceConfigBuilderModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_model" + "_gen.go"
}
