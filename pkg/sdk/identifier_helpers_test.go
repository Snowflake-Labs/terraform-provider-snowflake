package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccountIdentifierFromFullyQualifiedName(t *testing.T) {
	type test struct {
		input string
		want  AccountIdentifier
	}

	tests := []test{
		{input: "BSB98216", want: AccountIdentifier{accountLocator: "BSB98216"}},
		{input: "SNOW.MY_TEST_ACCOUNT", want: AccountIdentifier{organizationName: "SNOW", accountName: "MY_TEST_ACCOUNT"}},
		{input: "\"SNOW\".\"MY_TEST_ACCOUNT\"", want: AccountIdentifier{organizationName: "SNOW", accountName: "MY_TEST_ACCOUNT"}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			id := NewAccountIdentifierFromFullyQualifiedName(tc.input)
			require.Equal(t, tc.want, id)
		})
	}
}

func TestNewSchemaObjectIdentifierFromFullyQualifiedName(t *testing.T) {
	type test struct {
		input string
		want  SchemaObjectIdentifier
	}

	tests := []test{
		{input: "\"MY_DB\".\"MY_SCHEMA\".\"multiply\"(number, number)", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "multiply", arguments: []DataType{DataTypeNumber, DataTypeNumber}}},
		{input: "MY_DB.MY_SCHEMA.add(number, number)", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "add", arguments: []DataType{DataTypeNumber, DataTypeNumber}}},
		{input: "\"MY_DB\".\"MY_SCHEMA\".\"MY_UDF\"()", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "MY_UDF", arguments: []DataType{}}},
		{input: "\"MY_DB\".\"MY_SCHEMA\".\"MY_PIPE\"", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "MY_PIPE", arguments: nil}},
		{input: "MY_DB.MY_SCHEMA.MY_STAGE", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "MY_STAGE", arguments: nil}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			id := NewSchemaObjectIdentifierFromFullyQualifiedName(tc.input)
			require.Equal(t, tc.want, id)
		})
	}
}

func TestDatabaseObjectIdentifier(t *testing.T) {
	t.Run("create from strings", func(t *testing.T) {
		identifier := NewDatabaseObjectIdentifier("aaa", "bbb")

		assert.Equal(t, "aaa", identifier.DatabaseName())
		assert.Equal(t, "bbb", identifier.Name())
	})

	t.Run("create from quoted strings", func(t *testing.T) {
		identifier := NewDatabaseObjectIdentifier(`"aaa"`, `"bbb"`)

		assert.Equal(t, "aaa", identifier.DatabaseName())
		assert.Equal(t, "bbb", identifier.Name())
	})

	t.Run("create from fully qualified name", func(t *testing.T) {
		identifier := NewDatabaseObjectIdentifierFromFullyQualifiedName("aaa.bbb")

		assert.Equal(t, "aaa", identifier.DatabaseName())
		assert.Equal(t, "bbb", identifier.Name())
	})

	t.Run("create from quoted fully qualified name", func(t *testing.T) {
		identifier := NewDatabaseObjectIdentifierFromFullyQualifiedName(`"aaa"."bbb"`)

		assert.Equal(t, "aaa", identifier.DatabaseName())
		assert.Equal(t, "bbb", identifier.Name())
	})

	t.Run("get fully qualified name", func(t *testing.T) {
		identifier := DatabaseObjectIdentifier{"aaa", "bbb"}

		assert.Equal(t, `"aaa"."bbb"`, identifier.FullyQualifiedName())
	})
}

func TestNewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName(t *testing.T) {
	testCases := []struct {
		RawInput string
		Input    SchemaObjectIdentifierWithArguments
		Error    string
	}{
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, DataTypeNumber, DataTypeTimestampTZ)},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)")},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, "VECTOR(INT, 20)", DataTypeFloat)},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)", "VECTOR(INT, 10)")},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeTime, "VECTOR(INT, 20)", "VECTOR(FLOAT, 10)", DataTypeFloat)},
		// TODO(SNOW-1571674): Won't work, because of the assumption that identifiers are not containing '(' and ')' parentheses
		{Input: NewSchemaObjectIdentifierWithArguments(`ab()c`, `def()`, `()ghi`, DataTypeTime, "VECTOR(INT, 20)", "VECTOR(FLOAT, 10)", DataTypeFloat), Error: `unable to read identifier: "ab`},
		{Input: NewSchemaObjectIdentifierWithArguments(`ab(,)c`, `,def()`, `()ghi,`, DataTypeTime, "VECTOR(INT, 20)", "VECTOR(FLOAT, 10)", DataTypeFloat), Error: `unable to read identifier: "ab`},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`)},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`), RawInput: `abc.def.ghi()`},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)"), RawInput: `abc.def.ghi(FLOAT, VECTOR(INT, 20))`},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("processing %s", testCase.Input.FullyQualifiedName()), func(t *testing.T) {
			var id SchemaObjectIdentifierWithArguments
			var err error
			if testCase.RawInput != "" {
				id, err = NewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName(testCase.RawInput)
			} else {
				id, err = NewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName(testCase.Input.FullyQualifiedName())
			}

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.Input.FullyQualifiedName(), id.FullyQualifiedName())
			}
		})
	}
}
