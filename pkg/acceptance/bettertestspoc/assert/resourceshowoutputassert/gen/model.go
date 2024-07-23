package gen

import (
	"os"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type ResourceShowOutputAssertionsModel struct {
	Name    string
	SdkType string
	IdType  string
	Fields  []ResourceShowOutputAssertionModel
	PreambleModel
}

func (m ResourceShowOutputAssertionsModel) SomeFunc() {
}

type ResourceShowOutputAssertionModel struct {
	Name                  string
	ConcreteType          string
	IsOriginalTypePointer bool
	Mapper                gencommons.Mapper
}

func ModelFromSdkObjectDetails(sdkObject gencommons.SdkObjectDetails) ResourceShowOutputAssertionsModel {
	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	fields := make([]ResourceShowOutputAssertionModel, len(sdkObject.Fields))
	imports := make(map[string]struct{})
	for idx, field := range sdkObject.Fields {
		fields[idx] = MapToSnowflakeObjectFieldAssertion(field)
		additionalImport, isImportedType := field.GetImportedType()
		if isImportedType {
			imports[additionalImport] = struct{}{}
		}
	}
	additionalImports := make([]string, 0)
	for k := range imports {
		if !slices.Contains([]string{"sdk"}, k) {
			additionalImports = append(additionalImports, k)
		}
	}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceShowOutputAssertionsModel{
		Name:    name,
		SdkType: sdkObject.Name,
		IdType:  sdkObject.IdType,
		Fields:  fields,
		PreambleModel: PreambleModel{
			PackageName:               packageWithGenerateDirective,
			AdditionalStandardImports: additionalImports,
		},
	}
}

func MapToSnowflakeObjectFieldAssertion(field gencommons.Field) ResourceShowOutputAssertionModel {
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")

	// TODO [SNOW-1501905]: handle other mappings if needed
	mapper := gencommons.Identity
	if concreteTypeWithoutPtr == "sdk.AccountObjectIdentifier" {
		mapper = gencommons.Name
	}

	return ResourceShowOutputAssertionModel{
		Name:                  field.Name,
		ConcreteType:          field.ConcreteType,
		IsOriginalTypePointer: field.IsPointer(),
		Mapper:                mapper,
	}
}
