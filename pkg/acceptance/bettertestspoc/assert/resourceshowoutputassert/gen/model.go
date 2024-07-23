package gen

import (
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type ResourceShowOutputAssertionsModel struct {
	Name       string
	Attributes []ResourceShowOutputAssertionModel
	PreambleModel
}

func (m ResourceShowOutputAssertionsModel) SomeFunc() {
}

type ResourceShowOutputAssertionModel struct {
	Name             string
	ConcreteType     string
	AssertionCreator string
	Mapper           gencommons.Mapper
}

func ModelFromSdkObjectDetails(sdkObject gencommons.SdkObjectDetails) ResourceShowOutputAssertionsModel {
	attributes := make([]ResourceShowOutputAssertionModel, len(sdkObject.Fields))
	for idx, field := range sdkObject.Fields {
		attributes[idx] = MapToResourceShowOutputAssertion(field)
	}

	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceShowOutputAssertionsModel{
		Name:       name,
		Attributes: attributes,
		PreambleModel: PreambleModel{
			PackageName:               packageWithGenerateDirective,
			AdditionalStandardImports: gencommons.AdditionalStandardImports(sdkObject.Fields),
		},
	}
}

func MapToResourceShowOutputAssertion(field gencommons.Field) ResourceShowOutputAssertionModel {
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	// TODO [SNOW-1501905]: get a runtime name for the assertion creator
	var assertionCreator string
	switch {
	case concreteTypeWithoutPtr == "bool":
		assertionCreator = "ResourceShowOutputBoolValueSet"
	case concreteTypeWithoutPtr == "int":
		assertionCreator = "ResourceShowOutputIntValueSet"
	case concreteTypeWithoutPtr == "float64":
		assertionCreator = "ResourceShowOutputFloatValueSet"
	case concreteTypeWithoutPtr == "string":
		assertionCreator = "ResourceShowOutputValueSet"
	// TODO: distinguish between different enum types
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk."):
		assertionCreator = "ResourceShowOutputStringUnderlyingValueSet"
	default:
		assertionCreator = "ResourceShowOutputValueSet"
	}

	// TODO [SNOW-1501905]: handle other mappings if needed
	mapper := gencommons.Identity
	switch concreteTypeWithoutPtr {
	case "sdk.AccountObjectIdentifier":
		mapper = gencommons.Name
	case "time.Time":
		mapper = gencommons.ToString
	}

	return ResourceShowOutputAssertionModel{
		Name:             field.Name,
		ConcreteType:     concreteTypeWithoutPtr,
		AssertionCreator: assertionCreator,
		Mapper:           mapper,
	}
}
