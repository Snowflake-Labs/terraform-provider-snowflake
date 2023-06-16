package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrantPrivilegesToAccountRole(t *testing.T) {
	t.Run("on account", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
			},
			on: &AccountRoleGrantOn{
				Account: Bool(true),
			},
			accountRole:     NewAccountObjectIdentifier("role1"),
			WithGrantOption: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT MONITOR USAGE, APPLY TAG ON ACCOUNT TO ROLE "role1" WITH GRANT OPTION`
		assert.Equal(t, expected, actual)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Database: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT ALL PRIVILEGES ON DATABASE "db1" TO ROLE "role1"`
		assert.Equal(t, expected, actual)
	})
	t.Run("on schema", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(NewSchemaIdentifier("db1", "schema1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT CREATE ALERT ON SCHEMA "db1"."schema1" TO ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on all schemas in database", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE "db1" TO ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on all future schemas in database", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE "db1" TO ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       NewSchemaObjectIdentifier("db1", "schema1", "table1"),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT APPLY ON TABLE "db1"."schema1"."table1" TO ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT APPLY ON FUTURE TABLES IN DATABASE "db1" TO ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InSchema:         Pointer(NewSchemaIdentifier("db1", "schema1")),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT APPLY ON FUTURE TABLES IN SCHEMA "db1"."schema1" TO ROLE "role1"`
		assert.Equal(t, expected, actual)
	})
}

func TestRevokePrivilegesFromAccountRole(t *testing.T) {
	t.Run("on account", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
			},
			on: &AccountRoleGrantOn{
				Account: Bool(true),
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE MONITOR USAGE, APPLY TAG ON ACCOUNT FROM ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Database: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE ALL PRIVILEGES ON DATABASE "db1" FROM ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateDatabaseRole, AccountObjectPrivilegeModify},
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Database: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE CREATE DATABASE ROLE, MODIFY ON DATABASE "db1" FROM ROLE "role1"`
		assert.Equal(t, expected, actual)
	})
	t.Run("on schema", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(NewSchemaIdentifier("db1", "schema1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA "db1"."schema1" FROM ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on all schemas in database + restrict", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
			Restrict:    Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE "db1" FROM ROLE "role1" RESTRICT`
		assert.Equal(t, expected, actual)
	})

	t.Run("on all future schemas in database + cascade", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
			Cascade:     Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE "db1" FROM ROLE "role1" CASCADE`
		assert.Equal(t, expected, actual)
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       NewSchemaObjectIdentifier("db1", "schema1", "table1"),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE SELECT, UPDATE ON TABLE "db1"."schema1"."table1" FROM ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE "db1" FROM ROLE "role1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InSchema:         Pointer(NewSchemaIdentifier("db1", "schema1")),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA "db1"."schema1" FROM ROLE "role1"`
		assert.Equal(t, expected, actual)
	})
}

func TestGrantPrivilegeToShare(t *testing.T) {
	id := randomAccountObjectIdentifier(t)
	t.Run("on database", func(t *testing.T) {
		otherID := randomAccountObjectIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Database: otherID,
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON DATABASE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Schema: otherID,
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on table", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Table: &OnTable{
					Name: otherID,
				},
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON TABLE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON ALL TABLES IN SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on view", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				View: otherID,
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON VIEW %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestRevokePrivilegeFromShare(t *testing.T) {
	id := randomAccountObjectIdentifier(t)
	t.Run("on database", func(t *testing.T) {
		otherID := randomAccountObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Database: otherID,
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON DATABASE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Schema: otherID,
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on table", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Table: &OnTable{
					Name: otherID,
				},
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON TABLE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON ALL TABLES IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on view", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				View: &OnView{
					Name: otherID,
				},
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON VIEW %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on all views", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				View: &OnView{
					AllInSchema: otherID,
				},
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON ALL VIEWS IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestGrantShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &ShowGrantOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW GRANTS"
		assert.Equal(t, expected, actual)
	})

	t.Run("on account", func(t *testing.T) {
		opts := &ShowGrantOptions{
			On: &ShowGrantsOn{
				Account: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW GRANTS ON ACCOUNT"
		assert.Equal(t, expected, actual)
	})

	t.Run("on database", func(t *testing.T) {
		dbID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       dbID,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS ON DATABASE %s", dbID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS TO ROLE %s", roleID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to user", func(t *testing.T) {
		userID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				User: userID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS TO USER %s", userID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Share: shareID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS TO SHARE %s", shareID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("of role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Role: roleID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS OF ROLE %s", roleID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("of share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Share: shareID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS OF SHARE %s", shareID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}
