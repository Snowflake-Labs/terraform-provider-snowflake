package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_GrantAndRevokePrivilegesToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

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
		schemaTest, schemaCleanup := createSchema(t, client, testDb(t))
		t.Cleanup(schemaCleanup)
		privileges := &sdk.AccountRoleGrantPrivileges{
			SchemaPrivileges: []sdk.SchemaPrivilege{sdk.SchemaPrivilegeCreateAlert},
		}
		on := &sdk.AccountRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: sdk.Pointer(schemaTest.ID()),
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
		schemaTest, schemaCleanup := createSchema(t, client, testDb(t))
		t.Cleanup(schemaCleanup)
		tableTest, tableTestCleanup := createTable(t, client, testDb(t), schemaTest)
		t.Cleanup(tableTestCleanup)
		privileges := &sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect},
		}
		on := &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InSchema:         sdk.Pointer(schemaTest.ID()),
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
}

func TestInt_GrantAndRevokePrivilegesToDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

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
		schema, _ := createSchema(t, client, testDb(t))

		privileges := &sdk.DatabaseRoleGrantPrivileges{
			SchemaPrivileges: []sdk.SchemaPrivilege{sdk.SchemaPrivilegeCreateAlert},
		}
		on := &sdk.DatabaseRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: sdk.Pointer(schema.ID()),
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
		schema, _ := createSchema(t, client, testDb(t))
		table, _ := createTable(t, client, testDb(t), schema)

		privileges := &sdk.DatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect},
		}
		on := &sdk.DatabaseRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InSchema:         sdk.Pointer(schema.ID()),
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
}

func TestInt_GrantPrivilegeToShare(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	t.Run("without options", func(t *testing.T) {
		err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, nil, shareTest.ID())
		require.Error(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
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
		var shareGrant *sdk.Grant
		for i, grant := range grants {
			if grant.GranteeName.Name() == shareTest.ID().Name() {
				shareGrant = &grants[i]
				break
			}
		}
		assert.NotNil(t, shareGrant)
		assert.Equal(t, string(sdk.ObjectPrivilegeUsage), shareGrant.Privilege)
		assert.Equal(t, sdk.ObjectTypeDatabase, shareGrant.GrantedOn)
		assert.Equal(t, sdk.ObjectTypeShare, shareGrant.GrantedTo)
		assert.Equal(t, testDb(t).ID().Name(), shareGrant.Name.Name())
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
}

func TestInt_RevokePrivilegeToShare(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
		Database: testDb(t).ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Run("without options", func(t *testing.T) {
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, nil, shareTest.ID())
		require.Error(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
}

func TestInt_GrantOwnership(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("on schema object to database role", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, testDb(t))
		databaseRoleId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, databaseRole.Name)
		schema, _ := createSchema(t, client, testDb(t))
		table, _ := createTable(t, client, testDb(t), schema)

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
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		// role is deliberately created after warehouse, so that cleanup is done in reverse
		// because after ownership grant we lose privilege to drop object
		// with first dropping the role, we reacquire rights to do it - a little hacky trick
		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		roleId := role.ID()

		on := sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypeWarehouse,
				Name:       warehouse.ID(),
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
}

func TestInt_ShowGrants(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
		Database: testDb(t).ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
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
