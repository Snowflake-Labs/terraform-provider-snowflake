package sdk_integration_tests

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_UsersShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	userTest, userCleanup := sdk.createUserWithName(t, client, "USER_FOO")
	t.Cleanup(userCleanup)

	userTest2, user2Cleanup := sdk.createUserWithName(t, client, "USER_BAR")
	t.Cleanup(user2Cleanup)

	t.Run("with like options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(userTest.Name),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *userTest)
		assert.Equal(t, 1, len(users))
	})

	t.Run("with starts with options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			StartsWith: sdk.String("USER"),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *userTest)
		assert.Contains(t, users, *userTest2)
		assert.Equal(t, 2, len(users))
	})
	t.Run("with starts with, limit and from options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Limit:      sdk.Int(10),
			From:       sdk.String("USER_C"),
			StartsWith: sdk.String("USER"),
		}

		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *userTest)
		assert.Equal(t, 1, len(users))
	})

	t.Run("when searching a non-existent user", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(users))
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Limit: sdk.Int(1),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})
}

func TestInt_UserCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := sdk.createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := sdk.createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	tag, tagCleanup := sdk.createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleanup)

	t.Run("test complete case", func(t *testing.T) {
		id := sdk.randomAccountObjectIdentifier(t)
		tagValue := sdk.randomString(t)
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := sdk.randomString(t)
		loginName := sdk.randomString(t)

		opts := &sdk.CreateUserOptions{
			OrReplace: sdk.Bool(true),
			name:      id,
			ObjectProperties: &sdk.UserObjectProperties{
				Password:  &password,
				LoginName: &loginName,
			},
			ObjectParameters: &sdk.UserObjectParameters{
				EnableUnredactedQuerySyntaxError: sdk.Bool(true),
			},
			SessionParameters: &sdk.SessionParameters{
				Autocommit: sdk.Bool(true),
			},
			With: sdk.Bool(true),
			Tags: tags,
		}
		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.name, userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)

		user, err := client.Users.Show(ctx, &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: &id.name,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(user))
		assert.Equal(t, id.name, user[0].Name)
	})

	t.Run("test if not exists", func(t *testing.T) {
		id := sdk.randomAccountObjectIdentifier(t)
		tagValue := sdk.randomString(t)
		tags := []sdk.TagAssociation{
			{
				Name:  sdk.NewAccountObjectIdentifier(tag.Name),
				Value: tagValue,
			},
		}
		password := sdk.randomString(t)
		loginName := sdk.randomString(t)

		opts := &sdk.CreateUserOptions{
			IfNotExists: sdk.Bool(true),
			name:        id,
			ObjectProperties: &sdk.UserObjectProperties{
				Password:  &password,
				LoginName: &loginName,
			},
			ObjectParameters: &sdk.UserObjectParameters{
				EnableUnredactedQuerySyntaxError: sdk.Bool(true),
			},
			SessionParameters: &sdk.SessionParameters{
				Autocommit: sdk.Bool(true),
			},
			With: sdk.Bool(true),
			Tags: tags,
		}
		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.name, userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)

		user, err := client.Users.Show(ctx, &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: &id.name,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(user))
		assert.Equal(t, id.name, user[0].Name)
	})

	t.Run("test no options", func(t *testing.T) {
		id := sdk.randomAccountObjectIdentifier(t)

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.name, userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(id.name), userDetails.LoginName.Value)
		assert.Empty(t, userDetails.Password.Value)

		user, err := client.Users.Show(ctx, &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: &id.name,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(user))
		assert.Equal(t, id.name, user[0].Name)
	})
}

func TestInt_UserDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	user, userCleanup := sdk.createUser(t, client)
	t.Cleanup(userCleanup)

	t.Run("when user exists", func(t *testing.T) {
		userDetails, err := client.Users.Describe(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.Name, userDetails.Name.Value)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")
		_, err := client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})
}

func TestInt_UserDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when user exists", func(t *testing.T) {
		user, _ := sdk.createUser(t, client)
		id := user.ID()
		err := client.Users.Drop(ctx, id)
		require.NoError(t, err)
		_, err = client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")
		err := client.Users.Drop(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})
}
