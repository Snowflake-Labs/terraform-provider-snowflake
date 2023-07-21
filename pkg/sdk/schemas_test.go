package sdk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemasCreate(t *testing.T) {
	t.Run("clone", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			Clone: &Clone{
				SourceObject: NewAccountObjectIdentifier("sch1"),
				At: &TimeTravel{
					Timestamp: Pointer(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE SCHEMA CLONE "sch1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			OrReplace:                  Bool(true),
			Transient:                  Bool(true),
			IfNotExists:                Bool(true),
			WithManagedAccess:          Bool(true),
			DataRetentionTimeInDays:    Int(1),
			MaxDataExtensionTimeInDays: Int(1),
			DefaultDDLCollation:        String("en_US-trim"),
			Tag: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
					Value: "v1",
				},
			},
			Comment: String("comment"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE TRANSIENT SCHEMA IF NOT EXISTS WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 DEFAULT_DDL_COLLATION = 'en_US-trim' TAG ("db1"."schema1"."tag1" = 'v1') COMMENT = 'comment'`
		assert.Equal(t, expected, actual)
	})
}

func TestSchemasAlter(t *testing.T) {
	t.Run("rename to", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     NewSchemaIdentifier("database_name", "schema_name"),
			IfExists: Bool(true),
			NewName:  NewSchemaIdentifier("database_name", "new_schema_name"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SCHEMA IF EXISTS "database_name"."schema_name" RENAME TO "database_name"."new_schema_name"`
		assert.Equal(t, expected, actual)
	})

	t.Run("swap with", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     NewSchemaIdentifier("database_name", "schema_name"),
			IfExists: Bool(false),
			SwapWith: NewSchemaIdentifier("database_name", "target_schema_name"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SCHEMA "database_name"."schema_name" SWAP WITH "database_name"."target_schema_name"`
		assert.Equal(t, expected, actual)
	})

	t.Run("set options", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: NewSchemaIdentifier("database_name", "schema_name"),
			Set: &SchemaSet{
				DataRetentionTimeInDays:    Int(3),
				MaxDataExtensionTimeInDays: Int(2),
				DefaultDDLCollation:        String("en_US-trim"),
				Comment:                    String("comment"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SCHEMA "database_name"."schema_name" SET DATA_RETENTION_TIME_IN_DAYS = 3, MAX_DATA_EXTENSION_TIME_IN_DAYS = 2, DEFAULT_DDL_COLLATION = 'en_US-trim', COMMENT = 'comment'`
		assert.Equal(t, expected, actual)
	})

	t.Run("set tags", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: NewSchemaIdentifier("database_name", "schema_name"),
			Set: &SchemaSet{
				Tag: []TagAssociation{
					{
						Name:  NewAccountObjectIdentifier("tag1"),
						Value: "value1",
					},
					{
						Name:  NewAccountObjectIdentifier("tag2"),
						Value: "value2",
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SCHEMA "database_name"."schema_name" SET TAG "tag1" = 'value1', "tag2" = 'value2'`
		assert.Equal(t, expected, actual)
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: NewSchemaIdentifier("database_name", "schema_name"),
			Unset: &SchemaUnset{
				Tag: []ObjectIdentifier{
					NewAccountObjectIdentifier("tag1"),
					NewAccountObjectIdentifier("tag2"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SCHEMA "database_name"."schema_name" UNSET TAG "tag1", "tag2"`
		assert.Equal(t, expected, actual)
	})

	t.Run("unset options", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: NewSchemaIdentifier("database_name", "schema_name"),
			Unset: &SchemaUnset{
				DataRetentionTimeInDays:    Bool(true),
				MaxDataExtensionTimeInDays: Bool(true),
				DefaultDDLCollation:        Bool(true),
				Comment:                    Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SCHEMA "database_name"."schema_name" UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, DEFAULT_DDL_COLLATION, COMMENT`
		assert.Equal(t, expected, actual)
	})

	t.Run("enable managed access", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:                NewSchemaIdentifier("database_name", "schema_name"),
			EnableManagedAccess: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SCHEMA "database_name"."schema_name" ENABLE MANAGED ACCESS`
		assert.Equal(t, expected, actual)
	})

	t.Run("disable managed access", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:                 NewSchemaIdentifier("database_name", "schema_name"),
			DisableMangaedAccess: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SCHEMA "database_name"."schema_name" DISABLE MANAGED ACCESS`
		assert.Equal(t, expected, actual)
	})
}

func TestSchemasDrop(t *testing.T) {
	t.Run("cascade", func(t *testing.T) {
		opts := &DropSchemaOptions{
			IfExists: Bool(true),
			name:     NewSchemaIdentifier("database_name", "schema_name"),
			Cascade:  Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP SCHEMA IF EXISTS "database_name"."schema_name" CASCADE`
		assert.Equal(t, expected, actual)
	})

	t.Run("restrict", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name:     NewSchemaIdentifier("database_name", "schema_name"),
			Restrict: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP SCHEMA "database_name"."schema_name" RESTRICT`
		assert.Equal(t, expected, actual)
	})
}

func TestSchemasUndrop(t *testing.T) {
	opts := &undropSchemaOptions{
		name: NewSchemaIdentifier("database_name", "schema_name"),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `UNDROP SCHEMA "database_name"."schema_name"`
	assert.Equal(t, expected, actual)
}

func TestSchemasDescribe(t *testing.T) {
	opts := &describeSchemaOptions{
		name: NewSchemaIdentifier("database_name", "schema_name"),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `DESCRIBE SCHEMA "database_name"."schema_name"`
	assert.Equal(t, expected, actual)
}

func TestSchemasShow(t *testing.T) {
	t.Run("like", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			Like: &Like{
				Pattern: String("schema_pattern"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE SCHEMAS HISTORY LIKE 'schema_pattern'`
		assert.Equal(t, expected, actual)
	})

	t.Run("in account", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			In: &InSchema{
				Account: Bool(true),
				Name:    NewAccountObjectIdentifier("account_name"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE SCHEMAS HISTORY IN ACCOUNT "account_name"`
		assert.Equal(t, expected, actual)
	})

	t.Run("in database", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			In: &InSchema{
				Database: Bool(true),
				Name:     NewAccountObjectIdentifier("database_name"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE SCHEMAS HISTORY IN DATABASE "database_name"`
		assert.Equal(t, expected, actual)
	})

	t.Run("starts with", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:      Bool(true),
			History:    Bool(true),
			StartsWith: String("schema_pattern"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE SCHEMAS HISTORY STARTS WITH 'schema_pattern'`
		assert.Equal(t, expected, actual)
	})

	t.Run("limit", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			LimitFrom: &LimitFrom{
				Rows: Int(3),
				From: String("name_string"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE SCHEMAS HISTORY LIMIT 3 FROM 'name_string'`
		assert.Equal(t, expected, actual)
	})
}
