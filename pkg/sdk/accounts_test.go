package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestAccountCreate(t *testing.T) {
	t.Run("simplest case", func(t *testing.T) {
		opts := &CreateAccountOptions{
			name:          NewAccountObjectIdentifier("newaccount"),
			AdminName:     "someadmin",
			AdminPassword: String("v3rys3cr3t"),
			Email:         "admin@example.com",
			Edition:       EditionBusinessCritical,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE ACCOUNT "newaccount" ADMIN_NAME = 'someadmin' ADMIN_PASSWORD = 'v3rys3cr3t' EMAIL = 'admin@example.com' EDITION = BUSINESS_CRITICAL`
		assert.Equal(t, expected, actual)
	})

	t.Run("every option", func(t *testing.T) {
		opts := &CreateAccountOptions{
			name:               NewAccountObjectIdentifier("newaccount"),
			AdminName:          "someadmin",
			AdminRSAPublicKey:  String("s3cr3tk3y"),
			FirstName:          String("Ad"),
			LastName:           String("Min"),
			Email:              "admin@example.com",
			MustChangePassword: Bool(true),
			Edition:            EditionBusinessCritical,
			RegionGroup:        String("groupid"),
			Region:             String("regionid"),
			Comment:            String("Test account"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE ACCOUNT "newaccount" ADMIN_NAME = 'someadmin' ADMIN_RSA_PUBLIC_KEY = 's3cr3tk3y' FIRST_NAME = 'Ad' LAST_NAME = 'Min' EMAIL = 'admin@example.com' MUST_CHANGE_PASSWORD = true EDITION = BUSINESS_CRITICAL REGION_GROUP = 'groupid' REGION = 'regionid' COMMENT = 'Test account'`
		assert.Equal(t, expected, actual)
	})
}

func TestAccountAlter(t *testing.T) {
	t.Run("with set params", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				Parameters: &AccountLevelParameters{
					AccountParameters: &AccountParameters{
						ClientEncryptionKeySize:       Int(128),
						PreventUnloadToInternalStages: Bool(true),
					},
					SessionParameters: &SessionParameters{
						JSONIndent: Int(16),
					},
					ObjectParameters: &ObjectParameters{
						MaxDataExtensionTimeInDays: Int(30),
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET CLIENT_ENCRYPTION_KEY_SIZE = 128, PREVENT_UNLOAD_TO_INTERNAL_STAGES = true, JSON_INDENT = 16, MAX_DATA_EXTENSION_TIME_IN_DAYS = 30`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset params", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				Parameters: &AccountLevelParametersUnset{
					AccountParameters: &AccountParametersUnset{
						InitialReplicationSizeLimitInTB: Bool(true),
						SSOLoginPage:                    Bool(true),
					},
					SessionParameters: &SessionParametersUnset{
						SimulatedDataSharingConsumer: Bool(true),
						Timezone:                     Bool(true),
					},
					ObjectParameters: &ObjectParametersUnset{
						DefaultDDLCollation: Bool(true),
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT UNSET INITIAL_REPLICATION_SIZE_LIMIT_IN_TB, SSO_LOGIN_PAGE, SIMULATED_DATA_SHARING_CONSUMER, TIMEZONE, DEFAULT_DDL_COLLATION`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set resource monitor", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				ResourceMonitor: NewAccountObjectIdentifier("mymonitor"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET RESOURCE_MONITOR = "mymonitor"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set password policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				PasswordPolicy: NewSchemaObjectIdentifier("db", "schema", "passpol"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET PASSWORD POLICY "db"."schema"."passpol"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set session policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				SessionPolicy: NewSchemaObjectIdentifier("db", "schema", "sesspol"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET SESSION POLICY "db"."schema"."sesspol"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset password policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				PasswordPolicy: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT UNSET PASSWORD POLICY`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset session policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				SessionPolicy: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT UNSET SESSION POLICY`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				Tag: []TagAssociation{
					{
						Name:  NewSchemaObjectIdentifier("db", "schema", "tag1"),
						Value: "v1",
					},
					{
						Name:  NewSchemaObjectIdentifier("db", "schema", "tag2"),
						Value: "v2",
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET TAG "db"."schema"."tag1" = 'v1', "db"."schema"."tag2" = 'v2'`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				Tag: []ObjectIdentifier{
					NewSchemaObjectIdentifier("db", "schema", "tag1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT UNSET TAG "db"."schema"."tag1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("rename", func(t *testing.T) {
		oldName := NewAccountObjectIdentifier("oldname")
		newName := NewAccountObjectIdentifier("newname")
		opts := &AlterAccountOptions{
			Rename: &AccountRename{
				Name:       oldName,
				NewName:    newName,
				SaveOldURL: Bool(false),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT "oldname" RENAME TO "newname" SAVE_OLD_URL = false`
		assert.Equal(t, expected, actual)
	})

	t.Run("drop old url", func(t *testing.T) {
		oldName := NewAccountObjectIdentifier("oldname")
		opts := &AlterAccountOptions{
			Drop: &AccountDrop{
				Name:   oldName,
				OldURL: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT "oldname" DROP OLD URL`
		assert.Equal(t, expected, actual)
	})
}

func TestAccountShow(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		opts := &ShowAccountOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW ORGANIZATION ACCOUNTS`
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowAccountOptions{
			Like: &Like{
				Pattern: String("myaccount"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW ORGANIZATION ACCOUNTS LIKE 'myaccount'`
		assert.Equal(t, expected, actual)
	})
}
