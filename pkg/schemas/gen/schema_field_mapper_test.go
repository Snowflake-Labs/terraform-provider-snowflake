package gen

import (
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func Test_MapToSchemaField(t *testing.T) {
	type expectedValues struct {
		name       string
		schemaType schema.ValueType
		isPointer  bool
		mapper     gencommons.Mapper
	}

	testCases := []struct {
		field    gencommons.Field
		expected expectedValues
	}{
		{
			field:    gencommons.Field{Name: "unexportedString", ConcreteType: "string", UnderlyingType: "string"},
			expected: expectedValues{"unexported_string", schema.TypeString, false, gencommons.Identity},
		},
		{
			field:    gencommons.Field{Name: "unexportedInt", ConcreteType: "int", UnderlyingType: "int"},
			expected: expectedValues{"unexported_int", schema.TypeInt, false, gencommons.Identity},
		},
		{
			field:    gencommons.Field{Name: "unexportedBool", ConcreteType: "bool", UnderlyingType: "bool"},
			expected: expectedValues{"unexported_bool", schema.TypeBool, false, gencommons.Identity},
		},
		{
			field:    gencommons.Field{Name: "unexportedFloat64", ConcreteType: "float64", UnderlyingType: "float64"},
			expected: expectedValues{"unexported_float64", schema.TypeFloat, false, gencommons.Identity},
		},
		{
			field:    gencommons.Field{Name: "unexportedStringPtr", ConcreteType: "*string", UnderlyingType: "*string"},
			expected: expectedValues{"unexported_string_ptr", schema.TypeString, true, gencommons.Identity},
		},
		{
			field:    gencommons.Field{Name: "unexportedIntPtr", ConcreteType: "*int", UnderlyingType: "*int"},
			expected: expectedValues{"unexported_int_ptr", schema.TypeInt, true, gencommons.Identity},
		},
		{
			field:    gencommons.Field{Name: "unexportedBoolPtr", ConcreteType: "*bool", UnderlyingType: "*bool"},
			expected: expectedValues{"unexported_bool_ptr", schema.TypeBool, true, gencommons.Identity},
		},
		{
			field:    gencommons.Field{Name: "unexportedFloat64Ptr", ConcreteType: "*float64", UnderlyingType: "*float64"},
			expected: expectedValues{"unexported_float64_ptr", schema.TypeFloat, true, gencommons.Identity},
		},
		{
			field:    gencommons.Field{Name: "unexportedTime", ConcreteType: "time.Time", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_time", schema.TypeString, false, gencommons.ToString},
		},
		{
			field:    gencommons.Field{Name: "unexportedTimePtr", ConcreteType: "*time.Time", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_time_ptr", schema.TypeString, true, gencommons.ToString},
		},
		{
			field:    gencommons.Field{Name: "unexportedStringEnum", ConcreteType: "sdk.WarehouseType", UnderlyingType: "string"},
			expected: expectedValues{"unexported_string_enum", schema.TypeString, false, gencommons.CastToString},
		},
		{
			field:    gencommons.Field{Name: "unexportedStringEnumPtr", ConcreteType: "*sdk.WarehouseType", UnderlyingType: "*string"},
			expected: expectedValues{"unexported_string_enum_ptr", schema.TypeString, true, gencommons.CastToString},
		},
		{
			field:    gencommons.Field{Name: "unexportedIntEnum", ConcreteType: "sdk.ResourceMonitorLevel", UnderlyingType: "int"},
			expected: expectedValues{"unexported_int_enum", schema.TypeInt, false, gencommons.CastToInt},
		},
		{
			field:    gencommons.Field{Name: "unexportedIntEnumPtr", ConcreteType: "*sdk.ResourceMonitorLevel", UnderlyingType: "*int"},
			expected: expectedValues{"unexported_int_enum_ptr", schema.TypeInt, true, gencommons.CastToInt},
		},
		{
			field:    gencommons.Field{Name: "unexportedAccountIdentifier", ConcreteType: "sdk.AccountIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_account_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedExternalObjectIdentifier", ConcreteType: "sdk.ExternalObjectIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_external_object_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedAccountObjectIdentifier", ConcreteType: "sdk.AccountObjectIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_account_object_identifier", schema.TypeString, false, gencommons.Name},
		},
		{
			field:    gencommons.Field{Name: "unexportedDatabaseObjectIdentifier", ConcreteType: "sdk.DatabaseObjectIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_database_object_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedSchemaObjectIdentifier", ConcreteType: "sdk.SchemaObjectIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_schema_object_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedTableColumnIdentifier", ConcreteType: "sdk.TableColumnIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_table_column_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedAccountIdentifierPtr", ConcreteType: "*sdk.AccountIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_account_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedExternalObjectIdentifierPtr", ConcreteType: "*sdk.ExternalObjectIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_external_object_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedAccountObjectIdentifierPtr", ConcreteType: "*sdk.AccountObjectIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_account_object_identifier_ptr", schema.TypeString, true, gencommons.Name},
		},
		{
			field:    gencommons.Field{Name: "unexportedDatabaseObjectIdentifierPtr", ConcreteType: "*sdk.DatabaseObjectIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_database_object_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedSchemaObjectIdentifierPtr", ConcreteType: "*sdk.SchemaObjectIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_schema_object_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedTableColumnIdentifierPtr", ConcreteType: "*sdk.TableColumnIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_table_column_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{Name: "unexportedInterface", ConcreteType: "sdk.ObjectIdentifier", UnderlyingType: "interface"},
			expected: expectedValues{"unexported_interface", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
	}

	assertSchemaFieldMapped := func(schemaField SchemaField, originalField gencommons.Field, expected expectedValues) {
		assert.Equal(t, expected.name, schemaField.Name)
		assert.Equal(t, expected.schemaType, schemaField.SchemaType)
		assert.Equal(t, originalField.Name, schemaField.OriginalName)
		assert.Equal(t, expected.isPointer, schemaField.IsOriginalTypePointer)
		// TODO [SNOW-1501905]: ugly comparison of functions with the current implementation of mapper
		assert.Equal(t, reflect.ValueOf(expected.mapper).Pointer(), reflect.ValueOf(schemaField.Mapper).Pointer())
	}

	for _, tc := range testCases {
		t.Run(tc.field.Name, func(t *testing.T) {
			schemaField := MapToSchemaField(tc.field)

			assertSchemaFieldMapped(schemaField, tc.field, tc.expected)
		})
	}
}
