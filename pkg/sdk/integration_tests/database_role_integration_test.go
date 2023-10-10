package sdk_integration_tests

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DatabaseRoles(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := sdk.createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	assertDatabaseRole := func(t *testing.T, databaseRole *sdk.DatabaseRole, expectedName string, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, databaseRole.CreatedOn)
		assert.Equal(t, expectedName, databaseRole.Name)
		assert.Equal(t, "ACCOUNTADMIN", databaseRole.Owner)
		assert.Equal(t, expectedComment, databaseRole.Comment)
		assert.Equal(t, 0, databaseRole.GrantedToRoles)
		assert.Equal(t, 0, databaseRole.GrantedToDatabaseRoles)
		assert.Equal(t, 0, databaseRole.GrantedDatabaseRoles)
	}

	cleanupDatabaseRoleProvider := func(id sdk.DatabaseObjectIdentifier) func() {
		return func() {
			err := client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
			require.NoError(t, err)
		}
	}

	createDatabaseRole := func(t *testing.T) *sdk.DatabaseRole {
		t.Helper()
		name := sdk.randomString(t)
		id := sdk.NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		return databaseRole
	}

	t.Run("create database_role: complete case", func(t *testing.T) {
		name := sdk.randomString(t)
		id := sdk.NewDatabaseObjectIdentifier(database.Name, name)
		comment := sdk.randomComment(t)

		request := sdk.NewCreateDatabaseRoleRequest(id).WithComment(&comment).WithIfNotExists(true)
		err := client.DatabaseRoles.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)

		require.NoError(t, err)
		assertDatabaseRole(t, databaseRole, name, comment)
	})

	t.Run("create database_role: no optionals", func(t *testing.T) {
		name := sdk.randomString(t)
		id := sdk.NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assertDatabaseRole(t, databaseRole, name, "")
	})

	t.Run("drop database_role: existing", func(t *testing.T) {
		name := sdk.randomString(t)
		id := sdk.NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)

		err = client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
		require.NoError(t, err)

		_, err = client.DatabaseRoles.ShowByID(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("drop database_role: non-existing", func(t *testing.T) {
		id := sdk.NewDatabaseObjectIdentifier(database.Name, "does_not_exist")

		err := client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("alter database_role: set value and unset value", func(t *testing.T) {
		name := sdk.randomString(t)
		id := sdk.NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		alterRequest := sdk.NewAlterDatabaseRoleRequest(id).WithSetComment("new comment")
		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredDatabaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredDatabaseRole.Comment)

		alterRequest = sdk.NewAlterDatabaseRoleRequest(id).WithUnsetComment()
		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredDatabaseRole, err = client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredDatabaseRole.Comment)
	})

	t.Run("alter database_role: rename", func(t *testing.T) {
		name := sdk.randomString(t)
		id := sdk.NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)

		newName := sdk.randomString(t)
		newId := sdk.NewDatabaseObjectIdentifier(database.Name, newName)
		alterRequest := sdk.NewAlterDatabaseRoleRequest(id).WithRename(newId)

		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupDatabaseRoleProvider(id))
		} else {
			t.Cleanup(cleanupDatabaseRoleProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.DatabaseRoles.ShowByID(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertDatabaseRole(t, databaseRole, newName, "")
	})

	t.Run("alter database_role: rename to other database", func(t *testing.T) {
		secondDatabase, secondDatabaseCleanup := sdk.createDatabase(t, client)
		t.Cleanup(secondDatabaseCleanup)

		name := sdk.randomString(t)
		id := sdk.NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		newName := sdk.randomString(t)
		newId := sdk.NewDatabaseObjectIdentifier(secondDatabase.Name, newName)
		alterRequest := sdk.NewAlterDatabaseRoleRequest(id).WithRename(newId)

		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		assert.ErrorIs(t, err, sdk.errDifferentDatabase)
	})

	t.Run("show database_role: without like", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		role2 := createDatabaseRole(t)

		showRequest := sdk.NewShowDatabaseRoleRequest(database.ID())
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedDatabaseRoles))
		assert.Contains(t, returnedDatabaseRoles, *role1)
		assert.Contains(t, returnedDatabaseRoles, *role2)
	})

	t.Run("show database_role: with like", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		role2 := createDatabaseRole(t)

		showRequest := sdk.NewShowDatabaseRoleRequest(database.ID()).WithLike(role1.Name)
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedDatabaseRoles))
		assert.Contains(t, returnedDatabaseRoles, *role1)
		assert.NotContains(t, returnedDatabaseRoles, *role2)
	})

	t.Run("show database_role: no matches", func(t *testing.T) {
		showRequest := sdk.NewShowDatabaseRoleRequest(database.ID()).WithLike("non-existent")
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 0, len(returnedDatabaseRoles))
	})

	t.Run("grant and revoke database_role: to database role", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		id1 := sdk.NewDatabaseObjectIdentifier(database.Name, role1.Name)
		role2 := createDatabaseRole(t)
		id2 := sdk.NewDatabaseObjectIdentifier(database.Name, role2.Name)

		grantRequest := sdk.NewGrantDatabaseRoleRequest(id1).WithDatabaseRole(id2)
		err := client.DatabaseRoles.Grant(ctx, grantRequest)
		require.NoError(t, err)

		extractedRole, err := client.DatabaseRoles.ShowByID(ctx, id1)
		require.NoError(t, err)
		assert.Equal(t, 0, extractedRole.GrantedToRoles)
		assert.Equal(t, 1, extractedRole.GrantedToDatabaseRoles)
		assert.Equal(t, 0, extractedRole.GrantedDatabaseRoles)

		extractedRole, err = client.DatabaseRoles.ShowByID(ctx, id2)
		require.NoError(t, err)
		assert.Equal(t, 0, extractedRole.GrantedToRoles)
		assert.Equal(t, 0, extractedRole.GrantedToDatabaseRoles)
		assert.Equal(t, 1, extractedRole.GrantedDatabaseRoles)

		revokeRequest := sdk.NewRevokeDatabaseRoleRequest(id1).WithDatabaseRole(id2)
		err = client.DatabaseRoles.Revoke(ctx, revokeRequest)
		require.NoError(t, err)
	})

	t.Run("grant and revoke database_role: to account role", func(t *testing.T) {
		role := createDatabaseRole(t)
		roleId := sdk.NewDatabaseObjectIdentifier(database.Name, role.Name)

		accountRole, accountRoleCleanup := sdk.createRole(t, client)
		t.Cleanup(accountRoleCleanup)

		grantRequest := sdk.NewGrantDatabaseRoleRequest(roleId).WithAccountRole(accountRole.ID())
		err := client.DatabaseRoles.Grant(ctx, grantRequest)
		require.NoError(t, err)

		extractedRole, err := client.DatabaseRoles.ShowByID(ctx, roleId)
		require.NoError(t, err)
		assert.Equal(t, 1, extractedRole.GrantedToRoles)
		assert.Equal(t, 0, extractedRole.GrantedToDatabaseRoles)
		assert.Equal(t, 0, extractedRole.GrantedDatabaseRoles)

		revokeRequest := sdk.NewRevokeDatabaseRoleRequest(roleId).WithAccountRole(accountRole.ID())
		err = client.DatabaseRoles.Revoke(ctx, revokeRequest)
		require.NoError(t, err)
	})

	t.Run("grant and revoke database_role: to share", func(t *testing.T) {
		role := createDatabaseRole(t)
		roleId := sdk.NewDatabaseObjectIdentifier(database.Name, role.Name)

		share, shareCleanup := sdk.createShare(t, client)
		t.Cleanup(shareCleanup)

		err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{Database: database.ID()}, share.ID())
		require.NoError(t, err)

		grantRequest := sdk.NewGrantDatabaseRoleToShareRequest(roleId, share.ID())
		err = client.DatabaseRoles.GrantToShare(ctx, grantRequest)
		require.NoError(t, err)

		revokeRequest := sdk.NewRevokeDatabaseRoleFromShareRequest(roleId, share.ID())
		err = client.DatabaseRoles.RevokeFromShare(ctx, revokeRequest)
		require.NoError(t, err)
	})
}
