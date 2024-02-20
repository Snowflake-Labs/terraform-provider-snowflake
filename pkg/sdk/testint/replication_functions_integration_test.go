package testint

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ShowReplicationFunctions(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	accounts, err := client.ReplicationFunctions.ShowReplicationAccounts(ctx)
	if err != nil {
		t.Skip("replication not enabled in this account")
	}
	assert.NotEmpty(t, accounts)
}

func TestInt_ShowReplicationDatabases(t *testing.T) {
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	account, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	accountId := sdk.NewAccountIdentifierFromAccountLocator(account)

	secondaryAccount, err := secondaryClient.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	secondaryAccountId := sdk.NewAccountIdentifierFromAccountLocator(secondaryAccount)

	db1Name := random.AlphaN(10)
	db2Name := random.AlphaN(10)
	db3Name := random.AlphaN(10)
	db, dbCleanup := createDatabaseWithOptions(t, client, sdk.NewAccountObjectIdentifier(db1Name), &sdk.CreateDatabaseOptions{})
	t.Cleanup(dbCleanup)
	db2, dbCleanup2 := createDatabaseWithOptions(t, client, sdk.NewAccountObjectIdentifier(db2Name), &sdk.CreateDatabaseOptions{})
	t.Cleanup(dbCleanup2)

	err = client.Databases.AlterReplication(ctx, db.ID(), &sdk.AlterDatabaseReplicationOptions{EnableReplication: &sdk.EnableReplication{ToAccounts: []sdk.AccountIdentifier{secondaryAccountId}}})
	require.NoError(t, err)
	err = client.Databases.AlterReplication(ctx, db2.ID(), &sdk.AlterDatabaseReplicationOptions{EnableReplication: &sdk.EnableReplication{ToAccounts: []sdk.AccountIdentifier{secondaryAccountId}}})
	require.NoError(t, err)

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes creating secondary db right after primary creation resulted in error
	time.Sleep(1 * time.Second)
	db3, dbCleanup3 := createSecondaryDatabaseWithOptions(t, secondaryClient, sdk.NewAccountObjectIdentifier(db3Name), sdk.NewExternalObjectIdentifier(accountId, db.ID()), &sdk.CreateSecondaryDatabaseOptions{})
	t.Cleanup(dbCleanup3)

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes secondary database is not shown as SHOW REPLICATION DATABASES results right after creation
	time.Sleep(1 * time.Second)

	getByName := func(replicationDatabases []sdk.ReplicationDatabase, name sdk.AccountObjectIdentifier) *sdk.ReplicationDatabase {
		for _, rdb := range replicationDatabases {
			rdb := rdb
			if rdb.Name == name.Name() {
				return &rdb
			}
		}
		return nil
	}

	assertReplicationDatabase := func(rdb *sdk.ReplicationDatabase, expectedName string, expectedIsPrimary bool) {
		require.NotNil(t, rdb)
		require.Equal(t, expectedName, rdb.Name)
		require.Equal(t, expectedIsPrimary, rdb.IsPrimary)
		require.NotEmpty(t, rdb.SnowflakeRegion)
		require.NotEmpty(t, rdb.CreatedOn)
		require.NotEmpty(t, rdb.AccountName)
		require.NotEmpty(t, rdb.PrimaryDatabase)
		if expectedIsPrimary {
			require.NotEmpty(t, rdb.ReplicationAllowedToAccounts)
			require.NotEmpty(t, rdb.FailoverAllowedToAccounts)
		}
		require.NotEmpty(t, rdb.OrganizationName)
		require.NotEmpty(t, rdb.AccountName)
	}

	t.Run("no options", func(t *testing.T) {
		opts := &sdk.ShowReplicationDatabasesOptions{}
		replicationDatabases, err := client.ReplicationFunctions.ShowReplicationDatabases(ctx, opts)
		require.NoError(t, err)

		rdb := getByName(replicationDatabases, db.ID())
		assertReplicationDatabase(rdb, db.Name, true)

		rdb2 := getByName(replicationDatabases, db2.ID())
		assertReplicationDatabase(rdb2, db2.Name, true)

		rdb3 := getByName(replicationDatabases, db3.ID())
		assertReplicationDatabase(rdb3, db3.Name, false)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &sdk.ShowReplicationDatabasesOptions{
			Like: &sdk.Like{Pattern: &db.Name},
		}
		replicationDatabases, err := client.ReplicationFunctions.ShowReplicationDatabases(ctx, opts)
		require.NoError(t, err)

		require.Len(t, replicationDatabases, 1)
		require.Equal(t, db.Name, replicationDatabases[0].Name)

		opts = &sdk.ShowReplicationDatabasesOptions{
			Like: &sdk.Like{Pattern: &db2.Name},
		}
		replicationDatabases, err = client.ReplicationFunctions.ShowReplicationDatabases(ctx, opts)
		require.NoError(t, err)

		require.Len(t, replicationDatabases, 1)
		require.Equal(t, db2.Name, replicationDatabases[0].Name)
	})

	t.Run("with primary", func(t *testing.T) {
		opts := &sdk.ShowReplicationDatabasesOptions{
			WithPrimary: sdk.Pointer(sdk.NewExternalObjectIdentifier(accountId, db.ID())),
		}
		replicationDatabases, err := client.ReplicationFunctions.ShowReplicationDatabases(ctx, opts)
		require.NoError(t, err)

		require.Len(t, replicationDatabases, 2)

		primary := getByName(replicationDatabases, db.ID())
		require.Equal(t, db.Name, primary.Name)
		secondary := getByName(replicationDatabases, db3.ID())
		require.Equal(t, db3.Name, secondary.Name)
	})
}

func TestInt_ShowRegions(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	t.Run("no options", func(t *testing.T) {
		regions, err := client.ReplicationFunctions.ShowRegions(ctx, nil)
		require.NoError(t, err)
		assert.NotEmpty(t, regions)
	})

	t.Run("with options", func(t *testing.T) {
		regions, err := client.ReplicationFunctions.ShowRegions(ctx, &sdk.ShowRegionsOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("AWS_US_WEST_2"),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(regions))
		region := regions[0]
		assert.Equal(t, "AWS_US_WEST_2", region.SnowflakeRegion)
		assert.Equal(t, sdk.CloudTypeAWS, region.CloudType)
		assert.Equal(t, "us-west-2", region.Region)
		assert.Equal(t, "US West (Oregon)", region.DisplayName)
	})
}
