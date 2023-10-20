package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeSnowflakeParameterID(t *testing.T) {
	t.Run("decodes quoted account object identifier", func(t *testing.T) {
		id := `"test.name"`
		accObjId, err := DecodeSnowflakeParameterID(id)
		accObjId = accObjId.(sdk.AccountObjectIdentifier)
		require.NoError(t, err)
		require.Equal(t, id, accObjId.FullyQualifiedName())
	})

	t.Run("decodes quoted database object identifier", func(t *testing.T) {
		id := `"db"."test.name"`
		dbObjId, err := DecodeSnowflakeParameterID(id)
		dbObjId = dbObjId.(sdk.DatabaseObjectIdentifier)
		require.NoError(t, err)
		require.Equal(t, id, dbObjId.FullyQualifiedName())
	})

	t.Run("decodes quoted schema object identifier", func(t *testing.T) {
		id := `"db"."schema"."test.name"`
		schemaObjId, err := DecodeSnowflakeParameterID(id)
		schemaObjId = schemaObjId.(sdk.SchemaObjectIdentifier)
		require.NoError(t, err)
		require.Equal(t, id, schemaObjId.FullyQualifiedName())
	})

	t.Run("decodes quoted table column identifier", func(t *testing.T) {
		id := `"db"."schema"."table.name"."test.name"`
		schemaObjId, err := DecodeSnowflakeParameterID(id)
		schemaObjId = schemaObjId.(sdk.TableColumnIdentifier)
		require.NoError(t, err)
		require.Equal(t, id, schemaObjId.FullyQualifiedName())
	})

	t.Run("decodes unquoted account object identifier", func(t *testing.T) {
		id := `name`
		accObjId, err := DecodeSnowflakeParameterID(id)
		accObjId = accObjId.(sdk.AccountObjectIdentifier)
		require.NoError(t, err)
		require.Equal(t, `"name"`, accObjId.FullyQualifiedName())
	})

	t.Run("decodes unquoted database object identifier", func(t *testing.T) {
		id := `db.name`
		dbObjId, err := DecodeSnowflakeParameterID(id)
		dbObjId = dbObjId.(sdk.DatabaseObjectIdentifier)
		require.NoError(t, err)
		require.Equal(t, `"db"."name"`, dbObjId.FullyQualifiedName())
	})

	t.Run("decodes unquoted schema object identifier", func(t *testing.T) {
		id := `db.schema.name`
		schemaObjId, err := DecodeSnowflakeParameterID(id)
		schemaObjId = schemaObjId.(sdk.SchemaObjectIdentifier)
		require.NoError(t, err)
		require.Equal(t, `"db"."schema"."name"`, schemaObjId.FullyQualifiedName())
	})

	t.Run("decodes unquoted table column identifier", func(t *testing.T) {
		id := `db.schema.table.name`
		schemaObjId, err := DecodeSnowflakeParameterID(id)
		schemaObjId = schemaObjId.(sdk.TableColumnIdentifier)
		require.NoError(t, err)
		require.Equal(t, `"db"."schema"."table"."name"`, schemaObjId.FullyQualifiedName())
	})

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
