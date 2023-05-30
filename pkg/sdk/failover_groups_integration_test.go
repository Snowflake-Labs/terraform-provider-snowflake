package sdk

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestInt_FailoverGroupsCreate(t *testing.T) {
	if os.Getenv("SNOWFLAKE_TEST_BUSINESS_CRITICAL_FEATURES") != "1" {
		t.Skip("Skipping TestInt_FailoverGroupsCreate")
	}
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)

	t.Run("test complete", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)
		objectTypes := []PluralObjectType{
			PluralObjectTypeShares,
			PluralObjectTypeDatabases,
		}
		allowedAccounts := []AccountIdentifier{
			getSecondaryAccountIdentifier(t),
		}
		replicationSchedule := "10 MINUTE"
		err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &CreateFailoverGroupOptions{
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
		failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, id, nil)
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)
		assert.Equal(t, id.Name(), failoverGroup.Name)
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

	t.Run("test with allowed integration types", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)
		objectTypes := []PluralObjectType{
			PluralObjectTypeIntegrations,
		}
		allowedAccounts := []AccountIdentifier{
			getSecondaryAccountIdentifier(t),
		}
		allowedIntegrationTypes := []IntegrationType{
			IntegrationTypeAPIIntegrations,
			IntegrationTypeNotificationIntegrations,
		}
		err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &CreateFailoverGroupOptions{
			AllowedIntegrationTypes: allowedIntegrationTypes,
		})
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, id, nil)
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)
		failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), failoverGroup.Name)
		slices.Sort(failoverGroup.AllowedIntegrationTypes)
		slices.Sort(allowedIntegrationTypes)
		assert.Equal(t, allowedIntegrationTypes, failoverGroup.AllowedIntegrationTypes)
	})
}

func TestInt_CreateSecondaryReplicationGroup(t *testing.T) {
	if os.Getenv("SNOWFLAKE_TEST_BUSINESS_CRITICAL_FEATURES") != "1" {
		t.Skip("Skipping TestInt_FailoverGroupsCreate")
	}
	client := testClient(t)
	ctx := context.Background()
	primaryAccountID := getAccountIdentifier(t, client)
	secondaryClient := testSecondaryClient(t)
	secondaryClientID := getAccountIdentifier(t, secondaryClient)

	// create a temp share
	shareTest, cleanupDatabase := createShare(t, client)
	t.Cleanup(cleanupDatabase)

	// create a failover group in primary account and share with target account
	id := randomAccountObjectIdentifier(t)

	opts := &CreateFailoverGroupOptions{
		AllowedShares: []AccountObjectIdentifier{
			shareTest.ID(),
		},
		ReplicationSchedule: String("10 MINUTE"),
	}
	allowedAccounts := []AccountIdentifier{
		primaryAccountID,
		secondaryClientID,
	}
	objectTypes := []PluralObjectType{
		PluralObjectTypeShares,
	}
	err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, opts)
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)

	// there is a delay between creating a failover group and it being available for replication
	time.Sleep(1 * time.Second)

	// create a replica of failover group in target account
	err = secondaryClient.FailoverGroups.CreateSecondaryReplicationGroup(ctx, failoverGroup.ID(), failoverGroup.ExternalID(), &CreateSecondaryReplicationGroupOptions{
		IfNotExists: Bool(true),
	})
	require.NoError(t, err)

	// cleanup failover groups with retry (in case of replication delay)
	cleanupFailoverGroups := func() {
		err := retry.Do(
			func() error {
				return client.FailoverGroups.Drop(ctx, failoverGroup.ID(), nil)
			},
			retry.OnRetry(func(n uint, err error) {
				log.Printf("[DEBUG] Retrying client.FailoverGroups.Drop(): #%d", n+1)
			}),
			retry.Delay(1*time.Second),
			retry.Attempts(3),
		)
		require.NoError(t, err)
		err = retry.Do(
			func() error {
				return secondaryClient.FailoverGroups.Drop(ctx, failoverGroup.ID(), nil)
			},
			retry.OnRetry(func(n uint, err error) {
				log.Printf("[DEBUG] Retrying client.FailoverGroups.Drop(): #%d", n+1)
			}),
			retry.Delay(1*time.Second),
			retry.Attempts(3),
		)
		require.NoError(t, err)
	}
	t.Cleanup(cleanupFailoverGroups)

	failoverGroups, err := secondaryClient.FailoverGroups.Show(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, len(failoverGroups))
}

func TestInt_FailoverGroupsAlterSource(t *testing.T) {
	if os.Getenv("SNOWFLAKE_TEST_BUSINESS_CRITICAL_FEATURES") != "1" {
		t.Skip("Skipping TestInt_FailoverGroupsCreate")
	}
	client := testClient(t)
	ctx := context.Background()
	t.Run("rename the failover group", func(t *testing.T) {
		failoverGroup, _ := createFailoverGroup(t, client)
		oldID := failoverGroup.ID()
		newID := randomAccountObjectIdentifier(t)
		opts := &AlterSourceFailoverGroupOptions{
			NewName: newID,
		}
		err := client.FailoverGroups.AlterSource(ctx, oldID, opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), failoverGroup.Name)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, newID, nil)
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)
	})

	t.Run("reset the list of specified object types enabled for replication and failover.", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)
		objectTypes := []PluralObjectType{
			PluralObjectTypeDatabases,
		}
		opts := &AlterSourceFailoverGroupOptions{
			Set: &FailoverGroupSet{
				ObjectTypes: objectTypes,
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, objectTypes, failoverGroup.ObjectTypes)
	})

	t.Run("set or update the replication schedule for automatic refresh of secondary failover groups.", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)
		replicationSchedule := "USING CRON 0 0 10-20 * TUE,THU UTC"
		opts := &AlterSourceFailoverGroupOptions{
			Set: &FailoverGroupSet{
				ReplicationSchedule: &replicationSchedule,
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, replicationSchedule, failoverGroup.ReplicationSchedule)
	})

	t.Run("add and remove database account object", func(t *testing.T) {
		databaseTest, cleanupDatabase := createDatabase(t, client)
		t.Cleanup(cleanupDatabase)
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)

		// first add databases to allowed object types
		opts := &AlterSourceFailoverGroupOptions{
			Set: &FailoverGroupSet{
				ObjectTypes: []PluralObjectType{
					PluralObjectTypeDatabases,
				},
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(failoverGroup.ObjectTypes))
		assert.Equal(t, PluralObjectTypeDatabases, failoverGroup.ObjectTypes[0])

		// now add database to allowed databases
		opts = &AlterSourceFailoverGroupOptions{
			Add: &FailoverGroupAdd{
				AllowedDatabases: []AccountObjectIdentifier{
					databaseTest.ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		allowedDBs, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(allowedDBs))
		assert.Equal(t, databaseTest.ID().Name(), allowedDBs[0].Name())

		// now remove database from allowed databases
		opts = &AlterSourceFailoverGroupOptions{
			Remove: &FailoverGroupRemove{
				AllowedDatabases: []AccountObjectIdentifier{
					databaseTest.ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		allowedDBs, err = client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 0, len(allowedDBs))
	})

	t.Run("add and remove share account object", func(t *testing.T) {
		shareTest, cleanupDatabase := createShare(t, client)
		t.Cleanup(cleanupDatabase)
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)

		// first add shares to allowed object types
		opts := &AlterSourceFailoverGroupOptions{
			Set: &FailoverGroupSet{
				ObjectTypes: []PluralObjectType{
					PluralObjectTypeShares,
				},
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(failoverGroup.ObjectTypes))
		assert.Equal(t, shareTest.ObjectType().Plural(), failoverGroup.ObjectTypes[0])
		// now add share to allowed shares
		opts = &AlterSourceFailoverGroupOptions{
			Add: &FailoverGroupAdd{
				AllowedShares: []AccountObjectIdentifier{
					shareTest.ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		allowedShares, err := client.FailoverGroups.ShowShares(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(allowedShares))
		assert.Equal(t, shareTest.ID().Name(), allowedShares[0].Name())

		// now remove share from allowed shares
		opts = &AlterSourceFailoverGroupOptions{
			Remove: &FailoverGroupRemove{
				AllowedShares: []AccountObjectIdentifier{
					shareTest.ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		allowedShares, err = client.FailoverGroups.ShowShares(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 0, len(allowedShares))
	})

	t.Run("add and remove security integration account object", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)
		// first add security integrations to allowed object types
		opts := &AlterSourceFailoverGroupOptions{
			Set: &FailoverGroupSet{
				ObjectTypes: []PluralObjectType{
					PluralObjectTypeIntegrations,
				},
				AllowedIntegrationTypes: []IntegrationType{
					IntegrationTypeAPIIntegrations,
					IntegrationTypeNotificationIntegrations,
				},
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 2, len(failoverGroup.AllowedIntegrationTypes))
		assert.Equal(t, IntegrationTypeAPIIntegrations, failoverGroup.AllowedIntegrationTypes[0])
		assert.Equal(t, IntegrationTypeNotificationIntegrations, failoverGroup.AllowedIntegrationTypes[1])
		assert.Equal(t, 1, len(failoverGroup.ObjectTypes))
		assert.Equal(t, PluralObjectTypeIntegrations, failoverGroup.ObjectTypes[0])

		// now remove security integration from allowed security integrations
		opts = &AlterSourceFailoverGroupOptions{
			Set: &FailoverGroupSet{
				ObjectTypes: []PluralObjectType{
					PluralObjectTypeIntegrations,
				},
				AllowedIntegrationTypes: []IntegrationType{
					IntegrationTypeAPIIntegrations,
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(failoverGroup.AllowedIntegrationTypes))
		assert.Equal(t, IntegrationTypeAPIIntegrations, failoverGroup.AllowedIntegrationTypes[0])
	})

	t.Run("add or remove target accounts enabled for replication and failover", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)

		secondaryAccountID := getSecondaryAccountIdentifier(t)
		// first add target account
		opts := &AlterSourceFailoverGroupOptions{
			Add: &FailoverGroupAdd{
				AllowedAccounts: []AccountIdentifier{
					secondaryAccountID,
				},
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 2, len(failoverGroup.AllowedAccounts))
		assert.Contains(t, failoverGroup.AllowedAccounts, secondaryAccountID)

		// now remove target accounts
		opts = &AlterSourceFailoverGroupOptions{
			Remove: &FailoverGroupRemove{
				AllowedAccounts: []AccountIdentifier{
					secondaryAccountID,
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(failoverGroup.AllowedAccounts))
		assert.Contains(t, failoverGroup.AllowedAccounts, getAccountIdentifier(t, client))
	})

	t.Run("move shares to another failover group", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)

		// add "SHARES" to object types of both failover groups
		opts := &AlterSourceFailoverGroupOptions{
			Set: &FailoverGroupSet{
				ObjectTypes: []PluralObjectType{
					PluralObjectTypeShares,
				},
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		failoverGroup2, cleanupFailoverGroup2 := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup2)

		err = client.FailoverGroups.AlterSource(ctx, failoverGroup2.ID(), opts)
		require.NoError(t, err)

		// create a temp share
		shareTest, cleanupShare := createShare(t, client)
		t.Cleanup(cleanupShare)

		// now add share to allowed shares of failover group 1
		opts = &AlterSourceFailoverGroupOptions{
			Add: &FailoverGroupAdd{
				AllowedShares: []AccountObjectIdentifier{
					shareTest.ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// now move share to failover group 2
		opts = &AlterSourceFailoverGroupOptions{
			Move: &FailoverGroupMove{
				Shares: []AccountObjectIdentifier{
					shareTest.ID(),
				},
				To: failoverGroup2.ID(),
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// verify that share is now in failover group 2
		shares, err := client.FailoverGroups.ShowShares(ctx, failoverGroup2.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(shares))

		// verify that share is not in failover group 1
		shares, err = client.FailoverGroups.ShowShares(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 0, len(shares))
	})

	t.Run("move database to another failover group", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)

		// add "DATABASES" to object types of both failover groups
		opts := &AlterSourceFailoverGroupOptions{
			Set: &FailoverGroupSet{
				ObjectTypes: []PluralObjectType{
					PluralObjectTypeDatabases,
				},
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		failoverGroup2, cleanupFailoverGroup2 := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup2)

		err = client.FailoverGroups.AlterSource(ctx, failoverGroup2.ID(), opts)
		require.NoError(t, err)

		// create a temp database
		databaseTest, cleanupDatabase := createDatabase(t, client)
		t.Cleanup(cleanupDatabase)

		// now add database to allowed databases of failover group 1
		opts = &AlterSourceFailoverGroupOptions{
			Add: &FailoverGroupAdd{
				AllowedDatabases: []AccountObjectIdentifier{
					databaseTest.ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// now move database to failover group 2
		opts = &AlterSourceFailoverGroupOptions{
			Move: &FailoverGroupMove{
				Databases: []AccountObjectIdentifier{
					databaseTest.ID(),
				},
				To: failoverGroup2.ID(),
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// verify that database is now in failover group 2
		databases, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroup2.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(databases))

		// verify that database is not in failover group 1
		databases, err = client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 0, len(databases))
	})
}

func TestInt_FailoverGroupsAlterTarget(t *testing.T) {
	if os.Getenv("SNOWFLAKE_TEST_BUSINESS_CRITICAL_FEATURES") != "1" {
		t.Skip("Skipping TestInt_FailoverGroupsCreate")
	}
	client := testClient(t)
	ctx := context.Background()
	primaryAccountID := getAccountIdentifier(t, client)
	secondaryClient := testSecondaryClient(t)
	secondaryClientID := getAccountIdentifier(t, secondaryClient)

	// create a temp database
	databaseTest, cleanupDatabase := createDatabase(t, client)
	t.Cleanup(cleanupDatabase)

	// create a failover group in primary account and share with target account
	id := randomAccountObjectIdentifier(t)

	opts := &CreateFailoverGroupOptions{
		AllowedDatabases: []AccountObjectIdentifier{
			databaseTest.ID(),
		},
		ReplicationSchedule: String("10 MINUTE"),
	}
	allowedAccounts := []AccountIdentifier{
		primaryAccountID,
		secondaryClientID,
	}
	objectTypes := []PluralObjectType{
		PluralObjectTypeDatabases,
	}
	err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, opts)
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)

	// there is a delay between creating a failover group and it being available for replication
	time.Sleep(1 * time.Second)

	// create a replica of failover group in target account
	err = secondaryClient.FailoverGroups.CreateSecondaryReplicationGroup(ctx, failoverGroup.ID(), failoverGroup.ExternalID(), &CreateSecondaryReplicationGroupOptions{
		IfNotExists: Bool(true),
	})
	require.NoError(t, err)

	// cleanup failover groups with retry (in case of replication delay)
	cleanupFailoverGroups := func() {
		err := retry.Do(
			func() error {
				return client.FailoverGroups.Drop(ctx, failoverGroup.ID(), nil)
			},
			retry.OnRetry(func(n uint, err error) {
				log.Printf("[DEBUG] Retrying client.FailoverGroups.Drop(): #%d", n+1)
			}),
			retry.Delay(1*time.Second),
			retry.Attempts(3),
		)
		require.NoError(t, err)
		err = retry.Do(
			func() error {
				return secondaryClient.FailoverGroups.Drop(ctx, failoverGroup.ID(), nil)
			},
			retry.OnRetry(func(n uint, err error) {
				log.Printf("[DEBUG] Retrying client.FailoverGroups.Drop(): #%d", n+1)
			}),
			retry.Delay(1*time.Second),
			retry.Attempts(3),
		)
		require.NoError(t, err)
	}
	t.Cleanup(cleanupFailoverGroups)

	failoverGroups, err := secondaryClient.FailoverGroups.Show(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, len(failoverGroups))

	t.Run("perform suspend and resume", func(t *testing.T) {
		// suspend target failover group
		opts := &AlterTargetFailoverGroupOptions{
			Suspend: Bool(true),
		}
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// verify that target failover group is suspended
		fg, err := secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, FailoverGroupSecondaryStateSuspended, fg.SecondaryState)

		// resume target failover group
		opts = &AlterTargetFailoverGroupOptions{
			Resume: Bool(true),
		}
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// verify that target failover group is resumed
		failoverGroup, err = secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, FailoverGroupSecondaryStateStarted, failoverGroup.SecondaryState)
	})

	t.Run("refresh target failover group", func(t *testing.T) {
		// refresh target failover group
		opts := &AlterTargetFailoverGroupOptions{
			Refresh: Bool(true),
		}
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
	})

	t.Run("promote secondary to primary", func(t *testing.T) {
		// promote secondary to primary
		opts := &AlterTargetFailoverGroupOptions{
			Primary: Bool(true),
		}
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// verify that target failover group is promoted
		failoverGroup, err = secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, true, failoverGroup.IsPrimary)
	})
}

func TestInt_FailoverGroupsDrop(t *testing.T) {
	if os.Getenv("SNOWFLAKE_TEST_BUSINESS_CRITICAL_FEATURES") != "1" {
		t.Skip("Skipping TestInt_FailoverGroupsCreate")
	}
	client := testClient(t)
	ctx := context.Background()
	t.Run("no options", func(t *testing.T) {
		failoverGroup, _ := createFailoverGroup(t, client)
		err := client.FailoverGroups.Drop(ctx, failoverGroup.ID(), nil)
		require.NoError(t, err)
	})

	t.Run("with IfExists", func(t *testing.T) {
		failoverGroup, _ := createFailoverGroup(t, client)
		opts := &DropFailoverGroupOptions{
			IfExists: Bool(true),
		}
		err := client.FailoverGroups.Drop(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
	})
}

func TestInt_FailoverGroupsShow(t *testing.T) {
	if os.Getenv("SNOWFLAKE_TEST_BUSINESS_CRITICAL_FEATURES") != "1" {
		t.Skip("Skipping TestInt_FailoverGroupsCreate")
	}
	client := testClient(t)
	ctx := context.Background()
	failoverGroupTest, failoverGroupCleanup := createFailoverGroup(t, client)
	t.Cleanup(failoverGroupCleanup)

	t.Run("without show options", func(t *testing.T) {
		failoverGroups, err := client.FailoverGroups.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(failoverGroups))
		assert.Contains(t, failoverGroups, failoverGroupTest)
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &ShowFailoverGroupOptions{
			InAccount: NewAccountIdentifierFromAccountLocator(client.accountLocator),
		}
		failoverGroups, err := client.FailoverGroups.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(failoverGroups))
		assert.Contains(t, failoverGroups, failoverGroupTest)
	})

	t.Run("when searching a non-existent failover group", func(t *testing.T) {
		_, err := client.FailoverGroups.ShowByID(ctx, NewAccountObjectIdentifier("does-not-exist"))
		require.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_FailoverGroupsShowDatabases(t *testing.T) {
	if os.Getenv("SNOWFLAKE_TEST_BUSINESS_CRITICAL_FEATURES") != "1" {
		t.Skip("Skipping TestInt_FailoverGroupsCreate")
	}
	client := testClient(t)
	ctx := context.Background()
	failoverGroupTest, failoverGroupCleanup := createFailoverGroup(t, client)
	t.Cleanup(failoverGroupCleanup)

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	opts := &AlterSourceFailoverGroupOptions{
		Set: &FailoverGroupSet{
			ObjectTypes: []PluralObjectType{
				PluralObjectTypeDatabases,
			},
		},
	}
	err := client.FailoverGroups.AlterSource(ctx, failoverGroupTest.ID(), opts)
	require.NoError(t, err)
	opts = &AlterSourceFailoverGroupOptions{
		Add: &FailoverGroupAdd{
			AllowedDatabases: []AccountObjectIdentifier{
				databaseTest.ID(),
			},
		},
	}
	err = client.FailoverGroups.AlterSource(ctx, failoverGroupTest.ID(), opts)
	require.NoError(t, err)
	databases, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroupTest.ID())
	require.NoError(t, err)
	assert.Equal(t, 1, len(databases))
	assert.Equal(t, databaseTest.ID(), databases[0])
}

func TestInt_FailoverGroupsShowShares(t *testing.T) {
	if _, ok := os.LookupEnv("SNOWFLAKE_TEST_BUSINESS_CRITICAL_FEATURES"); !ok {
		t.Skip("Skipping TestInt_FailoverGroupsCreate")
	}
	client := testClient(t)
	ctx := context.Background()
	failoverGroupTest, failoverGroupCleanup := createFailoverGroup(t, client)
	t.Cleanup(failoverGroupCleanup)

	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	opts := &AlterSourceFailoverGroupOptions{
		Set: &FailoverGroupSet{
			ObjectTypes: []PluralObjectType{
				PluralObjectTypeShares,
			},
		},
	}
	err := client.FailoverGroups.AlterSource(ctx, failoverGroupTest.ID(), opts)
	require.NoError(t, err)
	opts = &AlterSourceFailoverGroupOptions{
		Add: &FailoverGroupAdd{
			AllowedShares: []AccountObjectIdentifier{
				shareTest.ID(),
			},
		},
	}
	err = client.FailoverGroups.AlterSource(ctx, failoverGroupTest.ID(), opts)
	require.NoError(t, err)
	shares, err := client.FailoverGroups.ShowShares(ctx, failoverGroupTest.ID())
	require.NoError(t, err)
	assert.Equal(t, 1, len(shares))
	assert.Equal(t, shareTest.ID(), shares[0])
}
