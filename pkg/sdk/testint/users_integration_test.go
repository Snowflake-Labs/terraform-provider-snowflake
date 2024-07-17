package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Users(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	randomPrefix := random.AlphaN(6)

	user, userCleanup := testClientHelper().User.CreateUserWithPrefix(t, randomPrefix+"_")
	t.Cleanup(userCleanup)

	user2, user2Cleanup := testClientHelper().User.CreateUserWithPrefix(t, randomPrefix)
	t.Cleanup(user2Cleanup)

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	t.Run("create: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToUpper(random.AlphaN(6))
		tagValue := random.String()
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := random.Password()
		loginName := random.String()

		opts := &sdk.CreateUserOptions{
			OrReplace: sdk.Bool(true),
			ObjectProperties: &sdk.UserObjectProperties{
				Password:    &password,
				LoginName:   &loginName,
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
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
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)
		assert.Equal(t, defaultRole, userDetails.DefaultRole.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), user.Name)
	})

	t.Run("create: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		tagValue := random.String()
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := random.Password()
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
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), user.Name)
	})

	t.Run("create: no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(id.Name()), userDetails.LoginName.Value)
		assert.Empty(t, userDetails.Password.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), user.Name)
	})

	t.Run("create: default role with hyphen", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToUpper(random.AlphaN(4) + "-" + random.AlphaN(4))

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		createdUser, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, defaultRole, createdUser.DefaultRole)
	})

	t.Run("create: default role in lowercase", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToLower(random.AlphaN(6))

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		createdUser, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, defaultRole, createdUser.DefaultRole)
	})

	t.Run("create: other params with hyphen and mixed cases", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		randomWithHyphenAndMixedCase := strings.ToUpper(random.AlphaN(4)) + "-" + strings.ToLower(random.AlphaN(4))
		var namespaceId sdk.ObjectIdentifier = sdk.NewDatabaseObjectIdentifier(randomWithHyphenAndMixedCase, randomWithHyphenAndMixedCase)

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				LoginName:        sdk.String(randomWithHyphenAndMixedCase),
				DisplayName:      sdk.String(randomWithHyphenAndMixedCase),
				FirstName:        sdk.String(randomWithHyphenAndMixedCase),
				MiddleName:       sdk.String(randomWithHyphenAndMixedCase),
				LastName:         sdk.String(randomWithHyphenAndMixedCase),
				DefaultWarehouse: sdk.Pointer(sdk.NewAccountObjectIdentifier(randomWithHyphenAndMixedCase)),
				DefaultNamespace: sdk.Pointer(namespaceId),
				DefaultRole:      sdk.Pointer(sdk.NewAccountObjectIdentifier(randomWithHyphenAndMixedCase)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		createdUser, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		// login name is always case-insensitive
		assert.Equal(t, strings.ToUpper(randomWithHyphenAndMixedCase), createdUser.LoginName)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.DisplayName)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.FirstName)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.LastName)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.DefaultWarehouse)
		assert.Equal(t, randomWithHyphenAndMixedCase+"."+randomWithHyphenAndMixedCase, createdUser.DefaultNamespace)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.DefaultRole)

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.MiddleName.Value)
		// login name is always case-insensitive
		assert.Equal(t, strings.ToUpper(randomWithHyphenAndMixedCase), userDetails.LoginName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DisplayName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.FirstName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.LastName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DefaultWarehouse.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase+"."+randomWithHyphenAndMixedCase, userDetails.DefaultNamespace.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DefaultRole.Value)
	})

	// TODO: add tests for alter
	// TODO: add tests for parameters

	t.Run("describe: when user exists", func(t *testing.T) {
		userDetails, err := client.Users.Describe(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.Name, userDetails.Name.Value)
	})

	t.Run("describe: when user does not exist", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		_, err := client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when user exists", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		id := user.ID()
		err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{})
		require.NoError(t, err)
		_, err = client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when user does not exist", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{})
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show: with like options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(user.Name),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Equal(t, 1, len(users))
	})

	t.Run("show: with starts with options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			StartsWith: sdk.String(randomPrefix),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Contains(t, users, *user2)
		assert.Equal(t, 2, len(users))
	})

	t.Run("show: with starts with, limit and from options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Limit:      sdk.Int(10),
			From:       sdk.String(randomPrefix + "_"),
			StartsWith: sdk.String(randomPrefix),
		}

		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Equal(t, 1, len(users))
	})

	t.Run("show: search for a non-existent user", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(NonExistingAccountObjectIdentifier.Name()),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(users))
	})

	t.Run("show: limit the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Limit: sdk.Int(1),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})
}
