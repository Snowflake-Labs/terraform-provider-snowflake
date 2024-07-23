package gen

import (
	"os"

	objectparametersassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type ResourceParametersAssertionsModel struct {
	Name string
	PreambleModel
}

func (m ResourceParametersAssertionsModel) SomeFunc() {
}

func ModelFromSnowflakeObjectParameters(snowflakeObjectParameters objectparametersassertgen.SnowflakeObjectParameters) ResourceParametersAssertionsModel {
	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceParametersAssertionsModel{
		Name: snowflakeObjectParameters.ObjectName(),
		PreambleModel: PreambleModel{
			PackageName: packageWithGenerateDirective,
		},
	}
}
