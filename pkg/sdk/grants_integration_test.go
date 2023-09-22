package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_GrantAndRevokePrivilegesToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("on account", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		privileges := &AccountRoleGrantPrivileges{
			GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
		}
		on := &AccountRoleGrantOn{
			Account: Bool(true),
		}
		opts := &GrantPrivilegesToAccountRoleOptions{
			WithGrantOption: Bool(true),
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), opts)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 2, len(grants))
		// The order of the grants is not guaranteed
		for _, grant := range grants {
			switch grant.Privilege {
			case GlobalPrivilegeMonitorUsage.String():
				assert.True(t, grant.GrantOption)
			case GlobalPrivilegeApplyTag.String():
				assert.True(t, grant.GrantOption)
			default:
				t.Errorf("unexpected privilege: %s", grant.Privilege)
			}
		}

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
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
		privileges := &AccountRoleGrantPrivileges{
			AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeMonitor},
		}
		on := &AccountRoleGrantOn{
			AccountObject: &GrantOnAccountObject{
				ResourceMonitor: Pointer(resourceMonitorTest.ID()),
			},
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, AccountObjectPrivilegeMonitor.String(), grants[0].Privilege)

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})

	t.Run("on schema", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
		t.Cleanup(schemaCleanup)
		privileges := &AccountRoleGrantPrivileges{
			SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
		}
		on := &AccountRoleGrantOn{
			Schema: &GrantOnSchema{
				Schema: Pointer(schemaTest.ID()),
			},
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, SchemaPrivilegeCreateAlert.String(), grants[0].Privilege)

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})

	t.Run("on schema object", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
		t.Cleanup(schemaCleanup)
		tableTest, tableTestCleanup := createTable(t, client, databaseTest, schemaTest)
		t.Cleanup(tableTestCleanup)
		privileges := &AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		on := &AccountRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				All: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
					InSchema:         Pointer(schemaTest.ID()),
				},
			},
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, SchemaObjectPrivilegeSelect.String(), grants[0].Privilege)
		assert.Equal(t, tableTest.ID().FullyQualifiedName(), grants[0].Name.FullyQualifiedName())

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})

	t.Run("on future schema object", func(t *testing.T) {
		roleTest, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		privileges := &AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		on := &AccountRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeExternalTables,
					InDatabase:       Pointer(databaseTest.ID()),
				},
			},
		}
		err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			Future: Bool(true),
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, SchemaObjectPrivilegeSelect.String(), grants[0].Privilege)

		// now revoke and verify that the grant(s) are gone
		err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privileges, on, roleTest.ID(), nil)
		require.NoError(t, err)
		grants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(grants))
	})
}

func TestInt_GrantAndRevokePrivilegesToDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	t.Run("on database", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, database)
		databaseRoleId := NewDatabaseObjectIdentifier(database.Name, databaseRole.Name)

		privileges := &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
		}
		on := &DatabaseRoleGrantOn{
			Database: Pointer(database.ID()),
		}

		err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Equal(t, 2, len(returnedGrants))

		usagePrivilege, err := findOne[Grant](returnedGrants, func(g Grant) bool { return g.Privilege == AccountObjectPrivilegeUsage.String() })
		require.NoError(t, err)
		assert.Equal(t, ObjectTypeDatabaseRole, usagePrivilege.GrantedTo)

		createSchemaPrivilege, err := findOne[Grant](returnedGrants, func(g Grant) bool { return g.Privilege == AccountObjectPrivilegeCreateSchema.String() })
		require.NoError(t, err)
		assert.Equal(t, ObjectTypeDatabase, createSchemaPrivilege.GrantedOn)
		assert.Equal(t, ObjectTypeDatabaseRole, createSchemaPrivilege.GrantedTo)

		// now revoke and verify that the new grant is gone
		err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedGrants))
		assert.Equal(t, AccountObjectPrivilegeUsage.String(), returnedGrants[0].Privilege)
	})

	t.Run("on schema", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, database)
		databaseRoleId := NewDatabaseObjectIdentifier(database.Name, databaseRole.Name)
		schema, _ := createSchema(t, client, database)

		privileges := &DatabaseRoleGrantPrivileges{
			SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
		}
		on := &DatabaseRoleGrantOn{
			Schema: &GrantOnSchema{
				Schema: Pointer(schema.ID()),
			},
		}

		err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Equal(t, 2, len(returnedGrants))

		usagePrivilege, err := findOne[Grant](returnedGrants, func(g Grant) bool { return g.Privilege == AccountObjectPrivilegeUsage.String() })
		require.NoError(t, err)
		assert.Equal(t, ObjectTypeDatabaseRole, usagePrivilege.GrantedTo)

		createAlertPrivilege, err := findOne[Grant](returnedGrants, func(g Grant) bool { return g.Privilege == SchemaPrivilegeCreateAlert.String() })
		require.NoError(t, err)
		assert.Equal(t, ObjectTypeSchema, createAlertPrivilege.GrantedOn)
		assert.Equal(t, ObjectTypeDatabaseRole, createAlertPrivilege.GrantedTo)

		// now revoke and verify that the new grant is gone
		err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedGrants))
		assert.Equal(t, AccountObjectPrivilegeUsage.String(), returnedGrants[0].Privilege)
	})

	t.Run("on schema object", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, database)
		databaseRoleId := NewDatabaseObjectIdentifier(database.Name, databaseRole.Name)
		schema, _ := createSchema(t, client, database)
		table, _ := createTable(t, client, database, schema)

		privileges := &DatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		on := &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				All: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
					InSchema:         Pointer(schema.ID()),
				},
			},
		}

		err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Equal(t, 2, len(returnedGrants))

		usagePrivilege, err := findOne[Grant](returnedGrants, func(g Grant) bool { return g.Privilege == AccountObjectPrivilegeUsage.String() })
		require.NoError(t, err)
		assert.Equal(t, ObjectTypeDatabaseRole, usagePrivilege.GrantedTo)

		selectPrivilege, err := findOne[Grant](returnedGrants, func(g Grant) bool { return g.Privilege == SchemaObjectPrivilegeSelect.String() })
		require.NoError(t, err)
		assert.Equal(t, ObjectTypeTable, selectPrivilege.GrantedOn)
		assert.Equal(t, ObjectTypeDatabaseRole, selectPrivilege.GrantedTo)
		assert.Equal(t, table.ID().FullyQualifiedName(), selectPrivilege.Name.FullyQualifiedName())

		// now revoke and verify that the new grant is gone
		err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedGrants))
		assert.Equal(t, AccountObjectPrivilegeUsage.String(), returnedGrants[0].Privilege)
	})

	t.Run("on future schema object", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, database)
		databaseRoleId := NewDatabaseObjectIdentifier(database.Name, databaseRole.Name)

		privileges := &DatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		on := &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeExternalTables,
					InDatabase:       Pointer(database.ID()),
				},
			},
		}
		err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			Future: Bool(true),
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(returnedGrants))

		assert.Equal(t, SchemaObjectPrivilegeSelect.String(), returnedGrants[0].Privilege)
		assert.Equal(t, ObjectTypeExternalTable, returnedGrants[0].GrantOn)
		assert.Equal(t, ObjectTypeDatabaseRole, returnedGrants[0].GrantTo)

		// now revoke and verify that the new grant is gone
		err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		returnedGrants, err = client.Grants.Show(ctx, &ShowGrantOptions{
			Future: Bool(true),
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(returnedGrants))
	})
}

func TestInt_GrantPrivilegeToShare(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	t.Run("without options", func(t *testing.T) {
		err := client.Grants.GrantPrivilegeToShare(ctx, ObjectPrivilegeUsage, nil, shareTest.ID())
		require.Error(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		err := client.Grants.GrantPrivilegeToShare(ctx, ObjectPrivilegeUsage, &GrantPrivilegeToShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       databaseTest.ID(),
				},
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(grants))
		var shareGrant *Grant
		for i, grant := range grants {
			if grant.GranteeName.Name() == shareTest.ID().Name() {
				shareGrant = &grants[i]
				break
			}
		}
		assert.NotNil(t, shareGrant)
		assert.Equal(t, string(ObjectPrivilegeUsage), shareGrant.Privilege)
		assert.Equal(t, ObjectTypeDatabase, shareGrant.GrantedOn)
		assert.Equal(t, ObjectTypeShare, shareGrant.GrantedTo)
		assert.Equal(t, databaseTest.ID().Name(), shareGrant.Name.Name())
		err = client.Grants.RevokePrivilegeFromShare(ctx, ObjectPrivilegeUsage, &RevokePrivilegeFromShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
}

func TestInt_RevokePrivilegeToShare(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, ObjectPrivilegeUsage, &GrantPrivilegeToShareOn{
		Database: databaseTest.ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Run("without options", func(t *testing.T) {
		err = client.Grants.RevokePrivilegeFromShare(ctx, ObjectPrivilegeUsage, nil, shareTest.ID())
		require.Error(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		err = client.Grants.RevokePrivilegeFromShare(ctx, ObjectPrivilegeUsage, &RevokePrivilegeFromShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
}

func TestInt_GrantOwnership(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	t.Run("on schema object to database role", func(t *testing.T) {
		databaseRole, _ := createDatabaseRole(t, client, database)
		databaseRoleId := NewDatabaseObjectIdentifier(database.Name, databaseRole.Name)
		schema, _ := createSchema(t, client, database)
		table, _ := createTable(t, client, database, schema)

		on := OwnershipGrantOn{
			SchemaObject: GrantOnSchemaObject{
				SchemaObject: &Object{
					ObjectType: ObjectTypeTable,
					Name:       table.ID(),
				},
			},
		}
		to := OwnershipGrantTo{
			DatabaseRoleName: &databaseRoleId,
		}

		err := client.Grants.GrantOwnership(ctx, on, to, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Equal(t, 2, len(returnedGrants))

		usagePrivilege, err := findOne[Grant](returnedGrants, func(g Grant) bool { return g.Privilege == AccountObjectPrivilegeUsage.String() })
		require.NoError(t, err)
		assert.Equal(t, ObjectTypeDatabaseRole, usagePrivilege.GrantedTo)

		ownership, err := findOne[Grant](returnedGrants, func(g Grant) bool { return g.Privilege == SchemaObjectOwnership.String() })
		require.NoError(t, err)
		assert.Equal(t, ObjectTypeTable, ownership.GrantedOn)
		assert.Equal(t, ObjectTypeDatabaseRole, ownership.GrantedTo)
		assert.Equal(t, table.ID().FullyQualifiedName(), ownership.Name.FullyQualifiedName())
	})

	t.Run("on future schema object in database to role", func(t *testing.T) {
		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)
		roleId := role.ID()

		on := OwnershipGrantOn{
			SchemaObject: GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeExternalTables,
					InDatabase:       Pointer(database.ID()),
				},
			},
		}
		to := OwnershipGrantTo{
			AccountRoleName: &roleId,
		}

		err := client.Grants.GrantOwnership(ctx, on, to, nil)
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			Future: Bool(true),
			To: &ShowGrantsTo{
				Role: roleId,
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(returnedGrants))

		assert.Equal(t, SchemaObjectOwnership.String(), returnedGrants[0].Privilege)
		assert.Equal(t, ObjectTypeExternalTable, returnedGrants[0].GrantOn)
		assert.Equal(t, ObjectTypeRole, returnedGrants[0].GrantTo)
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

		on := OwnershipGrantOn{
			SchemaObject: GrantOnSchemaObject{
				SchemaObject: &Object{
					ObjectType: ObjectTypeWarehouse,
					Name:       warehouse.ID(),
				},
			},
		}
		to := OwnershipGrantTo{
			AccountRoleName: &roleId,
		}
		currentGrants := OwnershipCurrentGrants{
			OutboundPrivileges: Copy,
		}

		err := client.Grants.GrantOwnership(ctx, on, to, &GrantOwnershipOptions{CurrentGrants: &currentGrants})
		require.NoError(t, err)

		returnedGrants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleId,
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(returnedGrants))

		assert.Equal(t, SchemaObjectOwnership.String(), returnedGrants[0].Privilege)
		assert.Equal(t, ObjectTypeWarehouse, returnedGrants[0].GrantedOn)
		assert.Equal(t, ObjectTypeRole, returnedGrants[0].GrantedTo)
		assert.Equal(t, roleId, returnedGrants[0].GranteeName)
	})
}

func TestInt_ShowGrants(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, ObjectPrivilegeUsage, &GrantPrivilegeToShareOn{
		Database: databaseTest.ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err = client.Grants.RevokePrivilegeFromShare(ctx, ObjectPrivilegeUsage, &RevokePrivilegeFromShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
	t.Run("without options", func(t *testing.T) {
		_, err := client.Grants.Show(ctx, nil)
		require.NoError(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       databaseTest.ID(),
				},
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(grants))
	})
}
