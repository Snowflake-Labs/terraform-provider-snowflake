package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_RolesCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("no options", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		err := client.Roles.Create(ctx, roleID, nil)
		require.NoError(t, err)
		role, err := client.Roles.ShowByID(ctx, roleID)
		require.NoError(t, err)
		assert.Equal(t, roleID.Name(), role.Name)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, roleID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("if not exists", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		opts := &RoleCreateOptions{
			IfNotExists: Bool(true),
		}
		err := client.Roles.Create(ctx, roleID, opts)
		require.NoError(t, err)
		role, err := client.Roles.ShowByID(ctx, roleID)
		require.NoError(t, err)
		assert.Equal(t, roleID.Name(), role.Name)
		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, roleID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("complete test case", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)

		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
		t.Cleanup(schemaCleanup)
		tagTest, tagCleanup := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tagCleanup)
		tag2Test, tag2Cleanup := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tag2Cleanup)
		comment := randomComment(t)

		opts := &RoleCreateOptions{
			OrReplace: Bool(true),
			Tag: []TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
				{
					Name:  tag2Test.ID(),
					Value: "v2",
				},
			},
			Comment: String(comment),
		}
		err := client.Roles.Create(ctx, roleID, opts)
		require.NoError(t, err)
		role, err := client.Roles.ShowByID(ctx, roleID)
		require.NoError(t, err)
		assert.Equal(t, roleID.Name(), role.Name)
		assert.Equal(t, comment, role.Comment)

		// verify tags
		tag1Value, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), role.ID(), ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, "v1", tag1Value)
		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2Test.ID(), role.ID(), ObjectTypeRole)
		require.NoError(t, err)
		assert.Equal(t, "v2", tag2Value)

		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, roleID, nil)
			require.NoError(t, err)
		})
	})
}

func TestInt_RolesAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, cleanupDatabase := createDatabase(t, client)
	t.Cleanup(cleanupDatabase)

	schema, cleanupSchema := createSchema(t, client, database)
	t.Cleanup(cleanupSchema)

	tag, cleanupTag := createTag(t, client, database, schema)
	t.Cleanup(cleanupTag)

	t.Run("renaming", func(t *testing.T) {
		role, _ := createRole(t, client)
		newName := randomAccountObjectIdentifier(t)

		err := client.Roles.Alter(ctx, role.ID(), &RoleAlterOptions{
			RenameTo: newName,
		})
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, newName)
		assert.Equal(t, newName.Name(), r.Name)

		t.Cleanup(func() {
			err = client.Roles.Drop(ctx, r.ID(), nil)
			require.NoError(t, err)
		})
	})

	t.Run("setting tags", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		err := client.Roles.Create(ctx, roleID, nil)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), roleID, "ROLE")
		require.Error(t, err)

		tagValue := "new-tag-value"
		err = client.Roles.Alter(ctx, roleID, &RoleAlterOptions{
			Set: &RoleSet{
				Tag: []TagAssociation{
					{
						Name:  tag.ID(),
						Value: tagValue,
					},
				},
			},
		})
		require.NoError(t, err)

		addedTag, err := client.SystemFunctions.GetTag(ctx, tag.ID(), roleID, "ROLE")
		require.NoError(t, err)
		assert.Equal(t, tagValue, addedTag)

		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, roleID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("unsetting tags", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		tagValue := "tagvalue"
		err := client.Roles.Create(ctx, roleID, &RoleCreateOptions{
			Tag: []TagAssociation{
				{
					Name:  tag.ID(),
					Value: tagValue,
				},
			},
		})
		require.NoError(t, err)

		value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), roleID, "ROLE")
		require.NoError(t, err)
		assert.Equal(t, tagValue, value)

		err = client.Roles.Alter(ctx, roleID, &RoleAlterOptions{
			Unset: &RoleUnset{
				Tag: []ObjectIdentifier{tag.ID()},
			},
		})
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), roleID, "ROLE")
		require.Error(t, err)

		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, roleID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("setting comment", func(t *testing.T) {
		role, cleanupRole := createRole(t, client)
		t.Cleanup(cleanupRole)

		comment := randomComment(t)
		err := client.Roles.Alter(ctx, role.ID(), &RoleAlterOptions{
			Set: &RoleSet{
				Comment: &comment,
			},
		})
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, role.ID())
		require.NoError(t, err)
		assert.Equal(t, comment, r.Comment)
	})

	t.Run("unsetting comment", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		comment := randomComment(t)
		err := client.Roles.Create(ctx, roleID, &RoleCreateOptions{
			Comment: &comment,
		})
		require.NoError(t, err)

		err = client.Roles.Alter(ctx, roleID, &RoleAlterOptions{
			Unset: &RoleUnset{
				Comment: Bool(true),
			},
		})
		require.NoError(t, err)

		role, err := client.Roles.ShowByID(ctx, roleID)
		require.NoError(t, err)
		assert.Equal(t, "", role.Comment)

		t.Cleanup(func() {
			err := client.Roles.Drop(ctx, roleID, nil)
			require.NoError(t, err)
		})
	})
}

func TestInt_RolesDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	role, _ := createRole(t, client)
	roleID := role.ID()

	t.Run("drop with nil options", func(t *testing.T) {
		err := client.Roles.Drop(ctx, roleID, nil)
		require.NoError(t, err)
	})
}

func TestInt_RolesShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	role, cleanup := createRole(t, client)
	t.Cleanup(cleanup)

	role2, cleanup2 := createRole(t, client)
	t.Cleanup(cleanup2)

	t.Run("no options", func(t *testing.T) {
		roles, err := client.Roles.Show(ctx, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(roles), 2)

		roleIDs := make([]AccountObjectIdentifier, len(roles))
		for i, r := range roles {
			roleIDs[i] = r.ID()
		}
		assert.Contains(t, roleIDs, role.ID())
		assert.Contains(t, roleIDs, role2.ID())
	})

	t.Run("with like", func(t *testing.T) {
		roles, err := client.Roles.Show(ctx, &RoleShowOptions{
			Like: &Like{
				Pattern: String(role.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(roles))
		assert.Equal(t, role.Name, roles[0].Name)
	})

	t.Run("by id", func(t *testing.T) {
		r, err := client.Roles.ShowByID(ctx, role.ID())
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, role.Name, r.Name)
	})
}

func TestInt_RolesGrant(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	parent_role, cleanup_parent_role := createRole(t, client)
	parent_roleID := parent_role.ID()
	t.Cleanup(cleanup_parent_role)

	role, cleanup_role := createRole(t, client)
	roleID := role.ID()
	t.Cleanup(cleanup_role)

	// user, cleanup_user := createUser(t, client)
	// t.Cleanup(cleanup_user)

	t.Run("grant role to user", func(t *testing.T) {
		// TODO: Wait for Users pr
	})

	t.Run("grant role to role", func(t *testing.T) {
		require.Equal(t, 0, role.GrantedToRoles)
		require.Equal(t, 0, parent_role.GrantedToRoles)

		opts := RoleGrantOptions{
			Grant: GrantRole{
				Role: &parent_roleID,
			},
		}
		err := client.Roles.Grant(ctx, roleID, &opts)
		require.NoError(t, err)

		r, err := client.Roles.ShowByID(ctx, roleID)
		require.NotNil(t, r)
		require.NoError(t, err)

		pr, err := client.Roles.ShowByID(ctx, parent_roleID)
		require.NotNil(t, pr)
		require.NoError(t, err)

		assert.Equal(t, 1, r.GrantedToRoles)
		assert.Equal(t, 1, pr.GrantedRoles)
	})
}

func TestInt_RolesRevoke(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	parent_role, cleanup_parent_role := createRole(t, client)
	parent_roleID := parent_role.ID()
	t.Cleanup(cleanup_parent_role)

	role, cleanup_role := createRole(t, client)
	roleID := role.ID()
	t.Cleanup(cleanup_role)

	err := client.Roles.Grant(ctx, roleID, &RoleGrantOptions{
		Grant: GrantRole{
			Role: &parent_roleID,
		},
	})
	require.NoError(t, err)

	t.Run("revoke role from user", func(t *testing.T) {
		// TODO: Wait for Users pr
	})

	t.Run("revoke role from role", func(t *testing.T) {
		role_before, err := client.Roles.ShowByID(ctx, roleID)
		require.NotNil(t, role_before)
		require.NoError(t, err)

		parent_role_before, err := client.Roles.ShowByID(ctx, parent_roleID)
		require.NotNil(t, parent_role_before)
		require.NoError(t, err)

		require.Equal(t, 1, role_before.GrantedToRoles)
		require.Equal(t, 1, parent_role_before.GrantedRoles)

		err = client.Roles.Revoke(ctx, roleID, &RoleRevokeOptions{
			Revoke: RevokeRole{
				Role: &parent_roleID,
			},
		})
		require.NoError(t, err)

		role_after, err := client.Roles.ShowByID(ctx, roleID)
		require.NotNil(t, role_after)
		require.NoError(t, err)

		parent_role_after, err := client.Roles.ShowByID(ctx, parent_roleID)
		require.NotNil(t, parent_role_after)
		require.NoError(t, err)

		assert.Equal(t, 0, role_after.GrantedToRoles)
		assert.Equal(t, 0, parent_role_after.GrantedRoles)
	})
}
