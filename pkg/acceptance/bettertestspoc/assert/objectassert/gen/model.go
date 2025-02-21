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
	IsOriginalTypeSlice   bool
	Mapper                genhelpers.Mapper
}

func ModelFromSdkObjectDetails(sdkObject genhelpers.SdkObjectDetails) SnowflakeObjectAssertionsModel {
	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	fields := make([]SnowflakeObjectFieldAssertion, len(sdkObject.Fields))
	containsSliceField := false
	for idx, field := range sdkObject.Fields {
		fields[idx] = MapToSnowflakeObjectFieldAssertion(field)
		if !containsSliceField && field.IsSlice() {
			containsSliceField = true
		}
	}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	additionalImports := genhelpers.AdditionalStandardImports(sdkObject.Fields)
	if containsSliceField {
		additionalImports = append(additionalImports, "slices", "errors")
	}
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
		IsOriginalTypeSlice:   field.IsSlice(),
		Mapper:                mapper,
	}
}
