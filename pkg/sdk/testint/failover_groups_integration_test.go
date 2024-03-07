package testint

import (
	"log"
	"slices"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_FailoverGroupsCreate(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)

	t.Run("test complete", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeShares,
			sdk.PluralObjectTypeDatabases,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			businessCriticalAccountId,
		}
		replicationSchedule := "10 MINUTE"
		err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &sdk.CreateFailoverGroupOptions{
			IfNotExists: sdk.Bool(true),
			AllowedDatabases: []sdk.AccountObjectIdentifier{
				testDb(t).ID(),
			},
			AllowedShares: []sdk.AccountObjectIdentifier{
				shareTest.ID(),
			},
			IgnoreEditionCheck:  sdk.Bool(true),
			ReplicationSchedule: sdk.String(replicationSchedule),
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
		assert.Equal(t, testDb(t).ID().Name(), fgDBS[0].Name())

		fgShares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 1, len(fgShares))
		assert.Equal(t, shareTest.ID().Name(), fgShares[0].Name())
	})

	t.Run("test with identifier containing a dot", func(t *testing.T) {
		shareId := sdk.NewAccountObjectIdentifier(random.AlphanumericN(6) + "." + random.AlphanumericN(6))

		shareWithDot, shareWithDotCleanup := createShareWithOptions(t, client, shareId, &sdk.CreateShareOptions{})
		t.Cleanup(shareWithDotCleanup)

		id := sdk.RandomAccountObjectIdentifier()
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeShares,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			businessCriticalAccountId,
		}
		err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &sdk.CreateFailoverGroupOptions{
			AllowedShares: []sdk.AccountObjectIdentifier{
				shareWithDot.ID(),
			},
			IgnoreEditionCheck: sdk.Bool(true),
		})
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, id, nil)
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)

		fgShares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 1, len(fgShares))
		assert.Equal(t, shareWithDot.ID().Name(), fgShares[0].Name())
	})

	t.Run("test with allowed integration types", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeIntegrations,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			businessCriticalAccountId,
		}
		allowedIntegrationTypes := []sdk.IntegrationType{
			sdk.IntegrationTypeAPIIntegrations,
			sdk.IntegrationTypeNotificationIntegrations,
		}
		err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &sdk.CreateFailoverGroupOptions{
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

func TestInt_Issue2544(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)

	t.Run("alter object types, replication schedule, and allowed integration types at the same time", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeIntegrations,
			sdk.PluralObjectTypeDatabases,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			businessCriticalAccountId,
		}
		allowedIntegrationTypes := []sdk.IntegrationType{
			sdk.IntegrationTypeAPIIntegrations,
			sdk.IntegrationTypeNotificationIntegrations,
		}
		replicationSchedule := "10 MINUTE"
		err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &sdk.CreateFailoverGroupOptions{
			AllowedDatabases: []sdk.AccountObjectIdentifier{
				testDb(t).ID(),
			},
			AllowedIntegrationTypes: allowedIntegrationTypes,
			ReplicationSchedule:     sdk.String(replicationSchedule),
		})
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, id, nil)
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)

		newObjectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeIntegrations,
		}
		newAllowedIntegrationTypes := []sdk.IntegrationType{
			sdk.IntegrationTypeAPIIntegrations,
		}
		newReplicationSchedule := "20 MINUTE"

		// does not work together:
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ObjectTypes:             newObjectTypes,
				AllowedIntegrationTypes: newAllowedIntegrationTypes,
				ReplicationSchedule:     sdk.String(newReplicationSchedule),
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, id, opts)
		require.Error(t, err)
		require.ErrorContains(t, err, "unexpected 'REPLICATION_SCHEDULE'")

		// works as two separate alters:
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ObjectTypes:             newObjectTypes,
				AllowedIntegrationTypes: newAllowedIntegrationTypes,
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, id, opts)
		require.NoError(t, err)

		opts = &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ReplicationSchedule: sdk.String(newReplicationSchedule),
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, id, opts)
		require.NoError(t, err)
	})
}

func TestInt_CreateSecondaryReplicationGroup(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	primaryAccountID := getAccountIdentifier(t, client)
	secondaryClient := testSecondaryClient(t)
	secondaryClientID := getAccountIdentifier(t, secondaryClient)

	// create a temp share
	shareTest, cleanupDatabase := createShare(t, client)
	t.Cleanup(cleanupDatabase)

	// create a failover group in primary account and share with target account
	id := sdk.RandomAccountObjectIdentifier()

	opts := &sdk.CreateFailoverGroupOptions{
		AllowedShares: []sdk.AccountObjectIdentifier{
			shareTest.ID(),
		},
		ReplicationSchedule: sdk.String("10 MINUTE"),
	}
	allowedAccounts := []sdk.AccountIdentifier{
		primaryAccountID,
		secondaryClientID,
	}
	objectTypes := []sdk.PluralObjectType{
		sdk.PluralObjectTypeShares,
	}
	err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, opts)
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)

	// there is a delay between creating a failover group and it being available for replication
	time.Sleep(1 * time.Second)

	// create a replica of failover group in target account
	err = secondaryClient.FailoverGroups.CreateSecondaryReplicationGroup(ctx, failoverGroup.ID(), failoverGroup.ExternalID(), &sdk.CreateSecondaryReplicationGroupOptions{
		IfNotExists: sdk.Bool(true),
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
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	t.Run("rename the failover group", func(t *testing.T) {
		failoverGroup, _ := createFailoverGroup(t, client)
		oldID := failoverGroup.ID()
		newID := sdk.RandomAccountObjectIdentifier()
		opts := &sdk.AlterSourceFailoverGroupOptions{
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
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeDatabases,
		}
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
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
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
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
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)

		// first add databases to allowed object types
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ObjectTypes: []sdk.PluralObjectType{
					sdk.PluralObjectTypeDatabases,
				},
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(failoverGroup.ObjectTypes))
		assert.Equal(t, sdk.PluralObjectTypeDatabases, failoverGroup.ObjectTypes[0])

		// now add database to allowed databases
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Add: &sdk.FailoverGroupAdd{
				AllowedDatabases: []sdk.AccountObjectIdentifier{
					testDb(t).ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		allowedDBs, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(allowedDBs))
		assert.Equal(t, testDb(t).ID().Name(), allowedDBs[0].Name())

		// now remove database from allowed databases
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Remove: &sdk.FailoverGroupRemove{
				AllowedDatabases: []sdk.AccountObjectIdentifier{
					testDb(t).ID(),
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
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ObjectTypes: []sdk.PluralObjectType{
					sdk.PluralObjectTypeShares,
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
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Add: &sdk.FailoverGroupAdd{
				AllowedShares: []sdk.AccountObjectIdentifier{
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
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Remove: &sdk.FailoverGroupRemove{
				AllowedShares: []sdk.AccountObjectIdentifier{
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
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ObjectTypes: []sdk.PluralObjectType{
					sdk.PluralObjectTypeIntegrations,
				},
				AllowedIntegrationTypes: []sdk.IntegrationType{
					sdk.IntegrationTypeAPIIntegrations,
					sdk.IntegrationTypeNotificationIntegrations,
				},
			},
		}
		err := client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 2, len(failoverGroup.AllowedIntegrationTypes))
		assert.Equal(t, sdk.IntegrationTypeAPIIntegrations, failoverGroup.AllowedIntegrationTypes[0])
		assert.Equal(t, sdk.IntegrationTypeNotificationIntegrations, failoverGroup.AllowedIntegrationTypes[1])
		assert.Equal(t, 1, len(failoverGroup.ObjectTypes))
		assert.Equal(t, sdk.PluralObjectTypeIntegrations, failoverGroup.ObjectTypes[0])

		// now remove security integration from allowed security integrations
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ObjectTypes: []sdk.PluralObjectType{
					sdk.PluralObjectTypeIntegrations,
				},
				AllowedIntegrationTypes: []sdk.IntegrationType{
					sdk.IntegrationTypeAPIIntegrations,
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(failoverGroup.AllowedIntegrationTypes))
		assert.Equal(t, sdk.IntegrationTypeAPIIntegrations, failoverGroup.AllowedIntegrationTypes[0])
	})

	t.Run("add or remove target accounts enabled for replication and failover", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := createFailoverGroup(t, client)
		t.Cleanup(cleanupFailoverGroup)

		secondaryAccountID := getAccountIdentifier(t, testSecondaryClient(t))
		// first add target account
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Add: &sdk.FailoverGroupAdd{
				AllowedAccounts: []sdk.AccountIdentifier{
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
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Remove: &sdk.FailoverGroupRemove{
				AllowedAccounts: []sdk.AccountIdentifier{
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
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ObjectTypes: []sdk.PluralObjectType{
					sdk.PluralObjectTypeShares,
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
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Add: &sdk.FailoverGroupAdd{
				AllowedShares: []sdk.AccountObjectIdentifier{
					shareTest.ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// now move share to failover group 2
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Move: &sdk.FailoverGroupMove{
				Shares: []sdk.AccountObjectIdentifier{
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
		opts := &sdk.AlterSourceFailoverGroupOptions{
			Set: &sdk.FailoverGroupSet{
				ObjectTypes: []sdk.PluralObjectType{
					sdk.PluralObjectTypeDatabases,
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
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Add: &sdk.FailoverGroupAdd{
				AllowedDatabases: []sdk.AccountObjectIdentifier{
					databaseTest.ID(),
				},
			},
		}
		err = client.FailoverGroups.AlterSource(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// now move database to failover group 2
		opts = &sdk.AlterSourceFailoverGroupOptions{
			Move: &sdk.FailoverGroupMove{
				Databases: []sdk.AccountObjectIdentifier{
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
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	primaryAccountID := getAccountIdentifier(t, client)
	secondaryClient := testSecondaryClient(t)
	secondaryClientID := getAccountIdentifier(t, secondaryClient)

	// create a temp database
	databaseTest, cleanupDatabase := createDatabase(t, client)
	t.Cleanup(cleanupDatabase)

	// create a failover group in primary account and share with target account
	id := sdk.RandomAccountObjectIdentifier()

	opts := &sdk.CreateFailoverGroupOptions{
		AllowedDatabases: []sdk.AccountObjectIdentifier{
			databaseTest.ID(),
		},
		ReplicationSchedule: sdk.String("10 MINUTE"),
	}
	allowedAccounts := []sdk.AccountIdentifier{
		primaryAccountID,
		secondaryClientID,
	}
	objectTypes := []sdk.PluralObjectType{
		sdk.PluralObjectTypeDatabases,
	}
	err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, opts)
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)

	// there is a delay between creating a failover group and it being available for replication
	time.Sleep(1 * time.Second)

	// create a replica of failover group in target account
	err = secondaryClient.FailoverGroups.CreateSecondaryReplicationGroup(ctx, failoverGroup.ID(), failoverGroup.ExternalID(), &sdk.CreateSecondaryReplicationGroupOptions{
		IfNotExists: sdk.Bool(true),
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
		opts := &sdk.AlterTargetFailoverGroupOptions{
			Suspend: sdk.Bool(true),
		}
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// verify that target failover group is suspended
		fg, err := secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.FailoverGroupSecondaryStateSuspended, fg.SecondaryState)

		// resume target failover group
		opts = &sdk.AlterTargetFailoverGroupOptions{
			Resume: sdk.Bool(true),
		}
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)

		// verify that target failover group is resumed
		failoverGroup, err = secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.FailoverGroupSecondaryStateStarted, failoverGroup.SecondaryState)
	})

	t.Run("refresh target failover group", func(t *testing.T) {
		// refresh target failover group
		opts := &sdk.AlterTargetFailoverGroupOptions{
			Refresh: sdk.Bool(true),
		}
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
	})

	t.Run("promote secondary to primary", func(t *testing.T) {
		// promote secondary to primary
		opts := &sdk.AlterTargetFailoverGroupOptions{
			Primary: sdk.Bool(true),
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
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	t.Run("no options", func(t *testing.T) {
		failoverGroup, _ := createFailoverGroup(t, client)
		err := client.FailoverGroups.Drop(ctx, failoverGroup.ID(), nil)
		require.NoError(t, err)
	})

	t.Run("with IfExists", func(t *testing.T) {
		failoverGroup, _ := createFailoverGroup(t, client)
		opts := &sdk.DropFailoverGroupOptions{
			IfExists: sdk.Bool(true),
		}
		err := client.FailoverGroups.Drop(ctx, failoverGroup.ID(), opts)
		require.NoError(t, err)
	})
}

func TestInt_FailoverGroupsShow(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := createFailoverGroup(t, client)
	t.Cleanup(failoverGroupCleanup)

	t.Run("without show options", func(t *testing.T) {
		failoverGroups, err := client.FailoverGroups.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(failoverGroups))
		assert.Contains(t, failoverGroups, failoverGroupTest)
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowFailoverGroupOptions{
			InAccount: sdk.NewAccountIdentifierFromAccountLocator(client.GetAccountLocator()),
		}
		failoverGroups, err := client.FailoverGroups.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(failoverGroups))
		assert.Contains(t, failoverGroups, failoverGroupTest)
	})

	t.Run("when searching a non-existent failover group", func(t *testing.T) {
		_, err := client.FailoverGroups.ShowByID(ctx, sdk.NewAccountObjectIdentifier("does-not-exist"))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_FailoverGroupsShowDatabases(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := createFailoverGroup(t, client)
	t.Cleanup(failoverGroupCleanup)

	opts := &sdk.AlterSourceFailoverGroupOptions{
		Set: &sdk.FailoverGroupSet{
			ObjectTypes: []sdk.PluralObjectType{
				sdk.PluralObjectTypeDatabases,
			},
		},
	}
	err := client.FailoverGroups.AlterSource(ctx, failoverGroupTest.ID(), opts)
	require.NoError(t, err)
	opts = &sdk.AlterSourceFailoverGroupOptions{
		Add: &sdk.FailoverGroupAdd{
			AllowedDatabases: []sdk.AccountObjectIdentifier{
				testDb(t).ID(),
			},
		},
	}
	err = client.FailoverGroups.AlterSource(ctx, failoverGroupTest.ID(), opts)
	require.NoError(t, err)
	databases, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroupTest.ID())
	require.NoError(t, err)
	assert.Equal(t, 1, len(databases))
	assert.Equal(t, testDb(t).ID(), databases[0])
}

func TestInt_FailoverGroupsShowShares(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := createFailoverGroup(t, client)
	t.Cleanup(failoverGroupCleanup)

	shareTest, shareCleanup := createShare(t, client)
	t.Cleanup(shareCleanup)
	opts := &sdk.AlterSourceFailoverGroupOptions{
		Set: &sdk.FailoverGroupSet{
			ObjectTypes: []sdk.PluralObjectType{
				sdk.PluralObjectTypeShares,
			},
		},
	}
	err := client.FailoverGroups.AlterSource(ctx, failoverGroupTest.ID(), opts)
	require.NoError(t, err)
	opts = &sdk.AlterSourceFailoverGroupOptions{
		Add: &sdk.FailoverGroupAdd{
			AllowedShares: []sdk.AccountObjectIdentifier{
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
