package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestInt_FailoverGroupsCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)

	t.Run("test complete", func(t *testing.T) {
		name := randomStringRange(t, 8, 28)
		id := NewAccountObjectIdentifier(name)
		objectTypes := []ObjectType{
			ObjectTypeShare,
			ObjectTypeDatabase,
		}
		allowedAccounts := []AccountIdentifier{
			secondaryAccountIdentifier(t),
		}
		replicationSchedule := "10 MINUTE"
		err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &FailoverGroupCreateOptions{
			IfNotExists: Bool(true),
			AllowedDatabases: []AccountObjectIdentifier{
				databaseTest.ID(),
			},
			AllowedShares: []AccountObjectIdentifier{
				shareTest.ID(),
			},
			IgnoreEditionCheck:  Bool(true),
			ReplicationSchedule: String(replicationSchedule),
		})
		require.NoError(t, err)
		failoverGroups, err := client.FailoverGroups.Show(ctx, nil)
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, id, nil)
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)
		assert.GreaterOrEqual(t, len(failoverGroups), 1)
		var failoverGroup *FailoverGroup
		for _, fg := range failoverGroups {
			if fg.Name == name {
				failoverGroup = fg
				break
			}
		}
		assert.NotNil(t, failoverGroup)
		assert.Equal(t, name, failoverGroup.Name)
		slices.Sort(objectTypes)
		slices.Sort(failoverGroup.ObjectTypes)
		assert.Equal(t, objectTypes, failoverGroup.ObjectTypes)
		assert.Equal(t, 0, len(failoverGroup.AllowedIntegrationTypes))
		// this is length 2 because it automatically adds the current account to allowed accounts list
		assert.Equal(t, 2, len(failoverGroup.AllowedAccounts))
		for _, allowedAccount := range allowedAccounts {
			assert.Contains(t, failoverGroup.AllowedAccounts, allowedAccount)
		}
		assert.Equal(t, replicationSchedule, failoverGroup.ReplicationSchedule)

		fgDBS, err := client.FailoverGroups.ShowDatabases(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 1, len(fgDBS))
		assert.Equal(t, databaseTest.ID().Name(), fgDBS[0].Name())

		fgShares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 1, len(fgShares))
		assert.Equal(t, shareTest.ID().Name(), fgShares[0].Name())
	})
}

/*
func TestFailoverGroups(t *testing.T) {
	// first create a new failover group
	client := testClient(t)
	ctx := context.Background()
	id := NewAccountObjectIdentifier("test_failover_group")
	objectTypes := []ObjectType{ObjectTypeDatabase,ObjectTypeIntegration}
	allowedAccounts := []AccountIdentifier{NewAccountIdentifier("sfdeverel","cloud_engineering4")}
	client.FailoverGroups.Create(ctx,id,objectTypes,allowedAccounts,&FailoverGroupCreateOptions{
		Comment: "test failover group",
	})
}
*/
