package sdk

import (
	"errors"
	"fmt"
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

	t.Run("on account object - external volume", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					ExternalVolume: Pointer(NewAccountObjectIdentifier("ex volume")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ALL PRIVILEGES ON EXTERNAL VOLUME "ex volume" TO ROLE "role1"`)
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
	dbId := NewAccountObjectIdentifier("db1")

	defaultGrantsForDb := func() *GrantPrivilegesToDatabaseRoleOptions {
		return &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			},
			on: &DatabaseRoleGrantOn{
				Database: &dbId,
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
	}

	defaultGrantsForSchema := func() *GrantPrivilegesToDatabaseRoleOptions {
		return &GrantPrivilegesToDatabaseRoleOptions{
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
	}

	defaultGrantsForSchemaObject := func() *GrantPrivilegesToDatabaseRoleOptions {
		return &GrantPrivilegesToDatabaseRoleOptions{
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
	}

	t.Run("validation: nil privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GrantPrivilegesToDatabaseRoleOptions", "privileges"))
	})

	t.Run("validation: no privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: too many privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			SchemaPrivileges:   []SchemaPrivilege{SchemaPrivilegeCreateAlert},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GrantPrivilegesToDatabaseRoleOptions", "on"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = &DatabaseRoleGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: too many ons set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = &DatabaseRoleGrantOn{
			Database: &dbId,
			Schema: &GrantOnSchema{
				Schema: Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: grant on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchema", "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase"))
	})

	t.Run("validation: grant on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObject", "SchemaObject", "All", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				All: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: unsupported database privilege", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges.DatabasePrivileges = []AccountObjectPrivilege{AccountObjectPrivilegeCreateDatabaseRole}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("privilege CREATE DATABASE ROLE is not allowed"))
	})

	t.Run("on database", func(t *testing.T) {
		opts := defaultGrantsForDb()
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE SCHEMA ON DATABASE "db1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON SCHEMA "db1"."schema1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on all schemas in database", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{
			AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE "db1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on all future schemas in database", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{
			FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE "db1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON TABLE "db1"."schema1"."table1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN DATABASE "db1" TO DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN SCHEMA "db1"."schema1" TO DATABASE ROLE "db1"."role1"`)
	})
}

func TestGrants_RevokePrivilegesFromDatabaseRoleRole(t *testing.T) {
	dbId := NewAccountObjectIdentifier("db1")

	defaultGrantsForDb := func() *RevokePrivilegesFromDatabaseRoleOptions {
		return &RevokePrivilegesFromDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			},
			on: &DatabaseRoleGrantOn{
				Database: &dbId,
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
	}

	defaultGrantsForSchema := func() *RevokePrivilegesFromDatabaseRoleOptions {
		return &RevokePrivilegesFromDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &DatabaseRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
				},
			},
			databaseRole: NewDatabaseObjectIdentifier("db1", "role1"),
		}
	}

	defaultGrantsForSchemaObject := func() *RevokePrivilegesFromDatabaseRoleOptions {
		return &RevokePrivilegesFromDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
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
	}

	t.Run("validation: nil privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("RevokePrivilegesFromDatabaseRoleOptions", "privileges"))
	})

	t.Run("validation: no privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: too many privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			SchemaPrivileges:   []SchemaPrivilege{SchemaPrivilegeCreateAlert},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: nil on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("RevokePrivilegesFromDatabaseRoleOptions", "on"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = &DatabaseRoleGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: too many ons set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = &DatabaseRoleGrantOn{
			Database: &dbId,
			Schema: &GrantOnSchema{
				Schema: Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: grant on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchema", "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase"))
	})

	t.Run("validation: grant on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObject", "SchemaObject", "All", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				All: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: unsupported database privilege", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges.DatabasePrivileges = []AccountObjectPrivilege{AccountObjectPrivilegeCreateDatabaseRole}
		assertOptsInvalidJoinedErrors(t, opts, errors.New("privilege CREATE DATABASE ROLE is not allowed"))
	})

	t.Run("on database", func(t *testing.T) {
		opts := defaultGrantsForDb()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE SCHEMA ON DATABASE "db1" FROM DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA "db1"."schema1" FROM DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on all schemas in database + restrict", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{
			AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
		}
		opts.Restrict = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE "db1" FROM DATABASE ROLE "db1"."role1" RESTRICT`)
	})

	t.Run("on all future schemas in database + cascade", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{
			FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
		}
		opts.Cascade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE "db1" FROM DATABASE ROLE "db1"."role1" CASCADE`)
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON TABLE "db1"."schema1"."table1" FROM DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE "db1" FROM DATABASE ROLE "db1"."role1"`)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(NewDatabaseObjectIdentifier("db1", "schema1")),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA "db1"."schema1" FROM DATABASE ROLE "db1"."role1"`)
	})
}

func TestGrantPrivilegeToShare(t *testing.T) {
	id := RandomAccountObjectIdentifier()
	t.Run("on database", func(t *testing.T) {
		otherID := RandomAccountObjectIdentifier()
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
		otherID := RandomDatabaseObjectIdentifier()
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
		otherID := RandomSchemaObjectIdentifier()
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
		otherID := RandomDatabaseObjectIdentifier()
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
		otherID := RandomSchemaObjectIdentifier()
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
	id := RandomAccountObjectIdentifier()
	t.Run("on database", func(t *testing.T) {
		otherID := RandomAccountObjectIdentifier()
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
		otherID := RandomDatabaseObjectIdentifier()
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
		otherID := RandomSchemaObjectIdentifier()
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
		otherID := RandomDatabaseObjectIdentifier()
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
		otherID := RandomSchemaObjectIdentifier()
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
		otherID := RandomDatabaseObjectIdentifier()
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

func TestGrants_GrantOwnership(t *testing.T) {
	dbId := NewAccountObjectIdentifier("db1")
	schemaId := NewDatabaseObjectIdentifier("db1", "schema1")
	roleId := NewAccountObjectIdentifier("role1")
	databaseRoleId := NewDatabaseObjectIdentifier("db1", "role1")
	tableId := NewSchemaObjectIdentifier("db1", "schema1", "table1")

	defaultOpts := func() *GrantOwnershipOptions {
		return &GrantOwnershipOptions{
			On: OwnershipGrantOn{
				Object: &Object{
					ObjectType: ObjectTypeTable,
					Name:       tableId,
				},
			},
			To: OwnershipGrantTo{
				AccountRoleName: Pointer(roleId),
			},
		}
	}

	t.Run("validation: grant on empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantOn", "Object", "AllIn", "Future"))
	})

	t.Run("validation: grant on too many", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Object: &Object{
				ObjectType: ObjectTypeTable,
				Name:       tableId,
			},
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantOn", "Object", "AllIn", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			All: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant to empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantTo", "databaseRoleName", "accountRoleName"))
	})

	t.Run("validation: grant to role and database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
			AccountRoleName:  Pointer(roleId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantTo", "databaseRoleName", "accountRoleName"))
	})

	t.Run("on schema object to role", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON TABLE %s TO ROLE %s`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on schema object to database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON FUTURE TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on all schema objects in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			All: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON ALL TABLES IN SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on schema object with current grants", func(t *testing.T) {
		opts := defaultOpts()
		opts.CurrentGrants = &OwnershipCurrentGrants{
			OutboundPrivileges: Copy,
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON TABLE %s TO ROLE %s COPY CURRENT GRANTS`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName())
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
		dbID := RandomAccountObjectIdentifier()
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
		roleID := RandomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("to user", func(t *testing.T) {
		userID := RandomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				User: userID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO USER %s", userID.FullyQualifiedName())
	})

	t.Run("to share", func(t *testing.T) {
		shareID := RandomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Share: shareID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO SHARE %s", shareID.FullyQualifiedName())
	})

	t.Run("of role", func(t *testing.T) {
		roleID := RandomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Role: roleID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS OF ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("of share", func(t *testing.T) {
		shareID := RandomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Share: shareID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS OF SHARE %s", shareID.FullyQualifiedName())
	})
}
