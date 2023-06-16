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
		for _, grant := range grants {
			if grant.GranteeName.Name() == shareTest.ID().Name() {
				shareGrant = grant
				break
			}
		}
		assert.NotNil(t, shareGrant)
		assert.Equal(t, ObjectPrivilegeUsage, shareGrant.Privilege)
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
		require.Error(t, err)
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
