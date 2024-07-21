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
	Name       string
	IdType     string
	Parameters []ParameterAssertionModel
	PreambleModel
}

type ParameterAssertionModel struct {
	Name         string
	DefaultLevel string
}

func (m SnowflakeObjectParametersAssertionsModel) SomeFunc() {
	return
}

func ModelFromSnowflakeObjectParameters(snowflakeObjectParameters SnowflakeObjectParameters) SnowflakeObjectParametersAssertionsModel {
	parameters := make([]ParameterAssertionModel, len(snowflakeObjectParameters.Parameters))
	for idx, p := range snowflakeObjectParameters.Parameters {
		parameters[idx] = ParameterAssertionModel{
			Name:         p.ParameterName,
			DefaultLevel: p.DefaultLevel,
		}
	}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return SnowflakeObjectParametersAssertionsModel{
		Name:       snowflakeObjectParameters.ObjectName(),
		IdType:     snowflakeObjectParameters.IdType,
		Parameters: parameters,
		PreambleModel: PreambleModel{
			PackageName:               packageWithGenerateDirective,
			AdditionalStandardImports: []string{},
		},
	}
}
