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
			field:    gencommons.Field{"unexportedString", "string", "string"},
			expected: expectedValues{"unexported_string", schema.TypeString, false, gencommons.Identity},
		},
		{
			field:    gencommons.Field{"unexportedInt", "int", "int"},
			expected: expectedValues{"unexported_int", schema.TypeInt, false, gencommons.Identity},
		},
		{
			field:    gencommons.Field{"unexportedBool", "bool", "bool"},
			expected: expectedValues{"unexported_bool", schema.TypeBool, false, gencommons.Identity},
		},
		{
			field:    gencommons.Field{"unexportedFloat64", "float64", "float64"},
			expected: expectedValues{"unexported_float64", schema.TypeFloat, false, gencommons.Identity},
		},
		{
			field:    gencommons.Field{"unexportedStringPtr", "*string", "*string"},
			expected: expectedValues{"unexported_string_ptr", schema.TypeString, true, gencommons.Identity},
		},
		{
			field:    gencommons.Field{"unexportedIntPtr", "*int", "*int"},
			expected: expectedValues{"unexported_int_ptr", schema.TypeInt, true, gencommons.Identity},
		},
		{
			field:    gencommons.Field{"unexportedBoolPtr", "*bool", "*bool"},
			expected: expectedValues{"unexported_bool_ptr", schema.TypeBool, true, gencommons.Identity},
		},
		{
			field:    gencommons.Field{"unexportedFloat64Ptr", "*float64", "*float64"},
			expected: expectedValues{"unexported_float64_ptr", schema.TypeFloat, true, gencommons.Identity},
		},
		{
			field:    gencommons.Field{"unexportedTime", "time.Time", "struct"},
			expected: expectedValues{"unexported_time", schema.TypeString, false, gencommons.ToString},
		},
		{
			field:    gencommons.Field{"unexportedTimePtr", "*time.Time", "*struct"},
			expected: expectedValues{"unexported_time_ptr", schema.TypeString, true, gencommons.ToString},
		},
		{
			field:    gencommons.Field{"unexportedStringEnum", "sdk.WarehouseType", "string"},
			expected: expectedValues{"unexported_string_enum", schema.TypeString, false, gencommons.CastToString},
		},
		{
			field:    gencommons.Field{"unexportedStringEnumPtr", "*sdk.WarehouseType", "*string"},
			expected: expectedValues{"unexported_string_enum_ptr", schema.TypeString, true, gencommons.CastToString},
		},
		{
			field:    gencommons.Field{"unexportedIntEnum", "sdk.ResourceMonitorLevel", "int"},
			expected: expectedValues{"unexported_int_enum", schema.TypeInt, false, gencommons.CastToInt},
		},
		{
			field:    gencommons.Field{"unexportedIntEnumPtr", "*sdk.ResourceMonitorLevel", "*int"},
			expected: expectedValues{"unexported_int_enum_ptr", schema.TypeInt, true, gencommons.CastToInt},
		},
		{
			field:    gencommons.Field{"unexportedAccountIdentifier", "sdk.AccountIdentifier", "struct"},
			expected: expectedValues{"unexported_account_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedExternalObjectIdentifier", "sdk.ExternalObjectIdentifier", "struct"},
			expected: expectedValues{"unexported_external_object_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedAccountObjectIdentifier", "sdk.AccountObjectIdentifier", "struct"},
			expected: expectedValues{"unexported_account_object_identifier", schema.TypeString, false, gencommons.Name},
		},
		{
			field:    gencommons.Field{"unexportedDatabaseObjectIdentifier", "sdk.DatabaseObjectIdentifier", "struct"},
			expected: expectedValues{"unexported_database_object_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedSchemaObjectIdentifier", "sdk.SchemaObjectIdentifier", "struct"},
			expected: expectedValues{"unexported_schema_object_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedTableColumnIdentifier", "sdk.TableColumnIdentifier", "struct"},
			expected: expectedValues{"unexported_table_column_identifier", schema.TypeString, false, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedAccountIdentifierPtr", "*sdk.AccountIdentifier", "*struct"},
			expected: expectedValues{"unexported_account_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedExternalObjectIdentifierPtr", "*sdk.ExternalObjectIdentifier", "*struct"},
			expected: expectedValues{"unexported_external_object_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedAccountObjectIdentifierPtr", "*sdk.AccountObjectIdentifier", "*struct"},
			expected: expectedValues{"unexported_account_object_identifier_ptr", schema.TypeString, true, gencommons.Name},
		},
		{
			field:    gencommons.Field{"unexportedDatabaseObjectIdentifierPtr", "*sdk.DatabaseObjectIdentifier", "*struct"},
			expected: expectedValues{"unexported_database_object_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedSchemaObjectIdentifierPtr", "*sdk.SchemaObjectIdentifier", "*struct"},
			expected: expectedValues{"unexported_schema_object_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedTableColumnIdentifierPtr", "*sdk.TableColumnIdentifier", "*struct"},
			expected: expectedValues{"unexported_table_column_identifier_ptr", schema.TypeString, true, gencommons.FullyQualifiedName},
		},
		{
			field:    gencommons.Field{"unexportedInterface", "sdk.ObjectIdentifier", "interface"},
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
