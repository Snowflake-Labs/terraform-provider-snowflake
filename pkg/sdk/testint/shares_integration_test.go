package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SharesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)

	_, shareCleanup2 := createShare(t, client)
	t.Cleanup(shareCleanup2)

	t.Run("without show options", func(t *testing.T) {
		shares, err := client.Shares.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(shares))
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		assert.Contains(t, shares, *shareTest)
	})

	t.Run("when searching a non-existent share", func(t *testing.T) {
		showOptions := &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(shares))
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowShareOptions{
			Limit: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
	})
}

func TestInt_SharesCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test complete", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()
		err := client.Shares.Create(ctx, id, &sdk.CreateShareOptions{
			OrReplace: sdk.Bool(true),
			Comment:   sdk.String("test comment"),
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(id.Name()),
			},
			Limit: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		assert.Equal(t, id.Name(), shares[0].Name.Name())
		assert.Equal(t, "test comment", shares[0].Comment)

		t.Cleanup(func() {
			err := client.Shares.Drop(ctx, id)
			require.NoError(t, err)
		})
	})

	t.Run("test no options", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()
		err := client.Shares.Create(ctx, id, &sdk.CreateShareOptions{
			OrReplace: sdk.Bool(true),
			Comment:   sdk.String("test comment"),
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(shares), 1)

		t.Cleanup(func() {
			err := client.Shares.Drop(ctx, id)
			require.NoError(t, err)
		})
	})
}

func TestInt_SharesDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when share exists", func(t *testing.T) {
		shareTest, _ := createShare(t, client)
		err := client.Shares.Drop(ctx, shareTest.ID())
		require.NoError(t, err)
	})

	t.Run("when share does not exist", func(t *testing.T) {
		err := client.Shares.Drop(ctx, sdk.NewAccountObjectIdentifier("does_not_exist"))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_SharesAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("add and remove accounts", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, client)
		t.Cleanup(shareCleanup)
		err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
				Database: testDb(t).ID(),
			}, shareTest.ID())
		})
		require.NoError(t, err)
		secondaryClient := testSecondaryClient(t)
		accountsToAdd := []sdk.AccountIdentifier{
			getAccountIdentifier(t, secondaryClient),
		}
		// first add the account.
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Add: &sdk.ShareAdd{
				Accounts:          accountsToAdd,
				ShareRestrictions: sdk.Bool(false),
			},
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		share := shares[0]
		assert.Equal(t, accountsToAdd, share.To)

		// now remove the account that was added.
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Remove: &sdk.ShareRemove{
				Accounts: accountsToAdd,
			},
		})
		require.NoError(t, err)
		shares, err = client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		share = shares[0]
		assert.Equal(t, 0, len(share.To))
	})

	t.Run("set accounts", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, client)
		t.Cleanup(shareCleanup)
		err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
				Database: testDb(t).ID(),
			}, shareTest.ID())
		})
		require.NoError(t, err)
		secondaryClient := testSecondaryClient(t)
		accountsToSet := []sdk.AccountIdentifier{
			getAccountIdentifier(t, secondaryClient),
		}
		// first add the account.
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Set: &sdk.ShareSet{
				Accounts: accountsToSet,
			},
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		share := shares[0]
		assert.Equal(t, accountsToSet, share.To)
	})

	t.Run("set and unset comment", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, client)
		t.Cleanup(shareCleanup)
		err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
				Database: testDb(t).ID(),
			}, shareTest.ID())
			require.NoError(t, err)
		})

		comment := random.Comment()
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Set: &sdk.ShareSet{
				Comment: sdk.String(comment),
			},
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		share := shares[0]
		assert.Equal(t, comment, share.Comment)

		// reset comment
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Unset: &sdk.ShareUnset{
				Comment: sdk.Bool(true),
			},
		})
		require.NoError(t, err)
		shares, err = client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		share = shares[0]
		assert.Equal(t, "", share.Comment)
	})

	t.Run("set and unset tags", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, client)
		t.Cleanup(shareCleanup)
		err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
				Database: testDb(t).ID(),
			}, shareTest.ID())
			require.NoError(t, err)
		})

		schemaTest, schemaCleanup := createSchema(t, client, testDb(t))
		t.Cleanup(schemaCleanup)
		tagTest, tagCleanup := createTag(t, client, testDb(t), schemaTest)
		t.Cleanup(tagCleanup)
		tagTest2, tagCleanup2 := createTag(t, client, testDb(t), schemaTest)
		t.Cleanup(tagCleanup2)
		tagAssociations := []sdk.TagAssociation{
			{
				Name:  tagTest.ID(),
				Value: random.String(),
			},
			{
				Name:  tagTest2.ID(),
				Value: random.String(),
			},
		}
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Set: &sdk.ShareSet{
				Tag: tagAssociations,
			},
		})
		require.NoError(t, err)
		tagValue, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), shareTest.ID(), sdk.ObjectTypeShare)
		require.NoError(t, err)
		assert.Equal(t, tagAssociations[0].Value, tagValue)
		tagValue, err = client.SystemFunctions.GetTag(ctx, tagTest2.ID(), shareTest.ID(), sdk.ObjectTypeShare)
		require.NoError(t, err)
		assert.Equal(t, tagAssociations[1].Value, tagValue)

		// unset tags
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Unset: &sdk.ShareUnset{
				Tag: []sdk.ObjectIdentifier{
					tagTest.ID(),
				},
			},
		})
		require.NoError(t, err)
		_, err = client.SystemFunctions.GetTag(ctx, tagTest.ID(), shareTest.ID(), sdk.ObjectTypeShare)
		require.Error(t, err)
		tagValue, err = client.SystemFunctions.GetTag(ctx, tagTest2.ID(), shareTest.ID(), sdk.ObjectTypeShare)
		require.NoError(t, err)
		assert.Equal(t, tagAssociations[1].Value, tagValue)
	})
}

func TestInt_ShareDescribeProvider(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("describe share", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, client)
		t.Cleanup(shareCleanup)

		err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
				Database: testDb(t).ID(),
			}, shareTest.ID())
			require.NoError(t, err)
		})

		t.Run("describe share by name", func(t *testing.T) {
			shareDetails, err := client.Shares.DescribeProvider(ctx, shareTest.ID())
			require.NoError(t, err)
			assert.Equal(t, 1, len(shareDetails.SharedObjects))
			sharedObject := shareDetails.SharedObjects[0]
			assert.Equal(t, sdk.ObjectTypeDatabase, sharedObject.Kind)
			assert.Equal(t, testDb(t).ID(), sharedObject.Name)
		})
	})
}

func TestInt_ShareDescribeConsumer(t *testing.T) {
	consumerClient := testSecondaryClient(t)
	ctx := testContext(t)
	providerClient := testClient(t)

	t.Run("describe share", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, providerClient)
		t.Cleanup(shareCleanup)

		err := providerClient.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
			Database: testDb(t).ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = providerClient.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
				Database: testDb(t).ID(),
			}, shareTest.ID())
			require.NoError(t, err)
		})

		// add consumer account to share.
		err = providerClient.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			Add: &sdk.ShareAdd{
				Accounts: []sdk.AccountIdentifier{
					getAccountIdentifier(t, consumerClient),
				},
			},
		})
		require.NoError(t, err)
		t.Run("describe consume share", func(t *testing.T) {
			shareDetails, err := consumerClient.Shares.DescribeConsumer(ctx, shareTest.ExternalID())
			require.NoError(t, err)
			assert.Equal(t, 1, len(shareDetails.SharedObjects))
			sharedObject := shareDetails.SharedObjects[0]
			assert.Equal(t, sdk.ObjectTypeDatabase, sharedObject.Kind)
			assert.Equal(t, sdk.NewAccountObjectIdentifier("<DB>"), sharedObject.Name)
		})
	})
}
