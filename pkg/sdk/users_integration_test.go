package sdk

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_UsersShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	userTest, userCleanup := createUserWithName(t, client, "USER_FOO")
	t.Cleanup(userCleanup)

	userTest2, user2Cleanup := createUserWithName(t, client, "USER_BAR")
	t.Cleanup(user2Cleanup)

	t.Run("with like options", func(t *testing.T) {
		showOptions := &ShowUserOptions{
			Like: &Like{
				Pattern: String(userTest.Name),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, userTest)
		assert.Equal(t, 1, len(users))
	})

	t.Run("with starts with options", func(t *testing.T) {
		showOptions := &ShowUserOptions{
			StartsWith: String("USER"),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, userTest)
		assert.Contains(t, users, userTest2)
		assert.Equal(t, 2, len(users))
	})
	t.Run("with starts with, limit and from options", func(t *testing.T) {
		showOptions := &ShowUserOptions{
			Limit:      Int(10),
			From:       String("USER_C"),
			StartsWith: String("USER"),
		}

		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, userTest)
		assert.Equal(t, 1, len(users))
	})

	t.Run("when searching a non-existent user", func(t *testing.T) {
		showOptions := &ShowUserOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(users))
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &ShowUserOptions{
			Limit: Int(1),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})
}

func TestInt_UserCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	tag, tagCleanup := createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleanup)

	t.Run("test complete case", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)
		tagValue := randomString(t)
		tags := []TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := randomString(t)
		loginName := randomString(t)

		opts := &CreateUserOptions{
			OrReplace: Bool(true),
			name:      id,
			ObjectProperties: &UserObjectProperties{
				Password:  &password,
				LoginName: &loginName,
			},
			ObjectParameters: &UserObjectParameters{
				EnableUnredactedQuerySyntaxError: Bool(true),
			},
			SessionParameters: &SessionParameters{
				Autocommit: Bool(true),
			},
			With: Bool(true),
			Tags: tags,
		}
		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.name, userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)

		user, err := client.Users.Show(ctx, &ShowUserOptions{
			Like: &Like{
				Pattern: &id.name,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(user))
		assert.Equal(t, id.name, user[0].Name)
	})

	t.Run("test if not exists", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)
		tagValue := randomString(t)
		tags := []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier(tag.Name),
				Value: tagValue,
			},
		}
		password := randomString(t)
		loginName := randomString(t)

		opts := &CreateUserOptions{
			IfNotExists: Bool(true),
			name:        id,
			ObjectProperties: &UserObjectProperties{
				Password:  &password,
				LoginName: &loginName,
			},
			ObjectParameters: &UserObjectParameters{
				EnableUnredactedQuerySyntaxError: Bool(true),
			},
			SessionParameters: &SessionParameters{
				Autocommit: Bool(true),
			},
			With: Bool(true),
			Tags: tags,
		}
		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.name, userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)

		user, err := client.Users.Show(ctx, &ShowUserOptions{
			Like: &Like{
				Pattern: &id.name,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(user))
		assert.Equal(t, id.name, user[0].Name)
	})

	t.Run("test no options", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.name, userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(id.name), userDetails.LoginName.Value)
		assert.Empty(t, userDetails.Password.Value)

		user, err := client.Users.Show(ctx, &ShowUserOptions{
			Like: &Like{
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
	ctx := context.Background()

	user, userCleanup := createUser(t, client)
	t.Cleanup(userCleanup)

	t.Run("when user exists", func(t *testing.T) {
		userDetails, err := client.Users.Describe(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.Name, userDetails.Name.Value)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		id := NewAccountObjectIdentifier("does_not_exist")
		_, err := client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})
}

func TestInt_UserDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("when user exists", func(t *testing.T) {
		user, _ := createUser(t, client)
		id := user.ID()
		err := client.Users.Drop(ctx, id)
		require.NoError(t, err)
		_, err = client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		id := NewAccountObjectIdentifier("does_not_exist")
		err := client.Users.Drop(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})
}
