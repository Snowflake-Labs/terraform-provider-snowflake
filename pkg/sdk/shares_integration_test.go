package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SharesShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
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
		showOptions := &ShareShowOptions{
			Like: &Like{
				Pattern: String(shareTest.Name.Name()),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		assert.Contains(t, shares, shareTest)
	})

	t.Run("when searching a non-existent share", func(t *testing.T) {
		showOptions := &ShareShowOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(shares))
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &ShareShowOptions{
			Limit: &LimitPagination{
				Rows: Int(1),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
	})
}

func TestInt_SharesCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("test complete", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)
		err := client.Shares.Create(ctx, id, &ShareCreateOptions{
			OrReplace: Bool(true),
			Comment:   String("test comment"),
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &ShareShowOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			Limit: &LimitPagination{
				Rows: Int(1),
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
		id := randomAccountObjectIdentifier(t)
		err := client.Shares.Create(ctx, id, &ShareCreateOptions{
			OrReplace: Bool(true),
			Comment:   String("test comment"),
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
	ctx := context.Background()

	t.Run("when share exists", func(t *testing.T) {
		shareTest, _ := createShare(t, client)
		err := client.Shares.Drop(ctx, shareTest.ID())
		require.NoError(t, err)
	})

	t.Run("when share does not exist", func(t *testing.T) {
		err := client.Shares.Drop(ctx, NewAccountObjectIdentifier("does_not_exist"))
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_Alter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	t.Run("add and remove accounts", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, client)
		t.Cleanup(shareCleanup)
		client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		t.Cleanup(func() {
			client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
				Database: databaseTest.ID(),
			}, shareTest.ID())
		})

		accountsToAdd := []AccountIdentifier{
			secondaryAccountIdentifier(t),
		}
		// first add the account.
		err := client.Shares.Alter(ctx, shareTest.ID(), &ShareAlterOptions{
			IfExists: Bool(true),
			Add: &ShareAdd{
				Accounts:          accountsToAdd,
				ShareRestrictions: Bool(false),
			},
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &ShareShowOptions{
			Like: &Like{
				Pattern: String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		share := shares[0]
		assert.Equal(t, accountsToAdd, share.To)

		// now remove the account that was added.
		err = client.Shares.Alter(ctx, shareTest.ID(), &ShareAlterOptions{
			IfExists: Bool(true),
			Remove: &ShareRemove{
				Accounts: accountsToAdd,
			},
		})
		require.NoError(t, err)
		shares, err = client.Shares.Show(ctx, &ShareShowOptions{
			Like: &Like{
				Pattern: String(shareTest.Name.Name()),
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
		client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		t.Cleanup(func() {
			client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
				Database: databaseTest.ID(),
			}, shareTest.ID())
		})

		accountsToSet := []AccountIdentifier{
			secondaryAccountIdentifier(t),
		}
		// first add the account.
		err := client.Shares.Alter(ctx, shareTest.ID(), &ShareAlterOptions{
			IfExists: Bool(true),
			Set: &ShareSet{
				Accounts: accountsToSet,
			},
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &ShareShowOptions{
			Like: &Like{
				Pattern: String(shareTest.Name.Name()),
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
		client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		t.Cleanup(func() {
			client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
				Database: databaseTest.ID(),
			}, shareTest.ID())
		})

		comment := randomComment(t)
		err := client.Shares.Alter(ctx, shareTest.ID(), &ShareAlterOptions{
			IfExists: Bool(true),
			Set: &ShareSet{
				Comment: String(comment),
			},
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &ShareShowOptions{
			Like: &Like{
				Pattern: String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))
		share := shares[0]
		assert.Equal(t, comment, share.Comment)

		// reset comment
		err = client.Shares.Alter(ctx, shareTest.ID(), &ShareAlterOptions{
			IfExists: Bool(true),
			Unset: &ShareUnset{
				Comment: Bool(true),
			},
		})
		require.NoError(t, err)
		shares, err = client.Shares.Show(ctx, &ShareShowOptions{
			Like: &Like{
				Pattern: String(shareTest.Name.Name()),
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
		client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		t.Cleanup(func() {
			client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
				Database: databaseTest.ID(),
			}, shareTest.ID())
		})

		schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
		t.Cleanup(schemaCleanup)
		tagTest, tagCleanup := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tagCleanup)
		tagTest2, tagCleanup2 := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tagCleanup2)
		tagAssociations := []TagAssociation{
			{
				Name:  tagTest.ID(),
				Value: randomString(t),
			},
			{
				Name:  tagTest2.ID(),
				Value: randomString(t),
			},
		}
		err := client.Shares.Alter(ctx, shareTest.ID(), &ShareAlterOptions{
			IfExists: Bool(true),
			Set: &ShareSet{
				Tag: tagAssociations,
			},
		})
		require.NoError(t, err)
		tagValue, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), shareTest.ID(), ObjectTypeShare)
		require.NoError(t, err)
		assert.Equal(t, tagAssociations[0].Value, tagValue)
		tagValue, err = client.SystemFunctions.GetTag(ctx, tagTest2.ID(), shareTest.ID(), ObjectTypeShare)
		require.NoError(t, err)
		assert.Equal(t, tagAssociations[1].Value, tagValue)

		// unset tags
		err = client.Shares.Alter(ctx, shareTest.ID(), &ShareAlterOptions{
			IfExists: Bool(true),
			Unset: &ShareUnset{
				Tag: []ObjectIdentifier{
					tagTest.ID(),
				},
			},
		})
		require.NoError(t, err)
		_, err = client.SystemFunctions.GetTag(ctx, tagTest.ID(), shareTest.ID(), ObjectTypeShare)
		require.Error(t, err)
		tagValue, err = client.SystemFunctions.GetTag(ctx, tagTest2.ID(), shareTest.ID(), ObjectTypeShare)
		require.NoError(t, err)
		assert.Equal(t, tagAssociations[1].Value, tagValue)
	})
}

func TestInt_ShareDescribeProvider(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("describe share", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, client)
		t.Cleanup(shareCleanup)

		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)

		client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		t.Cleanup(func() {
			client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
				Database: databaseTest.ID(),
			}, shareTest.ID())
		})

		t.Run("describe share by name", func(t *testing.T) {
			shareDetails, err := client.Shares.DescribeProvider(ctx, shareTest.ID())
			require.NoError(t, err)
			assert.Equal(t, 1, len(shareDetails.SharedObjects))
			sharedObject := shareDetails.SharedObjects[0]
			assert.Equal(t, ObjectTypeDatabase, sharedObject.Kind)
			assert.Equal(t, databaseTest.ID(), sharedObject.Name)
		})
	})
}

func TestInt_ShareDescribeConsumer(t *testing.T) {
	consumerClient, err := testClientFromProfile(t, "secondary_test_account")
	if err != nil {
		t.Skip("Skipping test due to missing profile credentials")
	}
	ctx := context.Background()
	providerClient := testClient(t)

	t.Run("describe share", func(t *testing.T) {
		shareTest, shareCleanup := createShare(t, providerClient)
		t.Cleanup(shareCleanup)

		databaseTest, databaseCleanup := createDatabase(t, providerClient)
		t.Cleanup(databaseCleanup)

		providerClient.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		t.Cleanup(func() {
			providerClient.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
				Database: databaseTest.ID(),
			}, shareTest.ID())
		})

		// add consumer account to share.
		err := providerClient.Shares.Alter(ctx, shareTest.ID(), &ShareAlterOptions{
			Add: &ShareAdd{
				Accounts: []AccountIdentifier{
					secondaryAccountIdentifier(t),
				},
			},
		})
		require.NoError(t, err)
		t.Run("describe consume share", func(t *testing.T) {
			shareDetails, err := consumerClient.Shares.DescribeConsumer(ctx, shareTest.ExternalID())
			require.NoError(t, err)
			assert.Equal(t, 1, len(shareDetails.SharedObjects))
			sharedObject := shareDetails.SharedObjects[0]
			assert.Equal(t, ObjectTypeDatabase, sharedObject.Kind)
			assert.Equal(t, NewAccountObjectIdentifier("<DB>"), sharedObject.Name)
		})
	})
}
