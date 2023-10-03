package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Roles(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schema, _ := createSchema(t, client, database)
	tag, _ := createTag(t, client, database, schema)
	tag2, _ := createTag(t, client, database, schema)

	t.Run("create no options", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		err := client.Roles.Create(ctx, NewCreateRoleRequest(roleID))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, NewDropRoleRequest(roleID))
			require.NoError(t, err)
		})

		role, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(roleID))
		require.NoError(t, err)

		assert.Equal(t, roleID.Name(), role.Name)
	})

	t.Run("create if not exists", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		err := client.Roles.Create(ctx, NewCreateRoleRequest(roleID).WithIfNotExists(true))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, NewDropRoleRequest(roleID))
			require.NoError(t, err)
		})

		role, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(roleID))
		require.NoError(t, err)
		assert.Equal(t, roleID.Name(), role.Name)
	})

	t.Run("create complete", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		comment := randomComment(t)
		createReq := NewCreateRoleRequest(roleID).
			WithOrReplace(true).
			WithTag([]TagAssociation{
				{
					Name:  tag.ID(),
					Value: "v1",
				},
				{
					Name:  tag2.ID(),
					Value: "v2",
				},
			}).
			WithComment(comment)
		err := client.Roles.Create(ctx, createReq)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, NewDropRoleRequest(roleID))
			require.NoError(t, err)
		})

		role, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(roleID))
		require.NoError(t, err)
		assert.Equal(t, roleID.Name(), role.Name)
		assert.Equal(t, comment, role.Comment)

		// verify tags
		tag1Value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, "v1", tag1Value)

		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2.ID(), role.ID(), ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, "v2", tag2Value)
	})

	t.Run("alter rename to", func(t *testing.T) {
		role, _ := createRole(t, client)
		newName := randomAccountObjectIdentifier(t)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, NewDropRoleRequest(newName))
			if err != nil {
				err = client.Roles.Drop(ctx, NewDropRoleRequest(role.ID()))
				require.NoError(t, err)
			}
		})

		err := client.Roles.Alter(ctx, NewAlterRoleRequest(role.ID()).WithRenameTo(newName))
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(newName))
		require.NoError(t, err)
		assert.Equal(t, newName.Name(), r.Name)
	})

	t.Run("alter set tags", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		_, err := client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), "ROLE")
		require.Error(t, err)

		tagValue := "new-tag-value"
		err = client.Roles.Alter(ctx, NewAlterRoleRequest(role.ID()).WithSetTags([]TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}))
		require.NoError(t, err)

		addedTag, err := client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, tagValue, addedTag)
	})

	t.Run("alter unset tags", func(t *testing.T) {
		tagValue := "tag-value"
		id := randomAccountObjectIdentifier(t)
		role, cleanup := createRoleWithRequest(t, client, NewCreateRoleRequest(id).
			WithTag([]TagAssociation{
				{
					Name:  tag.ID(),
					Value: tagValue,
				},
			}))
		t.Cleanup(cleanup)

		value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, tagValue, value)

		err = client.Roles.Alter(ctx, NewAlterRoleRequest(role.ID()).WithUnsetTags([]ObjectIdentifier{tag.ID()}))
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), ObjectTypeRole)
		require.Error(t, err)
	})

	t.Run("alter set comment", func(t *testing.T) {
		role, cleanupRole := createRole(t, client)
		t.Cleanup(cleanupRole)

		comment := randomComment(t)
		err := client.Roles.Alter(ctx, NewAlterRoleRequest(role.ID()).WithSetComment(comment))
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		assert.Equal(t, comment, r.Comment)
	})

	t.Run("alter unset comment", func(t *testing.T) {
		comment := randomComment(t)
		id := randomAccountObjectIdentifier(t)
		role, cleanup := createRoleWithRequest(t, client, NewCreateRoleRequest(id).WithComment(comment))
		t.Cleanup(cleanup)

		err := client.Roles.Alter(ctx, NewAlterRoleRequest(role.ID()).WithUnsetComment(true))
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		assert.Equal(t, "", r.Comment)
	})

	t.Run("drop no options", func(t *testing.T) {
		role, _ := createRole(t, client)
		err := client.Roles.Drop(ctx, NewDropRoleRequest(role.ID()))
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(role.ID()))
		require.Nil(t, r)
		require.Error(t, err)
	})

	t.Run("show no options", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		role2, cleanup2 := createRole(t, client)
		t.Cleanup(cleanup2)

		roles, err := client.Roles.Show(ctx, NewShowRoleRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(roles), 2)

		roleIDs := make([]AccountObjectIdentifier, len(roles))
		for i, r := range roles {
			roleIDs[i] = r.ID()
		}
		assert.Contains(t, roleIDs, role.ID())
		assert.Contains(t, roleIDs, role2.ID())
	})

	t.Run("show like", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		roles, err := client.Roles.Show(ctx, NewShowRoleRequest().WithLike(NewLikeRequest(role.Name)))
		require.NoError(t, err)
		assert.Equal(t, 1, len(roles))
		assert.Equal(t, role.Name, roles[0].Name)
	})

	t.Run("in class", func(t *testing.T) {
		roles, err := client.Roles.Show(ctx, NewShowRoleRequest().WithInClass(RolesInClass{
			Class: NewSchemaObjectIdentifier("SNOWFLAKE", "ML", "ANOMALY_DETECTION"),
		}))
		require.NoError(t, err)
		assert.Equal(t, 1, len(roles))
		assert.Equal(t, "USER", roles[0].Name)
	})

	t.Run("show by id", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		r, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, role.Name, r.Name)
	})

	t.Run("grant and revoke role from user", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		user, cleanupUser := createUser(t, client)
		t.Cleanup(cleanupUser)

		userID := user.ID()
		err := client.Roles.Grant(ctx, NewGrantRoleRequest(role.ID(), GrantRole{User: &userID}))
		require.NoError(t, err)

		roleBefore, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		assert.Equal(t, 1, roleBefore.AssignedToUsers)

		err = client.Roles.Revoke(ctx, NewRevokeRoleRequest(role.ID(), RevokeRole{User: &userID}))
		require.NoError(t, err)

		roleAfter, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		assert.Equal(t, 0, roleAfter.AssignedToUsers)
	})

	t.Run("grant and revoke role from role", func(t *testing.T) {
		parentRole, cleanupParentRole := createRole(t, client)
		t.Cleanup(cleanupParentRole)

		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		parentRoleID := parentRole.ID()
		err := client.Roles.Grant(ctx, NewGrantRoleRequest(role.ID(), GrantRole{Role: &parentRoleID}))
		require.NoError(t, err)

		roleBefore, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)

		parentRoleBefore, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(parentRole.ID()))
		require.NoError(t, err)

		require.Equal(t, 1, roleBefore.GrantedToRoles)
		require.Equal(t, 1, parentRoleBefore.GrantedRoles)

		err = client.Roles.Revoke(ctx, NewRevokeRoleRequest(role.ID(), RevokeRole{Role: &parentRoleID}))
		require.NoError(t, err)

		roleAfter, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)

		parentRoleAfter, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(parentRole.ID()))
		require.NoError(t, err)

		assert.Equal(t, 0, roleAfter.GrantedToRoles)
		assert.Equal(t, 0, parentRoleAfter.GrantedRoles)
	})
}
