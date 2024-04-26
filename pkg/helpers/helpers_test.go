package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestDecodeSnowflakeParameterID(t *testing.T) {
	testCases := map[string]struct {
		id                 string
		fullyQualifiedName string
	}{
		"decodes quoted account object identifier": {
			id:                 `"test.name"`,
			fullyQualifiedName: `"test.name"`,
		},
		"decodes quoted database object identifier": {
			id:                 `"db"."test.name"`,
			fullyQualifiedName: `"db"."test.name"`,
		},
		"decodes quoted schema object identifier": {
			id:                 `"db"."schema"."test.name"`,
			fullyQualifiedName: `"db"."schema"."test.name"`,
		},
		"decodes quoted table column identifier": {
			id:                 `"db"."schema"."table.name"."test.name"`,
			fullyQualifiedName: `"db"."schema"."table.name"."test.name"`,
		},
		"decodes unquoted account object identifier": {
			id:                 `name`,
			fullyQualifiedName: `"name"`,
		},
		"decodes unquoted database object identifier": {
			id:                 `db.name`,
			fullyQualifiedName: `"db"."name"`,
		},
		"decodes unquoted schema object identifier": {
			id:                 `db.schema.name`,
			fullyQualifiedName: `"db"."schema"."name"`,
		},
		"decodes unquoted table column identifier": {
			id:                 `db.schema.table.name`,
			fullyQualifiedName: `"db"."schema"."table"."name"`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			id, err := DecodeSnowflakeParameterID(tc.id)
			require.NoError(t, err)
			require.Equal(t, tc.fullyQualifiedName, id.FullyQualifiedName())
		})
	}

	t.Run("identifier with too many parts", func(t *testing.T) {
		id := `this.identifier.is.too.long.to.be.decoded`
		_, err := DecodeSnowflakeParameterID(id)
		require.Errorf(t, err, "unable to classify identifier: %s", id)
	})

	t.Run("incompatible empty identifier", func(t *testing.T) {
		id := ""
		_, err := DecodeSnowflakeParameterID(id)
		require.Errorf(t, err, "incompatible identifier: %s", id)
	})

	t.Run("incompatible multiline identifier", func(t *testing.T) {
		id := "db.\nname"
		_, err := DecodeSnowflakeParameterID(id)
		require.Errorf(t, err, "incompatible identifier: %s", id)
	})
}

// TODO: add tests for non object identifiers
func TestEncodeSnowflakeID(t *testing.T) {
	testCases := map[string]struct {
		identifier        sdk.ObjectIdentifier
		expectedEncodedID string
	}{
		"encodes account object identifier": {
			identifier:        sdk.NewAccountObjectIdentifier("database"),
			expectedEncodedID: `database`,
		},
		"encodes quoted account object identifier": {
			identifier:        sdk.NewAccountObjectIdentifier("\"database\""),
			expectedEncodedID: `database`,
		},
		"encodes account object identifier with a dot": {
			identifier:        sdk.NewAccountObjectIdentifier("data.base"),
			expectedEncodedID: `data.base`,
		},
		"encodes pointer to account object identifier": {
			identifier:        sdk.Pointer(sdk.NewAccountObjectIdentifier("database")),
			expectedEncodedID: `database`,
		},
		"encodes database object identifier": {
			identifier:        sdk.NewDatabaseObjectIdentifier("database", "schema"),
			expectedEncodedID: `database|schema`,
		},
		"encodes quoted database object identifier": {
			identifier:        sdk.NewDatabaseObjectIdentifier("\"database\"", "\"schema\""),
			expectedEncodedID: `database|schema`,
		},
		"encodes database object identifier with dots": {
			identifier:        sdk.NewDatabaseObjectIdentifier("data.base", "sche.ma"),
			expectedEncodedID: `data.base|sche.ma`,
		},
		"encodes pointer to database object identifier": {
			identifier:        sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database", "schema")),
			expectedEncodedID: `database|schema`,
		},
		"encodes schema object identifier": {
			identifier:        sdk.NewSchemaObjectIdentifier("database", "schema", "table"),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes quoted schema object identifier": {
			identifier:        sdk.NewSchemaObjectIdentifier("\"database\"", "\"schema\"", "\"table\""),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes schema object identifier with dots": {
			identifier:        sdk.NewSchemaObjectIdentifier("data.base", "sche.ma", "tab.le"),
			expectedEncodedID: `data.base|sche.ma|tab.le`,
		},
		"encodes pointer to schema object identifier": {
			identifier:        sdk.Pointer(sdk.NewSchemaObjectIdentifier("database", "schema", "table")),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes table column identifier": {
			identifier:        sdk.NewTableColumnIdentifier("database", "schema", "table", "column"),
			expectedEncodedID: `database|schema|table|column`,
		},
		"encodes quoted table column identifier": {
			identifier:        sdk.NewTableColumnIdentifier("\"database\"", "\"schema\"", "\"table\"", "\"column\""),
			expectedEncodedID: `database|schema|table|column`,
		},
		"encodes table column identifier with dots": {
			identifier:        sdk.NewTableColumnIdentifier("data.base", "sche.ma", "tab.le", "col.umn"),
			expectedEncodedID: `data.base|sche.ma|tab.le|col.umn`,
		},
		"encodes pointer to table column identifier": {
			identifier:        sdk.Pointer(sdk.NewTableColumnIdentifier("database", "schema", "table", "column")),
			expectedEncodedID: `database|schema|table|column`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			encodedID := EncodeSnowflakeID(tc.identifier)
			require.Equal(t, tc.expectedEncodedID, encodedID)
		})
	}

	t.Run("panics for unsupported object identifier", func(t *testing.T) {
		id := unsupportedObjectIdentifier{}
		require.PanicsWithValue(t, fmt.Sprintf("Unsupported object identifier: %v", id), func() {
			EncodeSnowflakeID(id)
		})
	})

	nilTestCases := []any{
		(*sdk.AccountObjectIdentifier)(nil),
		(*sdk.DatabaseObjectIdentifier)(nil),
		(*sdk.SchemaObjectIdentifier)(nil),
		(*sdk.TableColumnIdentifier)(nil),
	}

	for i, tt := range nilTestCases {
		t.Run(fmt.Sprintf("handle nil pointer to object identifier %d", i), func(t *testing.T) {
			require.PanicsWithValue(t, "Nil object identifier received", func() {
				EncodeSnowflakeID(tt)
			})
		})
	}
}

type unsupportedObjectIdentifier struct{}

func (i unsupportedObjectIdentifier) Name() string {
	return "name"
}

func (i unsupportedObjectIdentifier) FullyQualifiedName() string {
	return "fully qualified name"
}

func TestSafelyDecodeSnowflakeID(t *testing.T) {
	type TestCase struct {
		Name       string
		ID         string
		ExpectedId sdk.ObjectIdentifier
		Error      string
	}

	accountObjectIdentifierTestCases := []TestCase{
		{
			Name:       "basic",
			ID:         "simple_identifier",
			ExpectedId: sdk.NewAccountObjectIdentifier("simple_identifier"),
		},
		{
			Name:       "identifier with dots",
			ID:         "identifier.with.dots",
			ExpectedId: sdk.NewAccountObjectIdentifier("identifier.with.dots"),
		},
		{
			Name:       "quoted",
			ID:         "\"quoted.identifier\"",
			ExpectedId: sdk.NewAccountObjectIdentifier("quoted.identifier"),
		},
	}

	for _, tt := range accountObjectIdentifierTestCases {
		t.Run("account object identifiers - "+tt.Name, func(t *testing.T) {
			id, err := SafelyDecodeSnowflakeID[sdk.AccountObjectIdentifier](tt.ID)

			if err != nil || tt.Error != "" {
				assert.ErrorContains(t, err, tt.Error)
			} else {
				assert.NotNil(t, id)
				assert.Equal(t, tt.ExpectedId.FullyQualifiedName(), id.FullyQualifiedName())
			}
		})
	}

	databaseObjectIdentifierTestCases := []TestCase{
		{
			Name:       "basic",
			ID:         "database.schema",
			ExpectedId: sdk.NewDatabaseObjectIdentifier("database", "schema"),
		},
		{
			Name:       "quoted",
			ID:         "\"database\".\"schema\"",
			ExpectedId: sdk.NewDatabaseObjectIdentifier("database", "schema"),
		},
		{
			Name:       "with docs and quoted",
			ID:         "\"data.base\".\"sche.ma\"",
			ExpectedId: sdk.NewDatabaseObjectIdentifier("data.base", "sche.ma"),
		},
		{
			Name:  "validation: empty identifier",
			ID:    "",
			Error: "Unable to parse the identifier: . Make sure you are using the correct form of the fully qualified name for this field: <database_name>.<name>",
		},
		{
			Name:  "validation: too much parts of the identifier",
			ID:    "data.base.schema.name.table.name",
			Error: "Unable to parse the identifier: data.base.schema.name.table.name. Make sure you are using the correct form of the fully qualified name for this field: <database_name>.<name>",
		},
		{
			Name:  "validation: invalid object type",
			ID:    "database.schema.table",
			Error: "expected DatabaseObjectIdentifier identifier type, but got: sdk.SchemaObjectIdentifier. The correct form of the fully qualified name for this field is: <database_name>.<name>, but was <database_name>.<schema_name>.<name>",
		},
	}

	for _, tt := range databaseObjectIdentifierTestCases {
		t.Run("database object identifiers - "+tt.Name, func(t *testing.T) {
			id, err := SafelyDecodeSnowflakeID[sdk.DatabaseObjectIdentifier](tt.ID)

			if err != nil || tt.Error != "" {
				assert.ErrorContains(t, err, tt.Error)
			} else {
				assert.NotNil(t, id)
				assert.Equal(t, tt.ExpectedId.FullyQualifiedName(), id.FullyQualifiedName())
			}
		})
	}

	schemaObjectIdentifierTestCases := []TestCase{
		{
			Name:       "basic",
			ID:         "database.schema.table",
			ExpectedId: sdk.NewSchemaObjectIdentifier("database", "schema", "table"),
		},
		{
			Name:       "quoted",
			ID:         "\"database\".\"schema\".\"table\"",
			ExpectedId: sdk.NewSchemaObjectIdentifier("database", "schema", "table"),
		},
		{
			Name:       "with docs and quoted",
			ID:         "\"data.base\".\"sche.ma\".\"ta.ble\"",
			ExpectedId: sdk.NewSchemaObjectIdentifier("data.base", "sche.ma", "ta.ble"),
		},
		{
			Name:  "validation: empty identifier",
			ID:    "",
			Error: "Unable to parse the identifier: . Make sure you are using the correct form of the fully qualified name for this field: <database_name>.<schema_name>.<name>",
		},
		{
			Name:  "validation: too much parts of the identifier",
			ID:    "data.base.schema.name.table.name",
			Error: "Unable to parse the identifier: data.base.schema.name.table.name. Make sure you are using the correct form of the fully qualified name for this field: <database_name>.<schema_name>.<name>",
		},
		{
			Name:  "validation: invalid object type",
			ID:    "database.schema",
			Error: "expected SchemaObjectIdentifier identifier type, but got: sdk.DatabaseObjectIdentifier. The correct form of the fully qualified name for this field is: <database_name>.<schema_name>.<name>, but was <database_name>.<name>",
		},
	}

	for _, tt := range schemaObjectIdentifierTestCases {
		t.Run("schema object identifiers - "+tt.Name, func(t *testing.T) {
			id, err := SafelyDecodeSnowflakeID[sdk.SchemaObjectIdentifier](tt.ID)

			if err != nil || tt.Error != "" {
				assert.ErrorContains(t, err, tt.Error)
			} else {
				assert.NotNil(t, id)
				assert.Equal(t, tt.ExpectedId.FullyQualifiedName(), id.FullyQualifiedName())
			}
		})
	}

	tableColumnIdentifierTestCases := []TestCase{
		{
			Name:       "basic",
			ID:         "database.schema.table.column_name",
			ExpectedId: sdk.NewTableColumnIdentifier("database", "schema", "table", "column_name"),
		},
		{
			Name:       "quoted",
			ID:         "\"database\".\"schema\".\"table\".\"column_name\"",
			ExpectedId: sdk.NewTableColumnIdentifier("database", "schema", "table", "column_name"),
		},
		{
			Name:       "with docs and quoted",
			ID:         "\"data.base\".\"sche.ma\".\"ta.ble\".\"column.name\"",
			ExpectedId: sdk.NewTableColumnIdentifier("data.base", "sche.ma", "ta.ble", "column.name"),
		},
		{
			Name:  "validation: empty identifier",
			ID:    "",
			Error: "Unable to parse the identifier: . Make sure you are using the correct form of the fully qualified name for this field: <database_name>.<schema_name>.<table_name>.<column_name>",
		},
		{
			Name:  "validation: too much parts of the identifier",
			ID:    "data.base.schema.name.table.name",
			Error: "Unable to parse the identifier: data.base.schema.name.table.name. Make sure you are using the correct form of the fully qualified name for this field: <database_name>.<schema_name>.<table_name>.<column_name>",
		},
		{
			Name:  "validation: invalid object type",
			ID:    "database.schema",
			Error: "expected TableColumnIdentifier identifier type, but got: sdk.DatabaseObjectIdentifier. The correct form of the fully qualified name for this field is: <database_name>.<schema_name>.<table_name>.<column_name>, but was <database_name>.<name>",
		},
	}

	for _, tt := range tableColumnIdentifierTestCases {
		t.Run("table column object identifiers - "+tt.Name, func(t *testing.T) {
			id, err := SafelyDecodeSnowflakeID[sdk.TableColumnIdentifier](tt.ID)

			if err != nil || tt.Error != "" {
				assert.ErrorContains(t, err, tt.Error)
			} else {
				assert.NotNil(t, id)
				assert.Equal(t, tt.ExpectedId.FullyQualifiedName(), id.FullyQualifiedName())
			}
		})
	}
}
