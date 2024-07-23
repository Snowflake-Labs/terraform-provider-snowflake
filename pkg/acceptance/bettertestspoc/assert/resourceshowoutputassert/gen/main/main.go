//go:build exclude

package main

import (
	objectassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

func main() {
	gencommons.NewGenerator(
		objectassertgen.GetSdkObjectDetails,
		gen.ModelFromSdkObjectDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ gencommons.SdkObjectDetails, model gen.ResourceShowOutputAssertionsModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_show_output" + "_gen.go"
}
