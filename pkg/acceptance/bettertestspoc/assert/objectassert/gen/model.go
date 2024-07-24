package gen

import (
	"os"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type SnowflakeObjectAssertionsModel struct {
	Name    string
	SdkType string
	IdType  string
	Fields  []SnowflakeObjectFieldAssertion
	PreambleModel
}

func (m SnowflakeObjectAssertionsModel) SomeFunc() {
}

type SnowflakeObjectFieldAssertion struct {
	Name                  string
	ConcreteType          string
	IsOriginalTypePointer bool
	Mapper                genhelpers.Mapper
}

func ModelFromSdkObjectDetails(sdkObject genhelpers.SdkObjectDetails) SnowflakeObjectAssertionsModel {
	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	fields := make([]SnowflakeObjectFieldAssertion, len(sdkObject.Fields))
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
	return SnowflakeObjectAssertionsModel{
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

func MapToSnowflakeObjectFieldAssertion(field genhelpers.Field) SnowflakeObjectFieldAssertion {
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")

	// TODO [SNOW-1501905]: handle other mappings if needed
	mapper := genhelpers.Identity
	if concreteTypeWithoutPtr == "sdk.AccountObjectIdentifier" {
		mapper = genhelpers.Name
	}

	return SnowflakeObjectFieldAssertion{
		Name:                  field.Name,
		ConcreteType:          field.ConcreteType,
		IsOriginalTypePointer: field.IsPointer(),
		Mapper:                mapper,
	}
}
