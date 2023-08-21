package sdk

import (
	"testing"
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
		assertOptsValidAndSQLEquals(t, opts, `GRANT MONITOR USAGE, APPLY TAG ON ACCOUNT TO ROLE "role1" WITH GRANT OPTION`)
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
		assertOptsValidAndSQLEquals(t, opts, `GRANT ALL PRIVILEGES ON DATABASE "db1" TO ROLE "role1"`)
	})
	t.Run("on schema", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON SCHEMA "db1"."schema1" TO ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE "db1" TO ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE "db1" TO ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON TABLE "db1"."schema1"."table1" TO ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN DATABASE "db1" TO ROLE "role1"`)
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
						InSchema:         Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN SCHEMA "db1"."schema1" TO ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `REVOKE MONITOR USAGE, APPLY TAG ON ACCOUNT FROM ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `REVOKE ALL PRIVILEGES ON DATABASE "db1" FROM ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE DATABASE ROLE, MODIFY ON DATABASE "db1" FROM ROLE "role1"`)
	})
	t.Run("on schema", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA "db1"."schema1" FROM ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE "db1" FROM ROLE "role1" RESTRICT`)
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
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE "db1" FROM ROLE "role1" CASCADE`)
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
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON TABLE "db1"."schema1"."table1" FROM ROLE "role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE "db1" FROM ROLE "role1"`)
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
						InSchema:         Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA "db1"."schema1" FROM ROLE "role1"`)
	})
}

func TestGrants_GrantPrivilegesToDatabaseRole(t *testing.T) {
	// TODO: validation tests

	t.Run("on database", func(t *testing.T) {
		dbId := NewAccountObjectIdentifier("db1")
		opts := &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			},
			on: &DatabaseRoleGrantOn{
				Database: &dbId,
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE SCHEMA ON DATABASE "db1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on schema", func(t *testing.T) {
		opts := &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &DatabaseRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
				},
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON SCHEMA "db1"."schema1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on all schemas in database", func(t *testing.T) {
		opts := &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &DatabaseRoleGrantOn{
				Schema: &GrantOnSchema{
					AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE "db1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on all future schemas in database", func(t *testing.T) {
		opts := &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &DatabaseRoleGrantOn{
				Schema: &GrantOnSchema{
					FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE "db1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       NewSchemaObjectIdentifier("db1", "schema1", "table1"),
					},
				},
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON TABLE "db1"."schema1"."table1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
					},
				},
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN DATABASE "db1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InSchema:         Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
					},
				},
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN SCHEMA "db1"."schema1" TO DATABASE ROLE "db1"."role1"`)
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
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON DATABASE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Schema: otherID,
			},
			to: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON TABLE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			to: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON ALL TABLES IN SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON VIEW %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON DATABASE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Schema: otherID,
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON TABLE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON ALL TABLES IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON VIEW %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on all views", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			privilege: ObjectPrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				View: &OnView{
					AllInSchema: otherID,
				},
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON ALL VIEWS IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestGrantShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &ShowGrantOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS")
	})

	t.Run("on account", func(t *testing.T) {
		opts := &ShowGrantOptions{
			On: &ShowGrantsOn{
				Account: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS ON ACCOUNT")
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
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS ON DATABASE %s", dbID.FullyQualifiedName())
	})

	t.Run("to role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("to user", func(t *testing.T) {
		userID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				User: userID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO USER %s", userID.FullyQualifiedName())
	})

	t.Run("to share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Share: shareID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO SHARE %s", shareID.FullyQualifiedName())
	})

	t.Run("of role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Role: roleID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS OF ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("of share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Share: shareID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS OF SHARE %s", shareID.FullyQualifiedName())
	})
}
