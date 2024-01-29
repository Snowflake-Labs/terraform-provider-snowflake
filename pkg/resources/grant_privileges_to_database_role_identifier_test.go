package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestParseGrantPrivilegesToDatabaseRoleId(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier string
		Expected   GrantPrivilegesToDatabaseRoleId
		Error      string
	}{
		{
			Name:       "grant database role on database",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnDatabaseShareGrantKind|"on-database-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnDatabaseDatabaseRoleGrantKind,
				Data: &OnDatabaseGrantData{
					DatabaseName: sdk.NewAccountObjectIdentifier("on-database-name"),
				},
			},
		},
		{
			Name:       "grant database role on database - always apply with grant option",
			Identifier: `"database-name"."database-role"|true|true|CREATE SCHEMA,USAGE,MONITOR|OnDatabaseShareGrantKind|"on-database-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  true,
				AlwaysApply:      true,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnDatabaseDatabaseRoleGrantKind,
				Data: &OnDatabaseGrantData{
					DatabaseName: sdk.NewAccountObjectIdentifier("on-database-name"),
				},
			},
		},
		{
			Name:       "grant database role on database - all privileges",
			Identifier: `"database-name"."database-role"|false|false|ALL|OnDatabaseShareGrantKind|"on-database-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				AllPrivileges:    true,
				Privileges:       nil,
				Kind:             OnDatabaseDatabaseRoleGrantKind,
				Data: &OnDatabaseGrantData{
					DatabaseName: sdk.NewAccountObjectIdentifier("on-database-name"),
				},
			},
		},
		{
			Name:       "grant database role on schema with schema name",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnSchema|"database-name"."schema-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:       OnSchemaSchemaGrantKind,
					SchemaName: sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
				},
			},
		},
		{
			Name:       "grant database role on all schemas in database",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnAllSchemasInDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:         OnAllSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
				},
			},
		},
		{
			Name:       "grant database role on future schemas in database",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnFutureSchemasInDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:         OnFutureSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
				},
			},
		},
		{
			Name:       "grant database role on schema object with on object option",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject|TABLE|"database-name"."schema-name"."table-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnObjectSchemaObjectGrantKind,
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeTable,
						Name:       sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "table-name"),
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option in database",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             InDatabaseBulkOperationGrantKind,
						Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option in schema",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InSchema|"database-name"."schema-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             InSchemaBulkOperationGrantKind,
						Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on future option",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option in database",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES|InDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             InDatabaseBulkOperationGrantKind,
						Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option in schema",
			Identifier: `"database-name"."database-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES|InSchema|"database-name"."schema-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             InSchemaBulkOperationGrantKind,
						Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
					},
				},
			},
		},
		{
			Name:       "validation: grant database role not enough parts",
			Identifier: `"database-name"."role-name"|false|false`,
			Error:      "database role identifier should hold at least 6 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnDatabaseShareGrantKind kind",
			Identifier: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnDatabaseShareGrantKind`,
			Error:      "database role identifier should hold at least 6 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnSchema kind",
			Identifier: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnAllSchemasInDatabase`,
			Error:      "database role identifier should hold at least 7 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnSchemaObject kind",
			Identifier: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject`,
			Error:      "database role identifier should hold at least 7 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnSchemaObject kind",
			Identifier: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject|TABLE`,
			Error:      "database role identifier should hold 8 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnSchemaObject.InDatabase kind",
			Identifier: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InDatabase`,
			Error:      "database role identifier should hold 9 parts",
		},
		{
			Name:       "validation: grant database role invalid DatabaseRoleGrantKind kind",
			Identifier: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|some-kind|some-data`,
			Error:      "invalid DatabaseRoleGrantKind: some-kind",
		},
		{
			Name:       "validation: grant database role invalid OnSchemaGrantKind kind",
			Identifier: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|some-kind|some-data`,
			Error:      "invalid OnSchemaGrantKind: some-kind",
		},
		{
			Name:       "validation: grant database role invalid OnSchemaObjectGrantKind kind",
			Identifier: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|some-kind|some-data`,
			Error:      "invalid OnSchemaObjectGrantKind: some-kind",
		},
		{
			Name:       "validation: grant database role empty privileges",
			Identifier: `"database-name"."database-role"|false|false||OnDatabaseShareGrantKind|"on-database-name"`,
			Error:      `invalid Privileges value: , should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all privileges`,
		},
		{
			Name:       "validation: grant database role empty with grant option",
			Identifier: `"database-name"."database-role"||false|ALL PRIVILEGES|OnDatabaseShareGrantKind|"on-database-name"`,
			Error:      `invalid WithGrantOption value: , should be either "true" or "false"`,
		},
		{
			Name:       "validation: grant database role empty always apply",
			Identifier: `"database-name"."database-role"|false||ALL PRIVILEGES|OnDatabaseShareGrantKind|"on-database-name"`,
			Error:      `invalid AlwaysApply value: , should be either "true" or "false"`,
		},
		{
			Name:       "validation: grant database role empty database role name",
			Identifier: `|false|false|ALL PRIVILEGES|OnDatabaseShareGrantKind|"on-database-name"`,
			Error:      "invalid DatabaseRoleName value: , should be a fully qualified name of database object <database_name>.<name>",
		},
		{
			Name:       "validation: grant database role empty type",
			Identifier: `"database-name"."database-role"|false|false|ALL PRIVILEGES||"on-database-name"`,
			Error:      "invalid DatabaseRoleGrantKind: ",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			id, err := ParseGrantPrivilegesToDatabaseRoleId(tt.Identifier)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, id)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}

func TestGrantPrivilegesToDatabaseRoleIdString(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier GrantPrivilegesToDatabaseRoleId
		Expected   string
		Error      string
	}{
		{
			Name: "grant database role on database",
			Identifier: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "role-name"),
				WithGrantOption:  true,
				AllPrivileges:    true,
				Kind:             OnDatabaseDatabaseRoleGrantKind,
				AlwaysApply:      true,
				Data: &OnDatabaseGrantData{
					DatabaseName: sdk.NewAccountObjectIdentifier("database-name"),
				},
			},
			Expected: `"database-name"."role-name"|true|true|ALL|OnDatabaseShareGrantKind|"database-name"`,
		},
		{
			Name: "grant database role on schema on schema",
			Identifier: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "role-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:       OnSchemaSchemaGrantKind,
					SchemaName: sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
				},
			},
			Expected: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnSchema|"database-name"."schema-name"`,
		},
		{
			Name: "grant database role on all schemas in database",
			Identifier: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "role-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:         OnAllSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name")),
				},
			},
			Expected: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnAllSchemasInDatabase|"database-name"`,
		},
		{
			Name: "grant database role on future schemas in database",
			Identifier: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "role-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:         OnFutureSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name")),
				},
			},
			Expected: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnFutureSchemasInDatabase|"database-name"`,
		},
		{
			Name: "grant database role on schema object on object",
			Identifier: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "role-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnObjectSchemaObjectGrantKind,
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeTable,
						Name:       sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "table-name"),
					},
				},
			},
			Expected: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject|TABLE|"database-name"."schema-name"."table-name"`,
		},
		{
			Name: "grant database role on schema object on all tables in database",
			Identifier: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "role-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: sdk.PluralObjectTypeTables,
						Kind:             InDatabaseBulkOperationGrantKind,
						Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name")),
					},
				},
			},
			Expected: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InDatabase|"database-name"`,
		},
		{
			Name: "grant database role on schema object on all tables in schema",
			Identifier: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewDatabaseObjectIdentifier("database-name", "role-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: sdk.PluralObjectTypeTables,
						Kind:             InSchemaBulkOperationGrantKind,
						Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
					},
				},
			},
			Expected: `"database-name"."role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InSchema|"database-name"."schema-name"`,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Identifier.String())
		})
	}
}
