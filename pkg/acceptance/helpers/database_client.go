package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type DatabaseClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewDatabaseClient(context *TestClientContext, idsGenerator *IdsGenerator) *DatabaseClient {
	return &DatabaseClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *DatabaseClient) client() sdk.Databases {
	return c.context.client.Databases
}

func (c *DatabaseClient) CreateDatabase(t *testing.T) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithOptions(t, c.ids.RandomAccountObjectIdentifier(), &sdk.CreateDatabaseOptions{})
}

func (c *DatabaseClient) CreateDatabaseWithIdentifier(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithOptions(t, id, &sdk.CreateDatabaseOptions{})
}

func (c *DatabaseClient) CreateDatabaseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)

	database, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return database, c.DropDatabaseFunc(t, id)
}

func (c *DatabaseClient) DropDatabaseFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	return func() { require.NoError(t, c.DropDatabase(t, id)) }
}

func (c *DatabaseClient) DropDatabase(t *testing.T, id sdk.AccountObjectIdentifier) error {
	t.Helper()
	ctx := context.Background()

	if err := c.client().Drop(ctx, id, &sdk.DropDatabaseOptions{IfExists: sdk.Bool(true)}); err != nil {
		return err
	}
	if err := c.context.client.Sessions.UseSchema(ctx, c.ids.SchemaId()); err != nil {
		return err
	}
	return nil
}

func (c *DatabaseClient) CreateSecondaryDatabaseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, externalId sdk.ExternalObjectIdentifier, opts *sdk.CreateSecondaryDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes creating secondary db right after primary creation resulted in error
	time.Sleep(1 * time.Second)

	err := c.client().CreateSecondary(ctx, id, externalId, opts)
	require.NoError(t, err)

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes secondary database is not shown as SHOW REPLICATION DATABASES results right after creation
	time.Sleep(1 * time.Second)

	database, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return database, func() {
		err := c.client().Drop(ctx, id, nil)
		require.NoError(t, err)

		// TODO [926148]: make this wait better with tests stabilization
		// waiting because sometimes dropping primary db right after dropping the secondary resulted in error
		time.Sleep(1 * time.Second)
		err = c.context.client.Sessions.UseSchema(ctx, c.ids.SchemaId())
		require.NoError(t, err)
	}
}

func (c *DatabaseClient) CreatePrimaryDatabase(t *testing.T, enableReplicationTo []sdk.AccountIdentifier) (*sdk.Database, sdk.ExternalObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	primaryDatabase, primaryDatabaseCleanup := c.CreateDatabase(t)

	err := c.client().AlterReplication(ctx, primaryDatabase.ID(), &sdk.AlterDatabaseReplicationOptions{
		EnableReplication: &sdk.EnableReplication{
			ToAccounts:         enableReplicationTo,
			IgnoreEditionCheck: sdk.Bool(true),
		},
	})
	require.NoError(t, err)

	sessionDetails, err := c.context.client.ContextFunctions.CurrentSessionDetails(ctx)
	require.NoError(t, err)

	externalPrimaryId := sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName), primaryDatabase.ID())
	return primaryDatabase, externalPrimaryId, primaryDatabaseCleanup
}

func (c *DatabaseClient) UpdateDataRetentionTime(t *testing.T, id sdk.AccountObjectIdentifier, days int) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterDatabaseOptions{
		Set: &sdk.DatabaseSet{
			DataRetentionTimeInDays: sdk.Int(days),
		},
	})
	require.NoError(t, err)
}

func (c *DatabaseClient) UnsetCatalog(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterDatabaseOptions{
		Unset: &sdk.DatabaseUnset{
			Catalog: sdk.Bool(true),
		},
	})
	require.NoError(t, err)
}

func (c *DatabaseClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Database, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
