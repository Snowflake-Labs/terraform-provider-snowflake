package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_UsersShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	userTest, userCleanup := createUserWithName(t, client, "USER_FOO")
	t.Cleanup(userCleanup)

	userTest2, user2Cleanup := createUserWithName(t, client, "USER_BAR")
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

	tag, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
	t.Cleanup(tagCleanup)

	t.Run("test complete case", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()
		tagValue := random.String()
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := random.String()
		loginName := random.String()

		opts := &sdk.CreateUserOptions{
			OrReplace: sdk.Bool(true),
			ObjectProperties: &sdk.UserObjectProperties{
				Password:    &password,
				LoginName:   &loginName,
				DefaultRole: sdk.String("foo"),
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
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)
		assert.Equal(t, "FOO", userDetails.DefaultRole.Value)

		user, err := client.Users.Show(ctx, &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.Pointer(id.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(user))
		assert.Equal(t, id.Name(), user[0].Name)
	})

	t.Run("test if not exists", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()
		tagValue := random.String()
		tags := []sdk.TagAssociation{
			{
				Name:  sdk.NewAccountObjectIdentifier(tag.Name),
				Value: tagValue,
			},
		}
		password := random.String()
		loginName := random.String()

		opts := &sdk.CreateUserOptions{
			IfNotExists: sdk.Bool(true),
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
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)

		user, err := client.Users.Show(ctx, &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.Pointer(id.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(user))
		assert.Equal(t, id.Name(), user[0].Name)
	})

	t.Run("test no options", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(id.Name()), userDetails.LoginName.Value)
		assert.Empty(t, userDetails.Password.Value)

		user, err := client.Users.Show(ctx, &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.Pointer(id.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(user))
		assert.Equal(t, id.Name(), user[0].Name)
	})
}

func TestInt_UserDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	user, userCleanup := createUser(t, client)
	t.Cleanup(userCleanup)

	t.Run("when user exists", func(t *testing.T) {
		userDetails, err := client.Users.Describe(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.Name, userDetails.Name.Value)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")
		_, err := client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_UserDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when user exists", func(t *testing.T) {
		user, _ := createUser(t, client)
		id := user.ID()
		err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{})
		require.NoError(t, err)
		_, err = client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")
		err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{})
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
