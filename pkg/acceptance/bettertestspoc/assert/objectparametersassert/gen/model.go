package gen

import (
	"os"
)

// TODO: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type SnowflakeObjectParametersAssertionsModel struct {
	Name string
	PreambleModel
}

func (m SnowflakeObjectParametersAssertionsModel) SomeFunc() {
	return
}

func ModelFromSnowflakeObjectParameters(snowflakeObjectParameters SnowflakeObjectParameters) SnowflakeObjectParametersAssertionsModel {
	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return SnowflakeObjectParametersAssertionsModel{
		Name: snowflakeObjectParameters.ObjectName(),
		PreambleModel: PreambleModel{
			PackageName:               packageWithGenerateDirective,
			AdditionalStandardImports: []string{},
		},
	}
}
