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
	IsOriginalTypePointer bool
	Mapper                Mapper
}

var Identity = func(field string) string { return field }
var ToString = func(field string) string { return fmt.Sprintf("%s.String()", field) }
var FullyQualifiedName = func(field string) string { return fmt.Sprintf("%s.FullyQualifiedName()", field) }
var CastToString = func(field string) string { return fmt.Sprintf("string(%s)", field) }
var CastToInt = func(field string) string { return fmt.Sprintf("int(%s)", field) }

// TODO: handle other basic type variants
// TODO: handle any other interface (error)
// TODO: handle slices
// TODO: handle structs (chosen one or all)
func MapToSchemaField(field Field) SchemaField {
	isPointer := field.IsPointer()
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	name := ToSnakeCase(field.Name)
	switch concreteTypeWithoutPtr {
	case "string":
		return SchemaField{name, schema.TypeString, isPointer, Identity}
	case "int":
		return SchemaField{name, schema.TypeInt, isPointer, Identity}
	case "float64":
		return SchemaField{name, schema.TypeFloat, isPointer, Identity}
	case "bool":
		return SchemaField{name, schema.TypeBool, isPointer, Identity}
	case "time.Time":
		return SchemaField{name, schema.TypeString, isPointer, ToString}
	case "sdk.AccountIdentifier", "sdk.ExternalObjectIdentifier",
		"sdk.AccountObjectIdentifier", "sdk.DatabaseObjectIdentifier",
		"sdk.SchemaObjectIdentifier", "sdk.TableColumnIdentifier":
		return SchemaField{name, schema.TypeString, isPointer, FullyQualifiedName}
	case "sdk.ObjectIdentifier":
		return SchemaField{name, schema.TypeString, isPointer, FullyQualifiedName}
	}

	underlyingTypeWithoutPtr, _ := strings.CutPrefix(field.UnderlyingType, "*")
	switch {
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk.") && underlyingTypeWithoutPtr == "string":
		return SchemaField{name, schema.TypeString, isPointer, CastToString}
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk.") && underlyingTypeWithoutPtr == "int":
		return SchemaField{name, schema.TypeInt, isPointer, CastToInt}
	}
	return SchemaField{name, schema.TypeInvalid, isPointer, Identity}
}
