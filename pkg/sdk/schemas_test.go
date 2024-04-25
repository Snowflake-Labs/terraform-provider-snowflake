package sdk

import (
	"testing"
	"time"
)

func TestSchemasCreate(t *testing.T) {
	id := randomDatabaseObjectIdentifier()

	t.Run("clone", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			name:      id,
			OrReplace: Bool(true),
			Clone: &Clone{
				SourceObject: NewAccountObjectIdentifier("sch1"),
				At: &TimeTravel{
					Timestamp: Pointer(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SCHEMA %s CLONE "sch1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`, id.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			name:                       id,
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE TRANSIENT SCHEMA IF NOT EXISTS %s WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 DEFAULT_DDL_COLLATION = 'en_US-trim' TAG ("db1"."schema1"."tag1" = 'v1') COMMENT = 'comment'`, id.FullyQualifiedName())
	})
}

func TestSchemasAlter(t *testing.T) {
	t.Run("rename to", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     NewDatabaseObjectIdentifier("database_name", "schema_name"),
			IfExists: Bool(true),
			NewName:  NewDatabaseObjectIdentifier("database_name", "new_schema_name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA IF EXISTS "database_name"."schema_name" RENAME TO "database_name"."new_schema_name"`)
	})

	t.Run("swap with", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     NewDatabaseObjectIdentifier("database_name", "schema_name"),
			IfExists: Bool(false),
			SwapWith: NewDatabaseObjectIdentifier("database_name", "target_schema_name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA "database_name"."schema_name" SWAP WITH "database_name"."target_schema_name"`)
	})

	t.Run("set options", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: NewDatabaseObjectIdentifier("database_name", "schema_name"),
			Set: &SchemaSet{
				DataRetentionTimeInDays:    Int(3),
				MaxDataExtensionTimeInDays: Int(2),
				DefaultDDLCollation:        String("en_US-trim"),
				Comment:                    String("comment"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA "database_name"."schema_name" SET DATA_RETENTION_TIME_IN_DAYS = 3, MAX_DATA_EXTENSION_TIME_IN_DAYS = 2, DEFAULT_DDL_COLLATION = 'en_US-trim', COMMENT = 'comment'`)
	})

	t.Run("set tags", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: NewDatabaseObjectIdentifier("database_name", "schema_name"),
			SetTag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA "database_name"."schema_name" SET TAG "tag1" = 'value1', "tag2" = 'value2'`)
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: NewDatabaseObjectIdentifier("database_name", "schema_name"),
			UnsetTag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag1"),
				NewAccountObjectIdentifier("tag2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA "database_name"."schema_name" UNSET TAG "tag1", "tag2"`)
	})

	t.Run("unset options", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: NewDatabaseObjectIdentifier("database_name", "schema_name"),
			Unset: &SchemaUnset{
				DataRetentionTimeInDays:    Bool(true),
				MaxDataExtensionTimeInDays: Bool(true),
				DefaultDDLCollation:        Bool(true),
				Comment:                    Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA "database_name"."schema_name" UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, DEFAULT_DDL_COLLATION, COMMENT`)
	})

	t.Run("enable managed access", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:                NewDatabaseObjectIdentifier("database_name", "schema_name"),
			EnableManagedAccess: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA "database_name"."schema_name" ENABLE MANAGED ACCESS`)
	})

	t.Run("disable managed access", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:                 NewDatabaseObjectIdentifier("database_name", "schema_name"),
			DisableManagedAccess: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA "database_name"."schema_name" DISABLE MANAGED ACCESS`)
	})
}

func TestSchemasDrop(t *testing.T) {
	t.Run("cascade", func(t *testing.T) {
		opts := &DropSchemaOptions{
			IfExists: Bool(true),
			name:     NewDatabaseObjectIdentifier("database_name", "schema_name"),
			Cascade:  Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SCHEMA IF EXISTS "database_name"."schema_name" CASCADE`)
	})

	t.Run("restrict", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name:     NewDatabaseObjectIdentifier("database_name", "schema_name"),
			Restrict: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SCHEMA "database_name"."schema_name" RESTRICT`)
	})
}

func TestSchemasUndrop(t *testing.T) {
	opts := &undropSchemaOptions{
		name: NewDatabaseObjectIdentifier("database_name", "schema_name"),
	}
	assertOptsValidAndSQLEquals(t, opts, `UNDROP SCHEMA "database_name"."schema_name"`)
}

func TestSchemasDescribe(t *testing.T) {
	opts := &describeSchemaOptions{
		name: NewDatabaseObjectIdentifier("database_name", "schema_name"),
	}
	assertOptsValidAndSQLEquals(t, opts, `DESCRIBE SCHEMA "database_name"."schema_name"`)
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
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY LIKE 'schema_pattern'`)
	})

	t.Run("in account", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			In: &SchemaIn{
				Account: Bool(true),
				Name:    NewAccountObjectIdentifier("account_name"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY IN ACCOUNT "account_name"`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			In: &SchemaIn{
				Database: Bool(true),
				Name:     NewAccountObjectIdentifier("database_name"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY IN DATABASE "database_name"`)
	})

	t.Run("starts with", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:      Bool(true),
			History:    Bool(true),
			StartsWith: String("schema_pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY STARTS WITH 'schema_pattern'`)
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
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY LIMIT 3 FROM 'name_string'`)
	})
}
