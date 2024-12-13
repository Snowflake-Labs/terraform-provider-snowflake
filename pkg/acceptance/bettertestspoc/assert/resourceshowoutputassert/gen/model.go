package gen

import (
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
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
	Mapper           genhelpers.Mapper
}

func ModelFromSdkObjectDetails(sdkObject genhelpers.SdkObjectDetails) ResourceShowOutputAssertionsModel {
	attributes := make([]ResourceShowOutputAssertionModel, len(sdkObject.Fields))
	includeFmt := false
	for idx, field := range sdkObject.Fields {
		showOutputAssertions, inFmt := MapToResourceShowOutputAssertion(field)
		if !includeFmt && inFmt {
			includeFmt = true
		}
		attributes[idx] = showOutputAssertions
	}

	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	imports := genhelpers.AdditionalStandardImports(sdkObject.Fields)
	if includeFmt {
		imports = append(imports, "fmt")
	}
	return ResourceShowOutputAssertionsModel{
		Name:       name,
		Attributes: attributes,
		PreambleModel: PreambleModel{
			PackageName:               packageWithGenerateDirective,
			AdditionalStandardImports: imports,
		},
	}
}

func MapToResourceShowOutputAssertion(field genhelpers.Field) (ResourceShowOutputAssertionModel, bool) { // TODO: Temporary
	isPrimitive := true
	includeFmt := false

	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	// TODO [SNOW-1501905]: get a runtime name for the assertion creator
	var assertionCreator string
	switch {
	case concreteTypeWithoutPtr == "bool":
		assertionCreator = "ResourceShowOutputBoolValue"
	case concreteTypeWithoutPtr == "int":
		assertionCreator = "ResourceShowOutputIntValue"
	case concreteTypeWithoutPtr == "float64":
		assertionCreator = "ResourceShowOutputFloatValue"
	case concreteTypeWithoutPtr == "string":
		assertionCreator = "ResourceShowOutputValue"
	// TODO [SNOW-1501905]: distinguish between different enum types
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk."):
		assertionCreator = "ResourceShowOutputStringUnderlyingValue"
	default:
		assertionCreator = "ResourceShowOutputValue"
		isPrimitive = false
	}

	// TODO [SNOW-1501905]: handle other mappings if needed
	mapper := genhelpers.Identity
	switch concreteTypeWithoutPtr {
	case "sdk.AccountObjectIdentifier":
		mapper = genhelpers.Name
	case "sdk.AccountIdentifier", "sdk.DatabaseObjectIdentifier", "sdk.SchemaObjectIdentifier", "sdk.SchemaObjectIdentifierWithArguments", "sdk.ExternalObjectIdentifier":
		mapper = genhelpers.FullyQualifiedName
	case "time.Time":
		mapper = genhelpers.ToString
	default:
		if !isPrimitive {
			mapper = genhelpers.PrintToString
			includeFmt = true
		}
	}

	return ResourceShowOutputAssertionModel{
		Name:             field.Name,
		ConcreteType:     concreteTypeWithoutPtr,
		AssertionCreator: assertionCreator,
		Mapper:           mapper,
	}, includeFmt
}
