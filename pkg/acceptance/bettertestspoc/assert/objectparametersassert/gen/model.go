package gen

import (
	"os"
	"strings"
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
	Name             string
	Type             string
	DefaultLevel     string
	AssertionCreator string
}

func (m SnowflakeObjectParametersAssertionsModel) SomeFunc() {
	return
}

func ModelFromSnowflakeObjectParameters(snowflakeObjectParameters SnowflakeObjectParameters) SnowflakeObjectParametersAssertionsModel {
	parameters := make([]ParameterAssertionModel, len(snowflakeObjectParameters.Parameters))
	for idx, p := range snowflakeObjectParameters.Parameters {
		// TODO: get a runtime name for the assertion creator
		var assertionCreator string
		switch {
		case p.ParameterType == "bool":
			assertionCreator = "SnowflakeParameterBoolValueSet"
		case p.ParameterType == "int":
			assertionCreator = "SnowflakeParameterIntValueSet"
		case p.ParameterType == "string":
			assertionCreator = "SnowflakeParameterValueSet"
		case strings.HasPrefix(p.ParameterType, "sdk."):
			assertionCreator = "SnowflakeParameterStringUnderlyingValueSet"
		default:
			assertionCreator = "SnowflakeParameterValueSet"
		}

		parameters[idx] = ParameterAssertionModel{
			Name:             p.ParameterName,
			Type:             p.ParameterType,
			DefaultLevel:     p.DefaultLevel,
			AssertionCreator: assertionCreator,
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
