package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFailoverGroupsCreate(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		opts := &CreateFailoverGroupOptions{
			IfNotExists: Bool(true),
			name:        NewAccountObjectIdentifier("fg1"),
			objectTypes: []PluralObjectType{
				PluralObjectTypeShares,
				PluralObjectTypeDatabases,
			},
			AllowedDatabases: []AccountObjectIdentifier{
				NewAccountObjectIdentifier("db1"),
			},
			AllowedShares: []AccountObjectIdentifier{
				NewAccountObjectIdentifier("share1"),
			},
			allowedAccounts: []AccountIdentifier{
				NewAccountIdentifier("MY_ORG", "MY_ACCOUNT"),
			},
			IgnoreEditionCheck:  Bool(true),
			ReplicationSchedule: String("10 MINUTE"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE FAILOVER GROUP IF NOT EXISTS "fg1" OBJECT_TYPES = SHARES, DATABASES ALLOWED_DATABASES = "db1" ALLOWED_SHARES = "share1" ALLOWED_ACCOUNTS = "MY_ORG.MY_ACCOUNT" IGNORE EDITION CHECK REPLICATION_SCHEDULE = '10 MINUTE'`
		assert.Equal(t, expected, actual)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := &CreateFailoverGroupOptions{
			IfNotExists: Bool(true),
			name:        NewAccountObjectIdentifier("fg1"),
			objectTypes: []PluralObjectType{
				PluralObjectTypeRoles,
			},
			allowedAccounts: []AccountIdentifier{
				NewAccountIdentifier("MY_ORG", "MY_ACCOUNT"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE FAILOVER GROUP IF NOT EXISTS "fg1" OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = "MY_ORG.MY_ACCOUNT"`
		assert.Equal(t, expected, actual)
	})
}

func TestCreateSecondaryReplicationGroup(t *testing.T) {
	opts := &CreateSecondaryReplicationGroupOptions{
		IfNotExists:          Bool(true),
		name:                 NewAccountObjectIdentifier("fg1"),
		primaryFailoverGroup: NewExternalObjectIdentifierFromFullyQualifiedName("myorg.myaccount.fg1"),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `CREATE FAILOVER GROUP IF NOT EXISTS "fg1" AS REPLICA OF myorg.myaccount."fg1"`
	assert.Equal(t, expected, actual)
}

func TestFailoverGroupAlterSource(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name:    NewAccountObjectIdentifier("fg1"),
			NewName: NewAccountObjectIdentifier("myfg1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP "fg1" RENAME TO "myfg1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("set object types and replication schedule", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: NewAccountObjectIdentifier("fg1"),
			Set: &FailoverGroupSet{
				ObjectTypes:         []PluralObjectType{PluralObjectTypeShares},
				ReplicationSchedule: String("10 MINUTE"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP "fg1" SET OBJECT_TYPES = SHARES REPLICATION_SCHEDULE = '10 MINUTE'`
		assert.Equal(t, expected, actual)
	})

	t.Run("add database account object", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			Add: &FailoverGroupAdd{
				AllowedDatabases: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("db1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP ADD "db1" TO ALLOWED_DATABASES`
		assert.Equal(t, expected, actual)
	})

	t.Run("remove database account object", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			Remove: &FailoverGroupRemove{
				AllowedDatabases: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("db1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP REMOVE "db1" FROM ALLOWED_DATABASES`
		assert.Equal(t, expected, actual)
	})

	t.Run("add share account object", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			Add: &FailoverGroupAdd{
				AllowedShares: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("share1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP ADD "share1" TO ALLOWED_SHARES`
		assert.Equal(t, expected, actual)
	})

	t.Run("remove share account object", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			Remove: &FailoverGroupRemove{
				AllowedShares: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("share1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP REMOVE "share1" FROM ALLOWED_SHARES`
		assert.Equal(t, expected, actual)
	})

	t.Run("move shares to another failover group", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: NewAccountObjectIdentifier("fg1"),
			Move: &FailoverGroupMove{
				Shares: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("share1"),
				},
				To: NewAccountObjectIdentifier("fg2"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP "fg1" MOVE SHARES "share1" TO FAILOVER GROUP "fg2"`
		assert.Equal(t, expected, actual)
	})
}

func TestFailoverGroupsAlterTarget(t *testing.T) {
	t.Run("resume", func(t *testing.T) {
		opts := &AlterTargetFailoverGroupOptions{
			name:   NewAccountObjectIdentifier("fg1"),
			Resume: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP "fg1" RESUME`
		assert.Equal(t, expected, actual)
	})

	t.Run("primary", func(t *testing.T) {
		opts := &AlterTargetFailoverGroupOptions{
			name:    NewAccountObjectIdentifier("fg1"),
			Primary: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FAILOVER GROUP "fg1" PRIMARY`
		assert.Equal(t, expected, actual)
	})
}

func TestFailoverGroupsDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &DropFailoverGroupOptions{
			name: NewAccountObjectIdentifier("fg1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP FAILOVER GROUP "fg1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with IfExists", func(t *testing.T) {
		opts := &DropFailoverGroupOptions{
			name:     NewAccountObjectIdentifier("fg1"),
			IfExists: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP FAILOVER GROUP IF EXISTS "fg1"`
		assert.Equal(t, expected, actual)
	})
}

func TestFailoverGroupsShow(t *testing.T) {
	t.Run("without show options", func(t *testing.T) {
		opts := &ShowFailoverGroupOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW FAILOVER GROUPS`
		assert.Equal(t, expected, actual)
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &ShowFailoverGroupOptions{
			InAccount: NewAccountIdentifierFromAccountLocator("abcd123"),
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW FAILOVER GROUPS IN ACCOUNT "abcd123"`
		assert.Equal(t, expected, actual)
	})
}

func TestFailoverGroupsShowDatabases(t *testing.T) {
	opts := &showFailoverGroupDatabasesOptions{}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `SHOW DATABASES`
	assert.Equal(t, expected, actual)
}

func TestFailoverGroupsShowShares(t *testing.T) {
	opts := &showFailoverGroupSharesOptions{
		in: NewAccountObjectIdentifier("fg1"),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `SHOW SHARES IN FAILOVER GROUP "fg1"`
	assert.Equal(t, expected, actual)
}
