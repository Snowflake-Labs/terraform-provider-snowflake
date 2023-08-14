package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DatabaseRoles(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	assertDatabaseRole := func(t *testing.T, databaseRole *DatabaseRole, expectedName string, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, databaseRole.CreatedOn)
		assert.Equal(t, expectedName, databaseRole.Name)
		assert.Equal(t, "ACCOUNTADMIN", databaseRole.Owner)
		assert.Equal(t, expectedComment, databaseRole.Comment)
	}

	cleanupDatabaseRoleProvider := func(id DatabaseObjectIdentifier) func() {
		return func() {
			err := client.DatabaseRoles.Drop(ctx, NewDropDatabaseRoleRequest(id))
			require.NoError(t, err)
		}
	}

	createDatabaseRole := func(t *testing.T) *DatabaseRole {
		t.Helper()
		name := randomString(t)
		id := NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		return databaseRole
	}

	t.Run("create database_role: complete case", func(t *testing.T) {
		name := randomString(t)
		id := NewDatabaseObjectIdentifier(database.Name, name)
		comment := randomComment(t)

		request := NewCreateDatabaseRoleRequest(id).WithComment(&comment).WithIfNotExists(true)
		err := client.DatabaseRoles.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)

		require.NoError(t, err)
		assertDatabaseRole(t, databaseRole, name, comment)
	})

	t.Run("create database_role: no optionals", func(t *testing.T) {
		name := randomString(t)
		id := NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assertDatabaseRole(t, databaseRole, name, "")
	})

	t.Run("drop database_role: existing", func(t *testing.T) {
		name := randomString(t)
		id := NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)

		err = client.DatabaseRoles.Drop(ctx, NewDropDatabaseRoleRequest(id))
		require.NoError(t, err)

		_, err = client.DatabaseRoles.ShowByID(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop database_role: non-existing", func(t *testing.T) {
		id := NewDatabaseObjectIdentifier(database.Name, "does_not_exist")

		err := client.DatabaseRoles.Drop(ctx, NewDropDatabaseRoleRequest(id))
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter database_role: set value and unset value", func(t *testing.T) {
		name := randomString(t)
		id := NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		alterRequest := NewAlterDatabaseRoleRequest(id).WithSet(NewDatabaseRoleSetRequest("new comment"))
		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredDatabaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredDatabaseRole.Comment)

		alterRequest = NewAlterDatabaseRoleRequest(id).WithUnsetComment()
		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredDatabaseRole, err = client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredDatabaseRole.Comment)
	})

	t.Run("alter database_role: rename", func(t *testing.T) {
		name := randomString(t)
		id := NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)

		newName := randomString(t)
		newId := NewDatabaseObjectIdentifier(database.Name, newName)
		alterRequest := NewAlterDatabaseRoleRequest(id).WithRename(NewDatabaseRoleRenameRequest(newId))

		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupDatabaseRoleProvider(id))
		} else {
			t.Cleanup(cleanupDatabaseRoleProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.DatabaseRoles.ShowByID(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertDatabaseRole(t, databaseRole, newName, "")
	})

	t.Run("alter database_role: rename to other database", func(t *testing.T) {
		secondDatabase, secondDatabaseCleanup := createDatabase(t, client)
		t.Cleanup(secondDatabaseCleanup)

		name := randomString(t)
		id := NewDatabaseObjectIdentifier(database.Name, name)

		err := client.DatabaseRoles.Create(ctx, NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		newName := randomString(t)
		newId := NewDatabaseObjectIdentifier(secondDatabase.Name, newName)
		alterRequest := NewAlterDatabaseRoleRequest(id).WithRename(NewDatabaseRoleRenameRequest(newId))

		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		assert.ErrorIs(t, err, errDifferentDatabase)
	})

	t.Run("show database_role: without like", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		role2 := createDatabaseRole(t)

		showRequest := NewShowDatabaseRoleRequest(database.ID())
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedDatabaseRoles))
		assert.Contains(t, returnedDatabaseRoles, role1)
		assert.Contains(t, returnedDatabaseRoles, role2)
	})

	t.Run("show database_role: with like", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		role2 := createDatabaseRole(t)

		showRequest := NewShowDatabaseRoleRequest(database.ID()).WithLike(role1.Name)
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedDatabaseRoles))
		assert.Contains(t, returnedDatabaseRoles, role1)
		assert.NotContains(t, returnedDatabaseRoles, role2)
	})

	t.Run("show database_role: no matches", func(t *testing.T) {
		showRequest := NewShowDatabaseRoleRequest(database.ID()).WithLike("non-existent")
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 0, len(returnedDatabaseRoles))
	})
}
