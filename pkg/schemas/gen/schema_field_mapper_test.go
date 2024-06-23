package gen

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func Test_MapToSchemaField(t *testing.T) {

	assertSchemaFieldMapped := func(schemaField SchemaField, expectedName string, expectedSchemaType schema.ValueType, expectedOriginalPointer bool, expectedMapper Mapper) {
		assert.Equal(t, expectedName, schemaField.Name)
		assert.Equal(t, expectedSchemaType, schemaField.SchemaType)
		assert.Equal(t, expectedOriginalPointer, schemaField.IsOriginalTypePointer)
		// TODO: ugly comparison of functions with the current implementation of mapper
		assert.Equal(t, reflect.ValueOf(expectedMapper).Pointer(), reflect.ValueOf(schemaField.Mapper).Pointer())
	}

	t.Run("test schema field mapper", func(t *testing.T) {
		stringField := Field{"unexportedString", "string", "string"}

		schemaField := MapToSchemaField(stringField)

		assertSchemaFieldMapped(schemaField, "unexported_string", schema.TypeString, false, Identity)
	})
}
