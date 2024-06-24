package gen

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Mapper func(string) string

type SchemaField struct {
	Name                  string
	SchemaType            schema.ValueType
	OriginalName          string
	IsOriginalTypePointer bool
	Mapper                Mapper
}

var (
	Identity           = func(field string) string { return field }
	ToString           = func(field string) string { return fmt.Sprintf("%s.String()", field) }
	FullyQualifiedName = func(field string) string { return fmt.Sprintf("%s.FullyQualifiedName()", field) }
	CastToString       = func(field string) string { return fmt.Sprintf("string(%s)", field) }
	CastToInt          = func(field string) string { return fmt.Sprintf("int(%s)", field) }
)

// TODO [SNOW-1501905]: handle other basic type variants
// TODO [SNOW-1501905]: handle any other interface (error)
// TODO [SNOW-1501905]: handle slices
// TODO [SNOW-1501905]: handle structs (chosen one or all)
func MapToSchemaField(field Field) SchemaField {
	isPointer := field.IsPointer()
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	name := ToSnakeCase(field.Name)
	switch concreteTypeWithoutPtr {
	case "string":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, Identity}
	case "int":
		return SchemaField{name, schema.TypeInt, field.Name, isPointer, Identity}
	case "float64":
		return SchemaField{name, schema.TypeFloat, field.Name, isPointer, Identity}
	case "bool":
		return SchemaField{name, schema.TypeBool, field.Name, isPointer, Identity}
	case "time.Time":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, ToString}
	case "sdk.AccountIdentifier", "sdk.ExternalObjectIdentifier",
		"sdk.AccountObjectIdentifier", "sdk.DatabaseObjectIdentifier",
		"sdk.SchemaObjectIdentifier", "sdk.TableColumnIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, FullyQualifiedName}
	case "sdk.ObjectIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, FullyQualifiedName}
	}

	underlyingTypeWithoutPtr, _ := strings.CutPrefix(field.UnderlyingType, "*")
	switch {
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk.") && underlyingTypeWithoutPtr == "string":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, CastToString}
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk.") && underlyingTypeWithoutPtr == "int":
		return SchemaField{name, schema.TypeInt, field.Name, isPointer, CastToInt}
	}
	return SchemaField{name, schema.TypeInvalid, field.Name, isPointer, Identity}
}
