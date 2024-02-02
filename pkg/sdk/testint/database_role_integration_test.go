package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DatabaseRoles(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

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
		name := random.String()
		id := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		return databaseRole
	}

	t.Run("create database_role: complete case", func(t *testing.T) {
		name := random.String()
		id := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, name)
		comment := random.Comment()

		request := sdk.NewCreateDatabaseRoleRequest(id).WithComment(&comment).WithIfNotExists(true)
		err := client.DatabaseRoles.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)

		require.NoError(t, err)
		assertDatabaseRole(t, databaseRole, name, comment)
	})

	t.Run("create database_role: no optionals", func(t *testing.T) {
		name := random.String()
		id := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assertDatabaseRole(t, databaseRole, name, "")
	})

	t.Run("drop database_role: existing", func(t *testing.T) {
		name := random.String()
		id := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)

		err = client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
		require.NoError(t, err)

		_, err = client.DatabaseRoles.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop database_role: non-existing", func(t *testing.T) {
		id := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, "does_not_exist")

		err := client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter database_role: set value and unset value", func(t *testing.T) {
		name := random.String()
		id := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, name)

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
		name := random.String()
		id := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)

		newName := random.String()
		newId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, newName)
		alterRequest := sdk.NewAlterDatabaseRoleRequest(id).WithRename(newId)

		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupDatabaseRoleProvider(id))
		} else {
			t.Cleanup(cleanupDatabaseRoleProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.DatabaseRoles.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertDatabaseRole(t, databaseRole, newName, "")
	})

	t.Run("alter database_role: rename to other database", func(t *testing.T) {
		secondDatabase, secondDatabaseCleanup := createDatabase(t, client)
		t.Cleanup(secondDatabaseCleanup)

		name := random.String()
		id := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, name)

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		newName := random.String()
		newId := sdk.NewDatabaseObjectIdentifier(secondDatabase.Name, newName)
		alterRequest := sdk.NewAlterDatabaseRoleRequest(id).WithRename(newId)

		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		assert.ErrorIs(t, err, sdk.ErrDifferentDatabase)
	})

	t.Run("show database_role: without like", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		role2 := createDatabaseRole(t)

		showRequest := sdk.NewShowDatabaseRoleRequest(testDb(t).ID())
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedDatabaseRoles))
		assert.Contains(t, returnedDatabaseRoles, *role1)
		assert.Contains(t, returnedDatabaseRoles, *role2)
	})

	t.Run("show database_role: with like", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		role2 := createDatabaseRole(t)

		showRequest := sdk.NewShowDatabaseRoleRequest(testDb(t).ID()).WithLike(role1.Name)
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedDatabaseRoles))
		assert.Contains(t, returnedDatabaseRoles, *role1)
		assert.NotContains(t, returnedDatabaseRoles, *role2)
	})

	t.Run("show database_role: no matches", func(t *testing.T) {
		showRequest := sdk.NewShowDatabaseRoleRequest(testDb(t).ID()).WithLike("non-existent")
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 0, len(returnedDatabaseRoles))
	})

	t.Run("grant and revoke database_role: to database role", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		id1 := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, role1.Name)
		role2 := createDatabaseRole(t)
		id2 := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, role2.Name)

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
		roleId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, role.Name)

		accountRole, accountRoleCleanup := createRole(t, client)
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
		roleId := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, role.Name)

		share, shareCleanup := createShare(t, client)
		t.Cleanup(shareCleanup)

		err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{Database: testDb(t).ID()}, share.ID())
		require.NoError(t, err)

		grantRequest := sdk.NewGrantDatabaseRoleToShareRequest(roleId, share.ID())
		err = client.DatabaseRoles.GrantToShare(ctx, grantRequest)
		require.NoError(t, err)

		revokeRequest := sdk.NewRevokeDatabaseRoleFromShareRequest(roleId, share.ID())
		err = client.DatabaseRoles.RevokeFromShare(ctx, revokeRequest)
		require.NoError(t, err)
	})
}
