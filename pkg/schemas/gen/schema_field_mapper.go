package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SchemaField struct {
	Name                  string
	SchemaType            schema.ValueType
	OriginalName          string
	IsOriginalTypePointer bool
	Mapper                gencommons.Mapper
}

// TODO [SNOW-1501905]: handle other basic type variants
// TODO [SNOW-1501905]: handle any other interface (error)
// TODO [SNOW-1501905]: handle slices
// TODO [SNOW-1501905]: handle structs (chosen one or all)
func MapToSchemaField(field gencommons.Field) SchemaField {
	isPointer := field.IsPointer()
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	name := gencommons.ToSnakeCase(field.Name)
	switch concreteTypeWithoutPtr {
	case "string":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, gencommons.Identity}
	case "int":
		return SchemaField{name, schema.TypeInt, field.Name, isPointer, gencommons.Identity}
	case "float64":
		return SchemaField{name, schema.TypeFloat, field.Name, isPointer, gencommons.Identity}
	case "bool":
		return SchemaField{name, schema.TypeBool, field.Name, isPointer, gencommons.Identity}
	case "time.Time":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, gencommons.ToString}
	case "sdk.AccountObjectIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, gencommons.Name}
	case "sdk.AccountIdentifier", "sdk.ExternalObjectIdentifier", "sdk.DatabaseObjectIdentifier",
		"sdk.SchemaObjectIdentifier", "sdk.TableColumnIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, gencommons.FullyQualifiedName}
	case "sdk.ObjectIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, gencommons.FullyQualifiedName}
	}

	underlyingTypeWithoutPtr, _ := strings.CutPrefix(field.UnderlyingType, "*")
	isSdkDeclaredObject := strings.HasPrefix(concreteTypeWithoutPtr, "sdk.")
	switch {
	case isSdkDeclaredObject && underlyingTypeWithoutPtr == "string":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, gencommons.CastToString}
	case isSdkDeclaredObject && underlyingTypeWithoutPtr == "int":
		return SchemaField{name, schema.TypeInt, field.Name, isPointer, gencommons.CastToInt}
	}
	return SchemaField{name, schema.TypeInvalid, field.Name, isPointer, gencommons.Identity}
}
