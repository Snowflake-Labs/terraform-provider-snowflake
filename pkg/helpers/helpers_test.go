package helpers

import (
	"testing"

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
			identifier:        sdk.NewAccountObjectIdentifier("account"),
			expectedEncodedID: `account`,
		},
		"encodes quoted account object identifier": {
			identifier:        sdk.NewAccountObjectIdentifier("\"account\""),
			expectedEncodedID: `account`,
		},
		"encodes account object identifier with a dot": {
			identifier:        sdk.NewAccountObjectIdentifier("acc.ount"),
			expectedEncodedID: `acc.ount`,
		},
		"encodes database object identifier": {
			identifier:        sdk.NewDatabaseObjectIdentifier("account", "database"),
			expectedEncodedID: `account|database`,
		},
		"encodes quoted database object identifier": {
			identifier:        sdk.NewDatabaseObjectIdentifier("\"account\"", "\"database\""),
			expectedEncodedID: `account|database`,
		},
		"encodes database object identifier with dots": {
			identifier:        sdk.NewDatabaseObjectIdentifier("acc.ount", "data.base"),
			expectedEncodedID: `acc.ount|data.base`,
		},
		"encodes schema object identifier": {
			identifier:        sdk.NewSchemaObjectIdentifier("account", "database", "schema"),
			expectedEncodedID: `account|database|schema`,
		},
		"encodes quoted schema object identifier": {
			identifier:        sdk.NewSchemaObjectIdentifier("\"account\"", "\"database\"", "\"schema\""),
			expectedEncodedID: `account|database|schema`,
		},
		"encodes schema object identifier with dots": {
			identifier:        sdk.NewSchemaObjectIdentifier("acc.ount", "data.base", "sche.ma"),
			expectedEncodedID: `acc.ount|data.base|sche.ma`,
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
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			encodedID := EncodeSnowflakeID(tc.identifier)
			require.Equal(t, tc.expectedEncodedID, encodedID)
		})
	}
}
