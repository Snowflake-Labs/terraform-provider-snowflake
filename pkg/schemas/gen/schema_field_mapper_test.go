package gen

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func Test_MapToSchemaField(t *testing.T) {
	type expectedValues struct {
		name       string
		schemaType schema.ValueType
		isPointer  bool
		mapper     Mapper
	}

	testCases := []struct {
		field    Field
		expected expectedValues
	}{
		{
			field:    Field{"unexportedString", "string", "string"},
			expected: expectedValues{"unexported_string", schema.TypeString, false, Identity},
		},
		{
			field:    Field{"unexportedInt", "int", "int"},
			expected: expectedValues{"unexported_int", schema.TypeInt, false, Identity},
		},
		{
			field:    Field{"unexportedBool", "bool", "bool"},
			expected: expectedValues{"unexported_bool", schema.TypeBool, false, Identity},
		},
		{
			field:    Field{"unexportedFloat64", "float64", "float64"},
			expected: expectedValues{"unexported_float64", schema.TypeFloat, false, Identity},
		},
		{
			field:    Field{"unexportedStringPtr", "*string", "*string"},
			expected: expectedValues{"unexported_string_ptr", schema.TypeString, true, Identity},
		},
		{
			field:    Field{"unexportedIntPtr", "*int", "*int"},
			expected: expectedValues{"unexported_int_ptr", schema.TypeInt, true, Identity},
		},
		{
			field:    Field{"unexportedBoolPtr", "*bool", "*bool"},
			expected: expectedValues{"unexported_bool_ptr", schema.TypeBool, true, Identity},
		},
		{
			field:    Field{"unexportedFloat64Ptr", "*float64", "*float64"},
			expected: expectedValues{"unexported_float64_ptr", schema.TypeFloat, true, Identity},
		},
		{
			field:    Field{"unexportedTime", "time.Time", "struct"},
			expected: expectedValues{"unexported_time", schema.TypeString, false, ToString},
		},
		{
			field:    Field{"unexportedTimePtr", "*time.Time", "*struct"},
			expected: expectedValues{"unexported_time_ptr", schema.TypeString, true, ToString},
		},
		{
			field:    Field{"unexportedStringEnum", "sdk.WarehouseType", "string"},
			expected: expectedValues{"unexported_string_enum", schema.TypeString, false, CastToString},
		},
		{
			field:    Field{"unexportedStringEnumPtr", "*sdk.WarehouseType", "*string"},
			expected: expectedValues{"unexported_string_enum_ptr", schema.TypeString, true, CastToString},
		},
		{
			field:    Field{"unexportedIntEnum", "sdk.ResourceMonitorLevel", "int"},
			expected: expectedValues{"unexported_int_enum", schema.TypeInt, false, CastToInt},
		},
		{
			field:    Field{"unexportedIntEnumPtr", "*sdk.ResourceMonitorLevel", "*int"},
			expected: expectedValues{"unexported_int_enum_ptr", schema.TypeInt, true, CastToInt},
		},
		{
			field:    Field{"unexportedAccountIdentifier", "sdk.AccountIdentifier", "struct"},
			expected: expectedValues{"unexported_account_identifier", schema.TypeString, false, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedExternalObjectIdentifier", "sdk.ExternalObjectIdentifier", "struct"},
			expected: expectedValues{"unexported_external_object_identifier", schema.TypeString, false, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedAccountObjectIdentifier", "sdk.AccountObjectIdentifier", "struct"},
			expected: expectedValues{"unexported_account_object_identifier", schema.TypeString, false, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedDatabaseObjectIdentifier", "sdk.DatabaseObjectIdentifier", "struct"},
			expected: expectedValues{"unexported_database_object_identifier", schema.TypeString, false, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedSchemaObjectIdentifier", "sdk.SchemaObjectIdentifier", "struct"},
			expected: expectedValues{"unexported_schema_object_identifier", schema.TypeString, false, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedTableColumnIdentifier", "sdk.TableColumnIdentifier", "struct"},
			expected: expectedValues{"unexported_table_column_identifier", schema.TypeString, false, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedAccountIdentifierPtr", "*sdk.AccountIdentifier", "*struct"},
			expected: expectedValues{"unexported_account_identifier_ptr", schema.TypeString, true, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedExternalObjectIdentifierPtr", "*sdk.ExternalObjectIdentifier", "*struct"},
			expected: expectedValues{"unexported_external_object_identifier_ptr", schema.TypeString, true, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedAccountObjectIdentifierPtr", "*sdk.AccountObjectIdentifier", "*struct"},
			expected: expectedValues{"unexported_account_object_identifier_ptr", schema.TypeString, true, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedDatabaseObjectIdentifierPtr", "*sdk.DatabaseObjectIdentifier", "*struct"},
			expected: expectedValues{"unexported_database_object_identifier_ptr", schema.TypeString, true, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedSchemaObjectIdentifierPtr", "*sdk.SchemaObjectIdentifier", "*struct"},
			expected: expectedValues{"unexported_schema_object_identifier_ptr", schema.TypeString, true, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedTableColumnIdentifierPtr", "*sdk.TableColumnIdentifier", "*struct"},
			expected: expectedValues{"unexported_table_column_identifier_ptr", schema.TypeString, true, FullyQualifiedName},
		},
		{
			field:    Field{"unexportedInterface", "sdk.ObjectIdentifier", "interface"},
			expected: expectedValues{"unexported_interface", schema.TypeString, false, FullyQualifiedName},
		},
	}

	assertSchemaFieldMapped := func(schemaField SchemaField, originalField Field, expected expectedValues) {
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
