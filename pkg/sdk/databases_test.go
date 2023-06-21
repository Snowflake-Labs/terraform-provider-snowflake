package sdk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabasesCreate(t *testing.T) {
	t.Run("clone", func(t *testing.T) {
		opts := &CreateDatabaseOptions{
			Clone: &Clone{
				SourceObject: NewAccountObjectIdentifier("db1"),
				At: &TimeTravel{
					Timestamp: Pointer(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE DATABASE CLONE "db1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete", func(t *testing.T) {
		opts := &CreateDatabaseOptions{
			OrReplace:                  Bool(true),
			Transient:                  Bool(true),
			Comment:                    String("comment"),
			DataRetentionTimeInDays:    Int(1),
			MaxDataExtensionTimeInDays: Int(1),
			Tag: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
					Value: "v1",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE TRANSIENT DATABASE DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 COMMENT = 'comment' TAG ("db1"."schema1"."tag1" = 'v1')`
		assert.Equal(t, expected, actual)
	})
}

func TestDatabasesCreateShared(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		databaseID := NewAccountObjectIdentifier("db1")
		opts := &CreateSharedDatabaseOptions{
			name:      databaseID,
			fromShare: NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account1"), NewAccountObjectIdentifier("db1")),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE DATABASE "db1" FROM SHARE account1."db1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with comment", func(t *testing.T) {
		databaseID := NewAccountObjectIdentifier("db1")
		opts := &CreateSharedDatabaseOptions{
			name:      databaseID,
			fromShare: NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account1"), NewAccountObjectIdentifier("db1")),
			Comment:   String("comment"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE DATABASE "db1" FROM SHARE account1."db1" COMMENT = 'comment'`
		assert.Equal(t, expected, actual)
	})
}

func TestDatabasesCreateSecondary(t *testing.T) {
	opts := &CreateSecondaryDatabaseOptions{
		name:                    NewAccountObjectIdentifier("db1"),
		primaryDatabase:         NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account1"), NewAccountObjectIdentifier("db1")),
		DataRetentionTimeInDays: Int(1),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `CREATE DATABASE "db1" AS REPLICA OF account1."db1" DATA_RETENTION_TIME_IN_DAYS = 1`

	assert.Equal(t, expected, actual)
}

func TestDatabasesDrop(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		opts := &DropDatabaseOptions{
			name: NewAccountObjectIdentifier("db1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP DATABASE "db1"`
		assert.Equal(t, expected, actual)
	})
}

func TestDatabasesUndrop(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		opts := &undropDatabaseOptions{
			name: NewAccountObjectIdentifier("db1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `UNDROP DATABASE "db1"`
		assert.Equal(t, expected, actual)
	})
}

func TestDatabasesDescribe(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		opts := &describeDatabaseOptions{
			name: NewAccountObjectIdentifier("db1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE DATABASE "db1"`
		assert.Equal(t, expected, actual)
	})
}

func TestDatabasesAlter(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		opts := &AlterDatabaseOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("db1"),
			NewName:  NewAccountObjectIdentifier("db2"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE IF EXISTS "db1" RENAME TO "db2"`
		assert.Equal(t, expected, actual)
	})

	t.Run("swap with", func(t *testing.T) {
		opts := &AlterDatabaseOptions{
			name:     NewAccountObjectIdentifier("db1"),
			SwapWith: NewAccountObjectIdentifier("db2"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" SWAP WITH "db2"`
		assert.Equal(t, expected, actual)
	})

	t.Run("swap with", func(t *testing.T) {
		opts := &AlterDatabaseOptions{
			name:     NewAccountObjectIdentifier("db1"),
			SwapWith: NewAccountObjectIdentifier("db2"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" SWAP WITH "db2"`
		assert.Equal(t, expected, actual)
	})

	t.Run("set comment and retention time in days", func(t *testing.T) {
		opts := &AlterDatabaseOptions{
			name: NewAccountObjectIdentifier("db1"),
			Set: &DatabaseSet{
				DataRetentionTimeInDays: Int(1),
				Comment:                 String("comment"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" SET DATA_RETENTION_TIME_IN_DAYS = 1, COMMENT = 'comment'`
		assert.Equal(t, expected, actual)
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := &AlterDatabaseOptions{
			name: NewAccountObjectIdentifier("db1"),
			Unset: &DatabaseUnset{
				Comment: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" UNSET COMMENT`
		assert.Equal(t, expected, actual)
	})
}

func TestDatabasesAlterReplication(t *testing.T) {
	t.Run("enable replication", func(t *testing.T) {
		opts := &AlterDatabaseReplicationOptions{
			name: NewAccountObjectIdentifier("db1"),
			EnableReplication: &EnableReplication{
				ToAccounts: []AccountIdentifier{
					NewAccountIdentifierFromAccountLocator("account1"),
				},
				IgnoreEditionCheck: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" ENABLE REPLICATION TO ACCOUNTS "account1" IGNORE EDITION CHECK`
		assert.Equal(t, expected, actual)
	})

	t.Run("disable replication", func(t *testing.T) {
		opts := &AlterDatabaseReplicationOptions{
			name: NewAccountObjectIdentifier("db1"),
			DisableReplication: &DisableReplication{
				ToAccounts: []AccountIdentifier{
					NewAccountIdentifierFromAccountLocator("account1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" DISABLE REPLICATION TO ACCOUNTS "account1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("refresh", func(t *testing.T) {
		opts := &AlterDatabaseReplicationOptions{
			name:    NewAccountObjectIdentifier("db1"),
			Refresh: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" REFRESH`
		assert.Equal(t, expected, actual)
	})
}

func TestDatabasesAlterFailover(t *testing.T) {
	t.Run("enable failover", func(t *testing.T) {
		opts := &AlterDatabaseFailoverOptions{
			name: NewAccountObjectIdentifier("db1"),
			EnableFailover: &EnableFailover{
				ToAccounts: []AccountIdentifier{
					NewAccountIdentifierFromAccountLocator("account1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" ENABLE FAILOVER TO ACCOUNTS "account1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("disable failover", func(t *testing.T) {
		opts := &AlterDatabaseFailoverOptions{
			name: NewAccountObjectIdentifier("db1"),
			DisableFailover: &DisableFailover{
				ToAccounts: []AccountIdentifier{
					NewAccountIdentifierFromAccountLocator("account1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" DISABLE FAILOVER TO ACCOUNTS "account1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("primary", func(t *testing.T) {
		opts := &AlterDatabaseFailoverOptions{
			name:    NewAccountObjectIdentifier("db1"),
			Primary: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DATABASE "db1" PRIMARY`
		assert.Equal(t, expected, actual)
	})
}

func TestDatabasesShow(t *testing.T) {
	t.Run("without show options", func(t *testing.T) {
		opts := &ShowDatabasesOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW DATABASES`
		assert.Equal(t, expected, actual)
	})

	t.Run("terse", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			Terse: Bool(true),
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW TERSE DATABASES`
		assert.Equal(t, expected, actual)
	})

	t.Run("history", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			History: Bool(true),
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW DATABASES HISTORY`
		assert.Equal(t, expected, actual)
	})

	t.Run("like", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			Like: &Like{
				Pattern: String("db1"),
			},
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW DATABASES LIKE 'db1'`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			Terse:   Bool(true),
			History: Bool(true),
			Like: &Like{
				Pattern: String("db2"),
			},
			LimitFrom: &LimitFrom{
				Rows: Int(1),
				From: String("db1"),
			},
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW TERSE DATABASES HISTORY LIKE 'db2' LIMIT 1 FROM 'db1'`
		assert.Equal(t, expected, actual)
	})
}
