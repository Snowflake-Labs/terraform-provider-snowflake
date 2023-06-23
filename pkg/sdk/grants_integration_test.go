package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_GrantPrivilegeToShare(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	t.Run("without options", func(t *testing.T) {
		err := client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, nil, shareTest.ID())
		require.Error(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		err := client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       databaseTest.ID(),
				},
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(grants))
		var shareGrant *Grant
		for _, grant := range grants {
			if grant.GranteeName.Name() == shareTest.ID().Name() {
				shareGrant = grant
				break
			}
		}
		assert.NotNil(t, shareGrant)
		assert.Equal(t, PrivilegeUsage, shareGrant.Privilege)
		assert.Equal(t, ObjectTypeDatabase, shareGrant.GrantedOn)
		assert.Equal(t, ObjectTypeShare, shareGrant.GrantedTo)
		assert.Equal(t, databaseTest.ID().Name(), shareGrant.Name.Name())
		err = client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
}

func TestInt_RevokePrivilegeToShare(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
		Database: databaseTest.ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Run("without options", func(t *testing.T) {
		err = client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, nil, shareTest.ID())
		require.Error(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		err = client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
}

func TestInt_ShowGrants(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, PrivilegeUsage, &GrantPrivilegeToShareOn{
		Database: databaseTest.ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err = client.Grants.RevokePrivilegeFromShare(ctx, PrivilegeUsage, &RevokePrivilegeFromShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})
	t.Run("without options", func(t *testing.T) {
		_, err := client.Grants.Show(ctx, nil)
		require.Error(t, err)
	})
	t.Run("with options", func(t *testing.T) {
		grants, err := client.Grants.Show(ctx, &ShowGrantOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       databaseTest.ID(),
				},
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(grants))
	})
}
