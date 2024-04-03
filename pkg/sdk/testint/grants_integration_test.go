package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_GrantAndRevokePrivilegesToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertPrivilegeGrantedOnPipe := func(pipeId sdk.SchemaObjectIdentifier, privilege sdk.SchemaObjectPrivilege) {
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypePipe,
					Name:       pipeId,
				},
			},
		})
		require.NoError(t, err)
		require.Contains(t, grantsToPrivileges(grants), privilege.String())
	}

	assertPrivilegeNotGrantedOnPipe := func(pipeId sdk.SchemaObjectIdentifier, privilege sdk.SchemaObjectPrivilege) {
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypePipe,
					Name:       pipeId,
				},
			},
		})
		require.NoError(t, err)
		require.NotContains(t, grantsToPrivileges(grants), privilege.String())
	}

	t.Run("on account", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		privileges := &sdk.AccountRoleGrantPrivileges{
			GlobalPrivileges: []sdk.GlobalPrivilege{sdk.GlobalPrivilegeMonitorUsage, sdk.GlobalPrivilegeApplyTag},
		}
		on := &sdk.AccountRoleGrantOn{
			Account: sdk.Bool(true),
		}
		opts := &sdk.GrantPrivilegesToAccountRoleOptions{
			WithGrantOption: sdk.Bool(true),
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), opts)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 2, len(grants))
		// The order of the grants is not guaranteed
		for _, grant := range grants {
			switch grant.Privilege {
			case sdk.GlobalPrivilegeMonitorUsage.String():
				assert.True(t, grant.GrantOption)
			case sdk.GlobalPrivilegeApplyTag.String():
				assert.True(t, grant.GrantOption)
			default:
				t.Errorf("unexpected privilege: %s", grant.Privilege)
			}
		}

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})

	t.Run("on account object", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		resourceMonitorTest, resourceMonitorCleanup := createResourceMonitor(t, client)
		t.Cleanup(resourceMonitorCleanup)
		privileges := &sdk.AccountRoleGrantPrivileges{
			AccountObjectPrivileges: []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeMonitor},
		}
		on := &sdk.AccountRoleGrantOn{
			AccountObject: &sdk.GrantOnAccountObject{
				ResourceMonitor: sdk.Pointer(resourceMonitorTest.ID()),
			},
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, sdk.AccountObjectPrivilegeMonitor.String(), grants[0].Privilege)

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})

	t.Run("on schema", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		privileges := &sdk.AccountRoleGrantPrivileges{
			SchemaPrivileges: []sdk.SchemaPrivilege{sdk.SchemaPrivilegeCreateAlert},
		}
		on := &sdk.AccountRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: sdk.Pointer(testSchema(t).ID()),
			},
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, sdk.SchemaPrivilegeCreateAlert.String(), grants[0].Privilege)

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})

	t.Run("on schema object", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		tableTest, tableTestCleanup := createTable(t, client, testDb(t), testSchema(t))
		t.Cleanup(tableTestCleanup)
		privileges := &sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect},
		}
		on := &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InSchema:         sdk.Pointer(testSchema(t).ID()),
				},
			},
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, sdk.SchemaObjectPrivilegeSelect.String(), grants[0].Privilege)
		assert.Equal(t, tableTest.ID().FullyQualifiedName(), grants[0].Name.FullyQualifiedName())

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})

	t.Run("on future schema object", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		privileges := &sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect},
		}
		on := &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				Future: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeExternalTables,
					InDatabase:       sdk.Pointer(testDb(t).ID()),
				},
			},
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Future: sdk.Bool(true),
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, sdk.SchemaObjectPrivilegeSelect.String(), grants[0].Privilege)

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})

	t.Run("grant and revoke on all pipes", func(t *testing.T) {
		schema, schemaCleanup := createSchemaWithIdentifier(t, itc.client, testDb(t), random.AlphaN(20))
		t.Cleanup(schemaCleanup)

		table, tableCleanup := createTable(t, itc.client, testDb(t), schema)
		t.Cleanup(tableCleanup)

		stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, schema.Name, random.AlphaN(20)))
		t.Cleanup(stageCleanup)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), schema, random.AlphaN(20), createPipeCopyStatement(t, table, stage))
		t.Cleanup(pipeCleanup)

		secondPipe, secondPipeCleanup := createPipe(t, client, testDb(t), schema, random.AlphaN(20), createPipeCopyStatement(t, table, stage))
		t.Cleanup(secondPipeCleanup)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		err := client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor},
			},
			&sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: sdk.PluralObjectTypePipes,
						InSchema:         sdk.Pointer(schema.ID()),
					},
				},
			},
			role.ID(),
			&sdk.GrantPrivilegesToAccountRoleOptions{},
		)
		require.NoError(t, err)
		assertPrivilegeGrantedOnPipe(pipe.ID(), sdk.SchemaObjectPrivilegeMonitor)
		assertPrivilegeGrantedOnPipe(secondPipe.ID(), sdk.SchemaObjectPrivilegeMonitor)

		err = client.Grants.RevokePrivilegesFromAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor},
			},
			&sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: sdk.PluralObjectTypePipes,
						InSchema:         sdk.Pointer(schema.ID()),
					},
				},
			},
			role.ID(),
			&sdk.RevokePrivilegesFromAccountRoleOptions{},
		)
		require.NoError(t, err)
		assertPrivilegeNotGrantedOnPipe(pipe.ID(), sdk.SchemaObjectPrivilegeMonitor)
		assertPrivilegeNotGrantedOnPipe(secondPipe.ID(), sdk.SchemaObjectPrivilegeMonitor)
	})

	t.Run("grant and revoke on all pipes with multiple errors", func(t *testing.T) {
		schema, schemaCleanup := createSchemaWithIdentifier(t, itc.client, testDb(t), random.AlphaN(20))
		t.Cleanup(schemaCleanup)

		table, tableCleanup := createTable(t, itc.client, testDb(t), schema)
		t.Cleanup(tableCleanup)

		stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, schema.Name, random.AlphaN(20)))
		t.Cleanup(stageCleanup)

		_, pipeCleanup := createPipe(t, client, testDb(t), schema, random.AlphaN(20), createPipeCopyStatement(t, table, stage))
		t.Cleanup(pipeCleanup)

		_, secondPipeCleanup := createPipe(t, client, testDb(t), schema, random.AlphaN(20), createPipeCopyStatement(t, table, stage))
		t.Cleanup(secondPipeCleanup)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		err := client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeInsert}, // Invalid privilege
			},
			&sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: sdk.PluralObjectTypePipes,
						InSchema:         sdk.Pointer(schema.ID()),
					},
				},
			},
			role.ID(),
			&sdk.GrantPrivilegesToAccountRoleOptions{},
		)
		require.Error(t, err)
		joined, ok := err.(interface{ Unwrap() []error }) //nolint:all
		require.True(t, ok)
		require.Len(t, joined.Unwrap(), 2)

		err = client.Grants.RevokePrivilegesFromAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeInsert}, // Invalid privilege
			},
			&sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: sdk.PluralObjectTypePipes,
						InSchema:         sdk.Pointer(schema.ID()),
					},
				},
			},
			role.ID(),
			&sdk.RevokePrivilegesFromAccountRoleOptions{},
		)
		require.Error(t, err)
		joined, ok = err.(interface{ Unwrap() []error }) //nolint:all
		require.True(t, ok)
		require.Len(t, joined.Unwrap(), 2)
	})
}

func TestInt_GrantAndRevokePrivilegesToDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertPrivilegeGrantedOnPipe := func(pipeId sdk.SchemaObjectIdentifier, privilege sdk.SchemaObjectPrivilege) {
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypePipe,
					Name:       pipeId,
				},
			},
		})
		require.NoError(t, err)
		require.Contains(t, grantsToPrivileges(grants), privilege.String())
	}

	assertPrivilegeNotGrantedOnPipe := func(pipeId sdk.SchemaObjectIdentifier, privilege sdk.SchemaObjectPrivilege) {
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypePipe,
					Name:       pipeId,
				},
			},
		})
		require.NoError(t, err)
		require.NotContains(t, grantsToPrivileges(grants), privilege.String())
	}

	t.Run("on database", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, testDb(t))
		databaseRoleId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, databaseRole.Name)

		privileges := &sdk.DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeCreateSchema},
		}
		on := &sdk.DatabaseRoleGrantOn{
			Database: sdk.Pointer(testDb(t).ID()),
		}

		err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Equal(t, 2, len(returnedGrants))

		usagePrivilege, err := collections.FindOne[sdk.Grant](returnedGrants, func(g sdk.Grant) bool { return g.Privilege == sdk.AccountObjectPrivilegeUsage.String() })
		require.NoError(t, err)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, usagePrivilege.GrantedTo)

		createSchemaPrivilege, err := collections.FindOne[sdk.Grant](returnedGrants, func(g sdk.Grant) bool { return g.Privilege == sdk.AccountObjectPrivilegeCreateSchema.String() })
		require.NoError(t, err)
		assert.Equal(t, sdk.ObjectTypeDatabase, createSchemaPrivilege.GrantedOn)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, createSchemaPrivilege.GrantedTo)

		// now revoke and verify that the new grant is gone
		err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedGrants))
		assert.Equal(t, sdk.AccountObjectPrivilegeUsage.String(), returnedGrants[0].Privilege)
	})

	t.Run("on schema", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, testDb(t))
		databaseRoleId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, databaseRole.Name)

		privileges := &sdk.DatabaseRoleGrantPrivileges{
			SchemaPrivileges: []sdk.SchemaPrivilege{sdk.SchemaPrivilegeCreateAlert},
		}
		on := &sdk.DatabaseRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: sdk.Pointer(testSchema(t).ID()),
			},
		}

		err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Equal(t, 2, len(returnedGrants))

		usagePrivilege, err := collections.FindOne[sdk.Grant](returnedGrants, func(g sdk.Grant) bool { return g.Privilege == sdk.AccountObjectPrivilegeUsage.String() })
		require.NoError(t, err)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, usagePrivilege.GrantedTo)

		createAlertPrivilege, err := collections.FindOne[sdk.Grant](returnedGrants, func(g sdk.Grant) bool { return g.Privilege == sdk.SchemaPrivilegeCreateAlert.String() })
		require.NoError(t, err)
		assert.Equal(t, sdk.ObjectTypeSchema, createAlertPrivilege.GrantedOn)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, createAlertPrivilege.GrantedTo)

		// now revoke and verify that the new grant is gone
		err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedGrants))
		assert.Equal(t, sdk.AccountObjectPrivilegeUsage.String(), returnedGrants[0].Privilege)
	})

	t.Run("on schema object", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, testDb(t))
		databaseRoleId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, databaseRole.Name)
		table, _ := createTable(t, client, testDb(t), testSchema(t))

		privileges := &sdk.DatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect},
		}
		on := &sdk.DatabaseRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InSchema:         sdk.Pointer(testSchema(t).ID()),
				},
			},
		}

		err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Equal(t, 2, len(returnedGrants))

		usagePrivilege, err := collections.FindOne[sdk.Grant](returnedGrants, func(g sdk.Grant) bool { return g.Privilege == sdk.AccountObjectPrivilegeUsage.String() })
		require.NoError(t, err)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, usagePrivilege.GrantedTo)

		selectPrivilege, err := collections.FindOne[sdk.Grant](returnedGrants, func(g sdk.Grant) bool { return g.Privilege == sdk.SchemaObjectPrivilegeSelect.String() })
		require.NoError(t, err)
		assert.Equal(t, sdk.ObjectTypeTable, selectPrivilege.GrantedOn)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, selectPrivilege.GrantedTo)
		assert.Equal(t, table.ID().FullyQualifiedName(), selectPrivilege.Name.FullyQualifiedName())

		// now revoke and verify that the new grant is gone
		err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedGrants))
		assert.Equal(t, sdk.AccountObjectPrivilegeUsage.String(), returnedGrants[0].Privilege)
	})

	t.Run("on future schema object", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, testDb(t))
		databaseRoleId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, databaseRole.Name)

		privileges := &sdk.DatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect},
		}
		on := &sdk.DatabaseRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				Future: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeExternalTables,
					InDatabase:       sdk.Pointer(testDb(t).ID()),
				},
			},
		}
		err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Future: sdk.Bool(true),
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(returnedGrants))

		assert.Equal(t, sdk.SchemaObjectPrivilegeSelect.String(), returnedGrants[0].Privilege)
		assert.Equal(t, sdk.ObjectTypeExternalTable, returnedGrants[0].GrantOn)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, returnedGrants[0].GrantTo)

		// now revoke and verify that the new grant is gone
		err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Future: sdk.Bool(true),
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(returnedGrants))
	})

	t.Run("grant and revoke on all pipes", func(t *testing.T) {
		schema, schemaCleanup := createSchemaWithIdentifier(t, itc.client, testDb(t), random.AlphaN(20))
		t.Cleanup(schemaCleanup)

		table, tableCleanup := createTable(t, itc.client, testDb(t), schema)
		t.Cleanup(tableCleanup)

		stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, schema.Name, random.AlphaN(20)))
		t.Cleanup(stageCleanup)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), schema, random.StringN(20), createPipeCopyStatement(t, table, stage))
		t.Cleanup(pipeCleanup)

		secondPipe, secondPipeCleanup := createPipe(t, client, testDb(t), schema, random.StringN(20), createPipeCopyStatement(t, table, stage))
		t.Cleanup(secondPipeCleanup)

		role, roleCleanup := createDatabaseRole(t, client, testDb(t))
		t.Cleanup(roleCleanup)

		err := client.Grants.GrantPrivilegesToDatabaseRole(
			ctx,
			&sdk.DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor},
			},
			&sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: sdk.PluralObjectTypePipes,
						InSchema:         sdk.Pointer(schema.ID()),
					},
				},
			},
			sdk.NewDatabaseObjectIdentifier(testDb(t).Name, role.Name),
			&sdk.GrantPrivilegesToDatabaseRoleOptions{},
		)
		require.NoError(t, err)
		assertPrivilegeGrantedOnPipe(pipe.ID(), sdk.SchemaObjectPrivilegeMonitor)
		assertPrivilegeGrantedOnPipe(secondPipe.ID(), sdk.SchemaObjectPrivilegeMonitor)

		err = client.Grants.RevokePrivilegesFromDatabaseRole(
			ctx,
			&sdk.DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor},
			},
			&sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: sdk.PluralObjectTypePipes,
						InSchema:         sdk.Pointer(schema.ID()),
					},
				},
			},
			sdk.NewDatabaseObjectIdentifier(testDb(t).Name, role.Name),
			&sdk.RevokePrivilegesFromDatabaseRoleOptions{},
		)
		require.NoError(t, err)
		assertPrivilegeNotGrantedOnPipe(pipe.ID(), sdk.SchemaObjectPrivilegeMonitor)
		assertPrivilegeNotGrantedOnPipe(secondPipe.ID(), sdk.SchemaObjectPrivilegeMonitor)
	})

	t.Run("grant and revoke on all pipes with multiple errors", func(t *testing.T) {
		schema, schemaCleanup := createSchemaWithIdentifier(t, itc.client, testDb(t), random.AlphaN(20))
		t.Cleanup(schemaCleanup)

		table, tableCleanup := createTable(t, itc.client, testDb(t), schema)
		t.Cleanup(tableCleanup)

		stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, schema.Name, random.AlphaN(20)))
		t.Cleanup(stageCleanup)

		_, pipeCleanup := createPipe(t, client, testDb(t), schema, random.AlphaN(20), createPipeCopyStatement(t, table, stage))
		t.Cleanup(pipeCleanup)

		_, secondPipeCleanup := createPipe(t, client, testDb(t), schema, random.AlphaN(20), createPipeCopyStatement(t, table, stage))
		t.Cleanup(secondPipeCleanup)

		role, roleCleanup := createDatabaseRole(t, client, testDb(t))
		t.Cleanup(roleCleanup)

		err := client.Grants.GrantPrivilegesToDatabaseRole(
			ctx,
			&sdk.DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeInsert}, // Invalid privilege
			},
			&sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: sdk.PluralObjectTypePipes,
						InSchema:         sdk.Pointer(schema.ID()),
					},
				},
			},
			sdk.NewDatabaseObjectIdentifier(testDb(t).Name, role.Name),
			&sdk.GrantPrivilegesToDatabaseRoleOptions{},
		)
		require.Error(t, err)
		joined, ok := err.(interface{ Unwrap() []error }) //nolint:all
		require.True(t, ok)
		require.Len(t, joined.Unwrap(), 2)

		err = client.Grants.RevokePrivilegesFromDatabaseRole(
			ctx,
			&sdk.DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeInsert}, // Invalid privilege
			},
			&sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: sdk.PluralObjectTypePipes,
						InSchema:         sdk.Pointer(schema.ID()),
					},
				},
			},
			sdk.NewDatabaseObjectIdentifier(testDb(t).Name, role.Name),
			&sdk.RevokePrivilegesFromDatabaseRoleOptions{},
		)
		require.Error(t, err)
		joined, ok = err.(interface{ Unwrap() []error }) //nolint:all
		require.True(t, ok)
		require.Len(t, joined.Unwrap(), 2)
	})
}

func TestInt_GrantPrivilegeToShare(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)

	assertGrant := func(t *testing.T, grants []sdk.Grant, onId sdk.ObjectIdentifier, privilege sdk.ObjectPrivilege) {
		t.Helper()
		var shareGrant *sdk.Grant
		for i, grant := range grants {
			if grant.GranteeName.Name() == shareTest.ID().Name() && grant.Privilege == string(privilege) {
				shareGrant = &grants[i]
				break
			}
		}
		assert.NotNil(t, shareGrant)
		assert.Equal(t, sdk.ObjectTypeTable, shareGrant.GrantedOn)
		assert.Equal(t, sdk.ObjectTypeShare, shareGrant.GrantedTo)
		assert.Equal(t, onId.FullyQualifiedName(), shareGrant.Name.FullyQualifiedName())
	}

	t.Run("with options", func(t *testing.T) {
		err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
				Database: testDb(t).ID(),
			}, shareTest.ID())
			assert.NoError(t, err)
		})

		table, tableCleanup := createTable(t, client, testDb(t), testSchema(t))
		t.Cleanup(tableCleanup)

		err = client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, &sdk.ShareGrantOn{
			Table: &sdk.OnTable{
				AllInSchema: testSchema(t).ID(),
			},
		}, shareTest.ID())
		require.NoError(t, err)

		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypeTable,
					Name:       table.ID(),
				},
			},
		})
		require.NoError(t, err)
		assertGrant(t, grants, table.ID(), sdk.ObjectPrivilegeSelect)

		_, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Share: &sdk.ShowGrantsToShare{
					Name: shareTest.ID(),
				},
			},
		})
		require.NoError(t, err)

		appPackageName := random.AlphaN(8)
		cleanupAppPackage := createApplicationPackage(t, client, appPackageName)
		t.Cleanup(cleanupAppPackage)
		// TODO [SNOW-1284382]: alter the test when the syntax starts working
		// 2024/03/29 17:04:20 [DEBUG] sql-conn-query: [query SHOW GRANTS TO SHARE "0a8DMkl3NOx7" IN APPLICATION PACKAGE "hziiAtqY" err 001003 (42000): SQL compilation error:
		// syntax error line 1 at position 39 unexpected 'APPLICATION'. duration 445.248042ms args {}] (IYA62698)
		_, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Share: &sdk.ShowGrantsToShare{
					Name:                 shareTest.ID(),
					InApplicationPackage: sdk.Pointer(sdk.NewAccountObjectIdentifier(appPackageName)),
				},
			},
		})
		require.Error(t, err)

		err = client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, &sdk.ShareGrantOn{
			Table: &sdk.OnTable{
				AllInSchema: testSchema(t).ID(),
			},
		}, shareTest.ID())
		require.NoError(t, err)
	})
}

func TestInt_RevokePrivilegeToShare(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
		Database: testDb(t).ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Run("without options", func(t *testing.T) {
		err = client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, nil, shareTest.ID())
		require.Error(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		err = client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
}

func TestInt_GrantOwnership(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := createTable(t, itc.client, testDb(t), testSchema(t))
	t.Cleanup(tableCleanup)

	stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphaN(20)))
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)

	checkOwnershipOnObjectToRole := func(t *testing.T, on sdk.OwnershipGrantOn, role string) {
		t.Helper()
		if on.Object == nil {
			t.Error("only on.Object check is supported")
		}
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Object: on.Object,
			},
		})
		require.NoError(t, err)
		_, err = collections.FindOne(grants, func(grant sdk.Grant) bool {
			return grant.Privilege == "OWNERSHIP" && grant.GranteeName.Name() == role
		})
		require.NoError(t, err)
	}

	grantOwnershipToRole := func(t *testing.T, roleName string, on sdk.OwnershipGrantOn) {
		t.Helper()

		err := client.Grants.GrantOwnership(
			ctx,
			on,
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(sdk.NewAccountObjectIdentifier(roleName)),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)
	}

	grantDatabaseAndSchemaUsage := func(t *testing.T, role *sdk.Role) {
		t.Helper()

		err := client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				AccountObjectPrivileges: []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage},
			},
			&sdk.AccountRoleGrantOn{
				AccountObject: &sdk.GrantOnAccountObject{
					Database: sdk.Pointer(sdk.NewAccountObjectIdentifier(TestDatabaseName)),
				},
			},
			role.ID(),
			new(sdk.GrantPrivilegesToAccountRoleOptions),
		)
		require.NoError(t, err)

		err = client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				SchemaPrivileges: []sdk.SchemaPrivilege{sdk.SchemaPrivilegeUsage, sdk.SchemaPrivilegeCreatePipe, sdk.SchemaPrivilegeCreateTask},
			},
			&sdk.AccountRoleGrantOn{
				Schema: &sdk.GrantOnSchema{
					Schema: sdk.Pointer(sdk.NewDatabaseObjectIdentifier(TestDatabaseName, TestSchemaName)),
				},
			},
			role.ID(),
			new(sdk.GrantPrivilegesToAccountRoleOptions),
		)
		require.NoError(t, err)
	}

	grantPipeRole := func(t *testing.T, role *sdk.Role, table *sdk.Table, stage *sdk.Stage) {
		t.Helper()

		grantDatabaseAndSchemaUsage(t, role)

		err := client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect, sdk.SchemaObjectPrivilegeInsert},
			},
			&sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					SchemaObject: &sdk.Object{
						ObjectType: sdk.ObjectTypeTable,
						Name:       table.ID(),
					},
				},
			},
			role.ID(),
			new(sdk.GrantPrivilegesToAccountRoleOptions),
		)
		require.NoError(t, err)

		err = client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeRead},
			},
			&sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					SchemaObject: &sdk.Object{
						ObjectType: sdk.ObjectTypeStage,
						Name:       stage.ID(),
					},
				},
			},
			role.ID(),
			new(sdk.GrantPrivilegesToAccountRoleOptions),
		)
		require.NoError(t, err)
	}

	grantTaskRole := func(t *testing.T, role *sdk.Role) {
		t.Helper()

		grantDatabaseAndSchemaUsage(t, role)

		err := client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				AccountObjectPrivileges: []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage},
			},
			&sdk.AccountRoleGrantOn{
				AccountObject: &sdk.GrantOnAccountObject{
					Warehouse: sdk.Pointer(testWarehouse(t).ID()),
				},
			},
			role.ID(),
			new(sdk.GrantPrivilegesToAccountRoleOptions),
		)
		require.NoError(t, err)
	}

	makeAccountRoleOperableOnPipe := func(t *testing.T, grantingRole string, pipe *sdk.Pipe) {
		t.Helper()

		err := client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			&sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeOperate, sdk.SchemaObjectPrivilegeMonitor},
			},
			&sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					SchemaObject: &sdk.Object{
						ObjectType: sdk.ObjectTypePipe,
						Name:       pipe.ID(),
					},
				},
			},
			sdk.NewAccountObjectIdentifier(grantingRole),
			new(sdk.GrantPrivilegesToAccountRoleOptions),
		)
		require.NoError(t, err)
	}

	ownershipGrantOnObject := func(objectType sdk.ObjectType, identifier sdk.ObjectIdentifier) sdk.OwnershipGrantOn {
		return sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: objectType,
				Name:       identifier,
			},
		}
	}

	ownershipGrantOnPipe := func(pipe *sdk.Pipe) sdk.OwnershipGrantOn {
		return ownershipGrantOnObject(sdk.ObjectTypePipe, pipe.ID())
	}

	ownershipGrantOnTask := func(task *sdk.Task) sdk.OwnershipGrantOn {
		return ownershipGrantOnObject(sdk.ObjectTypeTask, task.ID())
	}

	t.Run("on schema object to database role", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, testDb(t))
		databaseRoleId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, databaseRole.Name)
		table, _ := createTable(t, client, testDb(t), testSchema(t))

		on := sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypeTable,
				Name:       table.ID(),
			},
		}
		to := sdk.OwnershipGrantTo{
			DatabaseRoleName: &databaseRoleId,
		}

		err := client.Grants.GrantOwnership(ctx, on, to, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Equal(t, 2, len(returnedGrants))

		usagePrivilege, err := collections.FindOne[sdk.Grant](returnedGrants, func(g sdk.Grant) bool { return g.Privilege == sdk.AccountObjectPrivilegeUsage.String() })
		require.NoError(t, err)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, usagePrivilege.GrantedTo)

		ownership, err := collections.FindOne[sdk.Grant](returnedGrants, func(g sdk.Grant) bool { return g.Privilege == sdk.SchemaObjectOwnership.String() })
		require.NoError(t, err)
		assert.Equal(t, sdk.ObjectTypeTable, ownership.GrantedOn)
		assert.Equal(t, sdk.ObjectTypeDatabaseRole, ownership.GrantedTo)
		assert.Equal(t, table.ID().FullyQualifiedName(), ownership.Name.FullyQualifiedName())
	})

	t.Run("on future schema object in database to role", func(t *testing.T) {
		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		roleId := role.ID()

		on := sdk.OwnershipGrantOn{
			Future: &sdk.GrantOnSchemaObjectIn{
				PluralObjectType: sdk.PluralObjectTypeExternalTables,
				InDatabase:       sdk.Pointer(testDb(t).ID()),
			},
		}
		to := sdk.OwnershipGrantTo{
			AccountRoleName: &roleId,
		}

		err := client.Grants.GrantOwnership(ctx, on, to, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Future: sdk.Bool(true),
			To: &sdk.ShowGrantsTo{
				Role: roleId,
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(returnedGrants))

		assert.Equal(t, sdk.SchemaObjectOwnership.String(), returnedGrants[0].Privilege)
		assert.Equal(t, sdk.ObjectTypeExternalTable, returnedGrants[0].GrantOn)
		assert.Equal(t, sdk.ObjectTypeRole, returnedGrants[0].GrantTo)
		assert.Equal(t, roleId, returnedGrants[0].GranteeName)
	})

	t.Run("on account level object to role", func(t *testing.T) {
		// role is deliberately created after warehouse, so that cleanup is done in reverse
		// because after ownership grant we lose privilege to drop object
		// with first dropping the role, we reacquire rights to do it - a little hacky trick
		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		roleId := role.ID()

		on := sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypeWarehouse,
				Name:       testWarehouse(t).ID(),
			},
		}
		to := sdk.OwnershipGrantTo{
			AccountRoleName: &roleId,
		}
		currentGrants := sdk.OwnershipCurrentGrants{
			OutboundPrivileges: sdk.Copy,
		}

		err := client.Grants.GrantOwnership(ctx, on, to, &sdk.GrantOwnershipOptions{CurrentGrants: &currentGrants})
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleId,
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(returnedGrants))

		assert.Equal(t, sdk.SchemaObjectOwnership.String(), returnedGrants[0].Privilege)
		assert.Equal(t, sdk.ObjectTypeWarehouse, returnedGrants[0].GrantedOn)
		assert.Equal(t, sdk.ObjectTypeRole, returnedGrants[0].GrantedTo)
		assert.Equal(t, roleId, returnedGrants[0].GranteeName)
	})

	t.Run("on pipe - with ownership", func(t *testing.T) {
		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(pipeCleanup)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.ID().Name())

		currentRole, err := client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)

		grantOwnershipToRole(t, currentRole, ownershipGrantOnPipe(pipe))
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), currentRole)

		pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)
	})

	t.Run("on pipe - with operate and monitor privileges granted", func(t *testing.T) {
		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		pipeRole, pipeRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(pipeRoleCleanup)

		// Role needs usage on the database and schema to later be able to remove pipe in the cleanup
		grantDatabaseAndSchemaUsage(t, role)
		// grantPipeRole grants the necessary privileges to a role to be able to create pipe
		grantPipeRole(t, pipeRole, table, stage)

		previousRole, err := client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)

		// Use a previously prepared role to create a pipe and grant MONITOR + OPERATE to the previously used role (ACCOUNTADMIN).
		usePreviousRole := useRole(t, client, pipeRole.Name)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(func() {
			usePreviousRole = useRole(t, client, role.Name)
			defer usePreviousRole()
			pipeCleanup()
		})

		// Grant MONITOR and OPERATE privileges to the role.
		makeAccountRoleOperableOnPipe(t, previousRole, pipe)

		usePreviousRole()

		err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				PipeExecutionPaused: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			&sdk.GrantOwnershipOptions{
				CurrentGrants: &sdk.OwnershipCurrentGrants{
					OutboundPrivileges: sdk.Revoke, // To revoke MONITOR privilege from ACCOUNTADMIN automatically
				},
			},
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.ID().Name())

		usePreviousRole()

		pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)
	})

	t.Run("on pipe - with operate privilege granted and copy current grants option", func(t *testing.T) {
		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		pipeRole, pipeRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(pipeRoleCleanup)

		// Role needs usage on the database and schema to later be able to remove pipe in the cleanup
		grantDatabaseAndSchemaUsage(t, role)
		// grantPipeRole grants the necessary privileges to a role to be able to create pipe
		grantPipeRole(t, pipeRole, table, stage)

		previousRole, err := client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)

		// Use a previously prepared role to create a pipe and grant MONITOR + OPERATE to the previously used role (ACCOUNTADMIN).
		usePreviousRole := useRole(t, client, pipeRole.Name)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(func() {
			usePreviousRole = useRole(t, client, role.Name)
			defer usePreviousRole()
			pipeCleanup()
		})

		// Grant MONITOR and OPERATE privileges to the role.
		makeAccountRoleOperableOnPipe(t, previousRole, pipe)

		usePreviousRole()

		err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				PipeExecutionPaused: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			&sdk.GrantOwnershipOptions{
				CurrentGrants: &sdk.OwnershipCurrentGrants{
					OutboundPrivileges: sdk.Copy, // With copy, we'll be able to resume the pipe after ownership transfer
				},
			},
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.ID().Name())

		usePreviousRole()

		pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)
	})

	t.Run("on pipe - with neither ownership nor operate", func(t *testing.T) {
		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		pipeRole, pipeRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(pipeRoleCleanup)

		// Role needs usage on the database and schema to later be able to remove pipe in the cleanup
		grantDatabaseAndSchemaUsage(t, role)
		// grantPipeRole grants the necessary privileges to a role to be able to create pipe
		grantPipeRole(t, pipeRole, table, stage)

		// Use a previously prepared role to create a pipe and grant MONITOR + OPERATE to the previously used role (ACCOUNTADMIN).
		usePreviousRole := useRole(t, client, pipeRole.Name)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(func() {
			usePreviousRole = useRole(t, client, pipeRole.Name)
			defer usePreviousRole()
			pipeCleanup()
		})

		err := client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				PipeExecutionPaused: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		usePreviousRole()

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.ErrorContains(t, err, fmt.Sprintf("Pipe %s not in paused state. To pause pipe run ALTER PIPE %s SET PIPE_EXECUTION_PAUSED=true", pipe.Name, pipe.Name))
	})

	t.Run("on pipe - with neither ownership nor operate on paused pipe", func(t *testing.T) {
		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		pipeRole, pipeRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(pipeRoleCleanup)

		// Role needs usage on the database and schema to later be able to remove pipe in the cleanup
		grantDatabaseAndSchemaUsage(t, role)
		// grantPipeRole grants the necessary privileges to a role to be able to create pipe
		grantPipeRole(t, pipeRole, table, stage)

		// Use a previously prepared role to create a pipe and grant MONITOR + OPERATE to the previously used role (ACCOUNTADMIN).
		usePreviousRole := useRole(t, client, pipeRole.Name)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(func() {
			usePreviousRole = useRole(t, client, role.Name)
			defer usePreviousRole()
			pipeCleanup()
		})

		err := client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				PipeExecutionPaused: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)

		usePreviousRole()

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.Name)
	})

	t.Run("on all pipes", func(t *testing.T) {
		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(pipeCleanup)

		secondPipe, secondPipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(secondPipeCleanup)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		secondPipeExecutionState, err := client.SystemFunctions.PipeStatus(secondPipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, secondPipeExecutionState)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		onAllPipesInSchema := sdk.OwnershipGrantOn{
			All: &sdk.GrantOnSchemaObjectIn{
				PluralObjectType: sdk.PluralObjectTypePipes,
				InSchema:         sdk.Pointer(testSchema(t).ID()),
			},
		}
		err = client.Grants.GrantOwnership(
			ctx,
			onAllPipesInSchema,
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)

		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.ID().Name())
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(secondPipe), role.ID().Name())

		currentRole, err := client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)
		grantOwnershipToRole(t, currentRole, onAllPipesInSchema)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), currentRole)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(secondPipe), currentRole)

		pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)

		secondPipeExecutionState, err = client.SystemFunctions.PipeStatus(secondPipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, secondPipeExecutionState)
	})

	t.Run("on pipe - with ownership", func(t *testing.T) {
		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(pipeCleanup)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.ID().Name())

		currentRole, err := client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)

		grantOwnershipToRole(t, currentRole, ownershipGrantOnPipe(pipe))
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), currentRole)

		pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)
	})

	t.Run("on pipe - with operate and monitor privileges granted", func(t *testing.T) {
		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		pipeRole, pipeRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(pipeRoleCleanup)

		// Role needs usage on the database and schema to later be able to remove pipe in the cleanup
		grantDatabaseAndSchemaUsage(t, role)
		// grantPipeRole grants the necessary privileges to a role to be able to create pipe
		grantPipeRole(t, pipeRole, table, stage)

		previousRole, err := client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)

		// Use a previously prepared role to create a pipe and grant MONITOR + OPERATE to the previously used role (ACCOUNTADMIN).
		usePreviousRole := useRole(t, client, pipeRole.Name)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(func() {
			usePreviousRole = useRole(t, client, role.Name)
			pipeCleanup()
			usePreviousRole()
		})

		// Grant MONITOR and OPERATE privileges to the role.
		makeAccountRoleOperableOnPipe(t, previousRole, pipe)

		usePreviousRole()

		err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				PipeExecutionPaused: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			&sdk.GrantOwnershipOptions{
				CurrentGrants: &sdk.OwnershipCurrentGrants{
					OutboundPrivileges: sdk.Revoke, // To revoke MONITOR privilege from ACCOUNTADMIN automatically
				},
			},
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.ID().Name())

		usePreviousRole()

		pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)
	})

	t.Run("on pipe - with operate privilege granted and copy current grants option", func(t *testing.T) {
		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		pipeRole, pipeRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(pipeRoleCleanup)

		// Role needs usage on the database and schema to later be able to remove pipe in the cleanup
		grantDatabaseAndSchemaUsage(t, role)
		// grantPipeRole grants the necessary privileges to a role to be able to create pipe
		grantPipeRole(t, pipeRole, table, stage)

		previousRole, err := client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)

		// Use a previously prepared role to create a pipe and grant MONITOR + OPERATE to the previously used role (ACCOUNTADMIN).
		usePreviousRole := useRole(t, client, pipeRole.Name)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(func() {
			usePreviousRole = useRole(t, client, role.Name)
			pipeCleanup()
			usePreviousRole()
		})

		// Grant MONITOR and OPERATE privileges to the role.
		makeAccountRoleOperableOnPipe(t, previousRole, pipe)

		usePreviousRole()

		err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				PipeExecutionPaused: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			&sdk.GrantOwnershipOptions{
				CurrentGrants: &sdk.OwnershipCurrentGrants{
					OutboundPrivileges: sdk.Copy, // With copy, we'll be able to resume the pipe after ownership transfer
				},
			},
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.ID().Name())

		usePreviousRole()

		pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)
	})

	t.Run("on pipe - with neither ownership nor operate", func(t *testing.T) {
		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		pipeRole, pipeRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(pipeRoleCleanup)

		// Role needs usage on the database and schema to later be able to remove pipe in the cleanup
		grantDatabaseAndSchemaUsage(t, role)
		// grantPipeRole grants the necessary privileges to a role to be able to create pipe
		grantPipeRole(t, pipeRole, table, stage)

		// Use a previously prepared role to create a pipe and grant MONITOR + OPERATE to the previously used role (ACCOUNTADMIN).
		usePreviousRole := useRole(t, client, pipeRole.Name)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(func() {
			usePreviousRole = useRole(t, client, pipeRole.Name)
			pipeCleanup()
			usePreviousRole()
		})

		err := client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				PipeExecutionPaused: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		usePreviousRole()

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.ErrorContains(t, err, fmt.Sprintf("Pipe %s not in paused state. To pause pipe run ALTER PIPE %s SET PIPE_EXECUTION_PAUSED=true", pipe.Name, pipe.Name))
	})

	t.Run("on pipe - with neither ownership nor operate on paused pipe", func(t *testing.T) {
		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		pipeRole, pipeRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(pipeRoleCleanup)

		// Role needs usage on the database and schema to later be able to remove pipe in the cleanup
		grantDatabaseAndSchemaUsage(t, role)
		// grantPipeRole grants the necessary privileges to a role to be able to create pipe
		grantPipeRole(t, pipeRole, table, stage)

		// Use a previously prepared role to create a pipe and grant MONITOR + OPERATE to the previously used role (ACCOUNTADMIN).
		usePreviousRole := useRole(t, client, pipeRole.Name)

		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(func() {
			usePreviousRole = useRole(t, client, role.Name)
			pipeCleanup()
			usePreviousRole()
		})

		err := client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				PipeExecutionPaused: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)

		usePreviousRole()

		err = client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnPipe(pipe),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.Name)
	})

	t.Run("on all pipes", func(t *testing.T) {
		pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(pipeCleanup)

		secondPipe, secondPipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
		t.Cleanup(secondPipeCleanup)

		pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

		secondPipeExecutionState, err := client.SystemFunctions.PipeStatus(secondPipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.RunningPipeExecutionState, secondPipeExecutionState)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		onAllPipesInSchema := sdk.OwnershipGrantOn{
			All: &sdk.GrantOnSchemaObjectIn{
				PluralObjectType: sdk.PluralObjectTypePipes,
				InSchema:         sdk.Pointer(testSchema(t).ID()),
			},
		}
		err = client.Grants.GrantOwnership(
			ctx,
			onAllPipesInSchema,
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)

		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), role.ID().Name())
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(secondPipe), role.ID().Name())

		currentRole, err := client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)
		grantOwnershipToRole(t, currentRole, onAllPipesInSchema)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(pipe), currentRole)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnPipe(secondPipe), currentRole)

		pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)

		secondPipeExecutionState, err = client.SystemFunctions.PipeStatus(secondPipe.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.PausedPipeExecutionState, secondPipeExecutionState)
	})

	t.Run("on task - with ownership", func(t *testing.T) {
		task, taskCleanup := createTask(t, client, testDb(t), testSchema(t))
		t.Cleanup(taskCleanup)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		err := client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnTask(task),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnTask(task), role.ID().Name())
	})

	t.Run("on task - without ownership", func(t *testing.T) {
		taskRole, taskRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(taskRoleCleanup)

		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		// Role needs usage on the database and schema to later be able to remove task in the cleanup
		grantDatabaseAndSchemaUsage(t, role)

		// grantTaskRole grants the necessary privileges to a role to be able to create task
		grantTaskRole(t, taskRole)

		// Use a previously prepared role to create a task
		usePreviousRole := useRole(t, client, taskRole.Name)

		taskId := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(20))
		withWarehouseReq := sdk.NewCreateTaskWarehouseRequest().WithWarehouse(sdk.Pointer(testWarehouse(t).ID()))
		task, taskCleanup := createTaskWithRequest(t, client, sdk.NewCreateTaskRequest(taskId, "SELECT CURRENT_TIMESTAMP").WithWarehouse(withWarehouseReq))
		t.Cleanup(taskCleanup)

		usePreviousRole()

		t.Cleanup(func() {
			currentRole, err := client.ContextFunctions.CurrentRole(ctx)
			require.NoError(t, err)

			grantOwnershipToRole(t, currentRole, ownershipGrantOnTask(task))
		})

		err := client.Grants.GrantOwnership(
			ctx,
			ownershipGrantOnTask(task),
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)
		checkOwnershipOnObjectToRole(t, ownershipGrantOnTask(task), role.ID().Name())
	})

	t.Run("on all tasks -  with ownership", func(t *testing.T) {
		task, taskCleanup := createTask(t, client, testDb(t), testSchema(t))
		t.Cleanup(taskCleanup)

		secondTask, secondTaskCleanup := createTask(t, client, testDb(t), testSchema(t))
		t.Cleanup(secondTaskCleanup)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		onAllTasks := sdk.OwnershipGrantOn{
			All: &sdk.GrantOnSchemaObjectIn{
				PluralObjectType: sdk.PluralObjectTypeTasks,
				InSchema:         sdk.Pointer(testSchema(t).ID()),
			},
		}
		err := client.Grants.GrantOwnership(
			ctx,
			onAllTasks,
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)

		checkOwnershipOnObjectToRole(t, ownershipGrantOnTask(task), role.ID().Name())
		checkOwnershipOnObjectToRole(t, ownershipGrantOnTask(secondTask), role.ID().Name())
	})

	t.Run("on all tasks -  without ownership", func(t *testing.T) {
		taskRole, taskRoleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(taskRoleCleanup)

		role, roleCleanup := createRoleGrantedToCurrentUser(t, client)
		t.Cleanup(roleCleanup)

		// Role needs usage on the database and schema to later be able to remove task in the cleanup
		grantDatabaseAndSchemaUsage(t, role)

		// grantTaskRole grants the necessary privileges to a role to be able to create task
		grantTaskRole(t, taskRole)

		// Use a previously prepared role to create a task
		usePreviousRole := useRole(t, client, taskRole.Name)

		taskId := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(20))
		withWarehouseReq := sdk.NewCreateTaskWarehouseRequest().WithWarehouse(sdk.Pointer(testWarehouse(t).ID()))
		task, taskCleanup := createTaskWithRequest(t, client, sdk.NewCreateTaskRequest(taskId, "SELECT CURRENT_TIMESTAMP").WithWarehouse(withWarehouseReq))
		t.Cleanup(taskCleanup)

		secondTaskId := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(20))
		secondTask, secondTaskCleanup := createTaskWithRequest(t, client, sdk.NewCreateTaskRequest(secondTaskId, "SELECT CURRENT_TIMESTAMP").WithWarehouse(withWarehouseReq))
		t.Cleanup(secondTaskCleanup)

		usePreviousRole()

		t.Cleanup(func() {
			currentRole, err := client.ContextFunctions.CurrentRole(ctx)
			require.NoError(t, err)

			grantOwnershipToRole(t, currentRole, ownershipGrantOnTask(task))
			grantOwnershipToRole(t, currentRole, ownershipGrantOnTask(secondTask))
		})

		onAllTasks := sdk.OwnershipGrantOn{
			All: &sdk.GrantOnSchemaObjectIn{
				PluralObjectType: sdk.PluralObjectTypeTasks,
				InSchema:         sdk.Pointer(testSchema(t).ID()),
			},
		}
		err := client.Grants.GrantOwnership(
			ctx,
			onAllTasks,
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(role.ID()),
			},
			new(sdk.GrantOwnershipOptions),
		)
		require.NoError(t, err)

		checkOwnershipOnObjectToRole(t, ownershipGrantOnTask(task), role.ID().Name())
		checkOwnershipOnObjectToRole(t, ownershipGrantOnTask(secondTask), role.ID().Name())
	})
}

func TestInt_ShowGrants(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
		Database: testDb(t).ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err = client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
	t.Run("without options", func(t *testing.T) {
		_, err := client.Grants.Show(ctx, nil)
		require.NoError(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypeDatabase,
					Name:       testDb(t).ID(),
				},
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(grants))
	})
}

func grantsToPrivileges(grants []sdk.Grant) []string {
	privileges := make([]string, len(grants))
	for i, grant := range grants {
		privileges[i] = grant.Privilege
	}
	return privileges
}
