package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestAccountCreate(t *testing.T) {
	t.Run("simplest case", func(t *testing.T) {
		opts := &AccountCreateOptions{
			name: AccountObjectIdentifier{
				name: "newaccount",
			},
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
		opts := &AccountCreateOptions{
			name: AccountObjectIdentifier{
				name: "newaccount",
			},
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
		opts := &AccountAlterOptions{
			Set: &AccountSet{
				ClientEncryptionKeySize:       Int(20),
				PreventUnloadToInternalStages: Bool(true),
				MaxDataExtensionTimeInDays:    Int(30),
				JsonIndent:                    Int(40),
				TimestampOutputFormat:         String("hello"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET CLIENT_ENCRYPTION_KEY_SIZE = 20 PREVENT_UNLOAD_TO_INTERNAL_STAGES = true MAX_DATA_EXTENSION_TIME_IN_DAYS = 30 JSON_INDENT = 40 TIMESTAMP_OUTPUT_FORMAT = 'hello'`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset params", func(t *testing.T) {
		opts := &AccountAlterOptions{
			Unset: &AccountUnset{
				InitialReplicationSizeLimitInTb: Bool(true),
				SsoLoginPage:                    Bool(true),
				DefaultDdlCollation:             Bool(true),
				SimulatedDataSharingConsumer:    Bool(true),
				Timezone:                        Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT UNSET INITIAL_REPLICATION_SIZE_LIMIT_IN_TB,SSO_LOGIN_PAGE,DEFAULT_DDL_COLLATION,SIMULATED_DATA_SHARING_CONSUMER,TIMEZONE`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set resource monitor", func(t *testing.T) {
		opts := &AccountAlterOptions{
			ResourceMonitor: NewAccountObjectIdentifier("mymonitor"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET RESOURCE_MONITOR = "mymonitor"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set password policy", func(t *testing.T) {
		opts := &AccountAlterOptions{
			PasswordPolicy: NewSchemaObjectIdentifier("db", "schema", "passpol"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET PASSWORD POLICY "db"."schema"."passpol"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set session policy", func(t *testing.T) {
		opts := &AccountAlterOptions{
			SessionPolicy: NewSchemaObjectIdentifier("db", "schema", "sesspol"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET SESSION POLICY "db"."schema"."sesspol"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset password policy", func(t *testing.T) {
		opts := &AccountAlterOptions{
			UnsetPasswordPolicy: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT UNSET PASSWORD POLICY`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset session policy", func(t *testing.T) {
		opts := &AccountAlterOptions{
			UnsetSessionPolicy: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT UNSET SESSION POLICY`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &AccountAlterOptions{
			SetTag: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db", "schema", "tag1"),
					Value: "v1",
				},
				{
					Name:  NewSchemaObjectIdentifier("db", "schema", "tag2"),
					Value: "v2",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT SET TAG "db"."schema"."tag1" = 'v1',"db"."schema"."tag2" = 'v2'`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &AccountAlterOptions{
			UnsetTag: []ObjectIdentifier{
				NewSchemaObjectIdentifier("db", "schema", "tag1"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT UNSET TAG "db"."schema"."tag1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("rename", func(t *testing.T) {
		oldName := NewAccountObjectIdentifier("oldname")
		opts := &AccountAlterOptions{
			Name:       &oldName,
			NewName:    NewAccountObjectIdentifier("newname"),
			SaveOldURL: Bool(false),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT "oldname" RENAME TO "newname" SAVE_OLD_URL = false`
		assert.Equal(t, expected, actual)
	})

	t.Run("drop old url", func(t *testing.T) {
		oldName := NewAccountObjectIdentifier("oldname")
		opts := &AccountAlterOptions{
			Name:       &oldName,
			DropOldURL: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ACCOUNT "oldname" DROP OLD URL`
		assert.Equal(t, expected, actual)
	})
}

func TestAccountShow(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		opts := &AccountShowOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW ORGANIZATION ACCOUNTS`
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &AccountShowOptions{
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
