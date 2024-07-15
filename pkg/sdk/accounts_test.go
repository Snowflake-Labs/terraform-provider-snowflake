package sdk

import (
	"testing"
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE ACCOUNT "newaccount" ADMIN_NAME = 'someadmin' ADMIN_PASSWORD = 'v3rys3cr3t' EMAIL = 'admin@example.com' EDITION = BUSINESS_CRITICAL`)
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE ACCOUNT "newaccount" ADMIN_NAME = 'someadmin' ADMIN_RSA_PUBLIC_KEY = 's3cr3tk3y' FIRST_NAME = 'Ad' LAST_NAME = 'Min' EMAIL = 'admin@example.com' MUST_CHANGE_PASSWORD = true EDITION = BUSINESS_CRITICAL REGION_GROUP = 'groupid' REGION = 'regionid' COMMENT = 'Test account'`)
	})

	t.Run("static password", func(t *testing.T) {
		opts := &CreateAccountOptions{
			name:               NewAccountObjectIdentifier("newaccount"),
			AdminName:          "someadmin",
			AdminPassword:      String("v3rys3cr3t"),
			FirstName:          String("Ad"),
			LastName:           String("Min"),
			Email:              "admin@example.com",
			MustChangePassword: Bool(false),
			Edition:            EditionBusinessCritical,
			RegionGroup:        String("groupid"),
			Region:             String("regionid"),
			Comment:            String("Test account"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE ACCOUNT "newaccount" ADMIN_NAME = 'someadmin' ADMIN_PASSWORD = 'v3rys3cr3t' FIRST_NAME = 'Ad' LAST_NAME = 'Min' EMAIL = 'admin@example.com' MUST_CHANGE_PASSWORD = false EDITION = BUSINESS_CRITICAL REGION_GROUP = 'groupid' REGION = 'regionid' COMMENT = 'Test account'`)
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET CLIENT_ENCRYPTION_KEY_SIZE = 128, PREVENT_UNLOAD_TO_INTERNAL_STAGES = true, JSON_INDENT = 16, MAX_DATA_EXTENSION_TIME_IN_DAYS = 30`)
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET INITIAL_REPLICATION_SIZE_LIMIT_IN_TB, SSO_LOGIN_PAGE, SIMULATED_DATA_SHARING_CONSUMER, TIMEZONE, DEFAULT_DDL_COLLATION`)
	})

	t.Run("with set resource monitor", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				ResourceMonitor: NewAccountObjectIdentifier("mymonitor"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET RESOURCE_MONITOR = "mymonitor"`)
	})

	t.Run("with set password policy", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				PasswordPolicy: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET PASSWORD POLICY %s`, id.FullyQualifiedName())
	})

	t.Run("with set session policy", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				SessionPolicy: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET SESSION POLICY %s`, id.FullyQualifiedName())
	})

	t.Run("with set authentication policy", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				AuthenticationPolicy: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET AUTHENTICATION POLICY %s`, id.FullyQualifiedName())
	})

	t.Run("with unset password policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				PasswordPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET PASSWORD POLICY`)
	})

	t.Run("with unset session policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				SessionPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET SESSION POLICY`)
	})

	t.Run("with unset authentication policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				AuthenticationPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET AUTHENTICATION POLICY`)
	})

	t.Run("with set tag", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := &AlterAccountOptions{
			SetTag: []TagAssociation{
				{
					Name:  tagId1,
					Value: "v1",
				},
				{
					Name:  tagId2,
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET TAG %s = 'v1', %s = 'v2'`, tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("with unset tag", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			UnsetTag: []ObjectIdentifier{
				id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET TAG %s`, id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT "oldname" RENAME TO "newname" SAVE_OLD_URL = false`)
	})

	t.Run("drop old url", func(t *testing.T) {
		oldName := NewAccountObjectIdentifier("oldname")
		opts := &AlterAccountOptions{
			Drop: &AccountDrop{
				Name:   oldName,
				OldURL: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT "oldname" DROP OLD URL`)
	})
}

func TestAccountShow(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		opts := &ShowAccountOptions{}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ORGANIZATION ACCOUNTS`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowAccountOptions{
			Like: &Like{
				Pattern: String("myaccount"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ORGANIZATION ACCOUNTS LIKE 'myaccount'`)
	})
}
