package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Roles(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schema, _ := createSchema(t, client, database)
	tag, _ := createTag(t, client, database, schema)
	tag2, _ := createTag(t, client, database, schema)

	t.Run("create no options", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		err := client.Roles.Create(ctx, sdk.NewCreateRoleRequest(roleID))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, sdk.NewDropRoleRequest(roleID))
			require.NoError(t, err)
		})

		role, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(roleID))
		require.NoError(t, err)

		assert.Equal(t, roleID.Name(), role.Name)
	})

	t.Run("create if not exists", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		err := client.Roles.Create(ctx, sdk.NewCreateRoleRequest(roleID).WithIfNotExists(true))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, sdk.NewDropRoleRequest(roleID))
			require.NoError(t, err)
		})

		role, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(roleID))
		require.NoError(t, err)
		assert.Equal(t, roleID.Name(), role.Name)
	})

	t.Run("create complete", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		comment := random.RandomComment(t)
		createReq := sdk.NewCreateRoleRequest(roleID).
			WithOrReplace(true).
			WithTag([]sdk.TagAssociation{
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
			err := client.Roles.Drop(ctx, sdk.NewDropRoleRequest(roleID))
			require.NoError(t, err)
		})

		role, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(roleID))
		require.NoError(t, err)
		assert.Equal(t, roleID.Name(), role.Name)
		assert.Equal(t, comment, role.Comment)

		// verify tags
		tag1Value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), sdk.ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, "v1", tag1Value)

		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2.ID(), role.ID(), sdk.ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, "v2", tag2Value)
	})

	t.Run("alter rename to", func(t *testing.T) {
		role, _ := createRole(t, client)
		newName := randomAccountObjectIdentifier(t)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, sdk.NewDropRoleRequest(newName))
			if err != nil {
				err = client.Roles.Drop(ctx, sdk.NewDropRoleRequest(role.ID()))
				require.NoError(t, err)
			}
		})

		err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(role.ID()).WithRenameTo(newName))
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(newName))
		require.NoError(t, err)
		assert.Equal(t, newName.Name(), r.Name)
	})

	t.Run("alter set tags", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		_, err := client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), "ROLE")
		require.Error(t, err)

		tagValue := "new-tag-value"
		err = client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(role.ID()).WithSetTags([]sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}))
		require.NoError(t, err)

		addedTag, err := client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), sdk.ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, tagValue, addedTag)
	})

	t.Run("alter unset tags", func(t *testing.T) {
		tagValue := "tag-value"
		id := randomAccountObjectIdentifier(t)
		role, cleanup := createRoleWithRequest(t, client, sdk.NewCreateRoleRequest(id).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: tagValue,
				},
			}))
		t.Cleanup(cleanup)

		value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), sdk.ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, tagValue, value)

		err = client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(role.ID()).WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}))
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), role.ID(), sdk.ObjectTypeRole)
		require.Error(t, err)
	})

	t.Run("alter set comment", func(t *testing.T) {
		role, cleanupRole := createRole(t, client)
		t.Cleanup(cleanupRole)

		comment := random.RandomComment(t)
		err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(role.ID()).WithSetComment(comment))
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		assert.Equal(t, comment, r.Comment)
	})

	t.Run("alter unset comment", func(t *testing.T) {
		comment := random.RandomComment(t)
		id := randomAccountObjectIdentifier(t)
		role, cleanup := createRoleWithRequest(t, client, sdk.NewCreateRoleRequest(id).WithComment(comment))
		t.Cleanup(cleanup)

		err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(role.ID()).WithUnsetComment(true))
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		assert.Equal(t, "", r.Comment)
	})

	t.Run("drop no options", func(t *testing.T) {
		role, _ := createRole(t, client)
		err := client.Roles.Drop(ctx, sdk.NewDropRoleRequest(role.ID()))
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(role.ID()))
		require.Nil(t, r)
		require.Error(t, err)
	})

	t.Run("show no options", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		role2, cleanup2 := createRole(t, client)
		t.Cleanup(cleanup2)

		roles, err := client.Roles.Show(ctx, sdk.NewShowRoleRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(roles), 2)

		roleIDs := make([]sdk.AccountObjectIdentifier, len(roles))
		for i, r := range roles {
			roleIDs[i] = r.ID()
		}
		assert.Contains(t, roleIDs, role.ID())
		assert.Contains(t, roleIDs, role2.ID())
	})

	t.Run("show like", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		roles, err := client.Roles.Show(ctx, sdk.NewShowRoleRequest().WithLike(sdk.NewLikeRequest(role.Name)))
		require.NoError(t, err)
		assert.Equal(t, 1, len(roles))
		assert.Equal(t, role.Name, roles[0].Name)
	})

	t.Run("in class", func(t *testing.T) {
		roles, err := client.Roles.Show(ctx, sdk.NewShowRoleRequest().WithInClass(sdk.RolesInClass{
			Class: sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "ML", "ANOMALY_DETECTION"),
		}))
		require.NoError(t, err)
		assert.Equal(t, 1, len(roles))
		assert.Equal(t, "USER", roles[0].Name)
	})

	t.Run("show by id", func(t *testing.T) {
		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		r, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(role.ID()))
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
		err := client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{User: &userID}))
		require.NoError(t, err)

		roleBefore, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		assert.Equal(t, 1, roleBefore.AssignedToUsers)

		err = client.Roles.Revoke(ctx, sdk.NewRevokeRoleRequest(role.ID(), sdk.RevokeRole{User: &userID}))
		require.NoError(t, err)

		roleAfter, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)
		assert.Equal(t, 0, roleAfter.AssignedToUsers)
	})

	t.Run("grant and revoke role from role", func(t *testing.T) {
		parentRole, cleanupParentRole := createRole(t, client)
		t.Cleanup(cleanupParentRole)

		role, cleanup := createRole(t, client)
		t.Cleanup(cleanup)

		parentRoleID := parentRole.ID()
		err := client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{Role: &parentRoleID}))
		require.NoError(t, err)

		roleBefore, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)

		parentRoleBefore, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(parentRole.ID()))
		require.NoError(t, err)

		require.Equal(t, 1, roleBefore.GrantedToRoles)
		require.Equal(t, 1, parentRoleBefore.GrantedRoles)

		err = client.Roles.Revoke(ctx, sdk.NewRevokeRoleRequest(role.ID(), sdk.RevokeRole{Role: &parentRoleID}))
		require.NoError(t, err)

		roleAfter, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(role.ID()))
		require.NoError(t, err)

		parentRoleAfter, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(parentRole.ID()))
		require.NoError(t, err)

		assert.Equal(t, 0, roleAfter.GrantedToRoles)
		assert.Equal(t, 0, parentRoleAfter.GrantedRoles)
	})
}
