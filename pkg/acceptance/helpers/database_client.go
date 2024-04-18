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
}

func NewDatabaseClient(context *TestClientContext) *DatabaseClient {
	return &DatabaseClient{
		context: context,
	}
}

func (d *DatabaseClient) client() sdk.Databases {
	return d.context.client.Databases
}

func (d *DatabaseClient) CreateDatabase(t *testing.T) (*sdk.Database, func()) {
	t.Helper()
	return d.CreateDatabaseWithOptions(t, sdk.RandomAccountObjectIdentifier(), &sdk.CreateDatabaseOptions{})
}

func (d *DatabaseClient) CreateDatabaseWithName(t *testing.T, name string) (*sdk.Database, func()) {
	t.Helper()
	return d.CreateDatabaseWithOptions(t, sdk.NewAccountObjectIdentifier(name), &sdk.CreateDatabaseOptions{})
}

func (d *DatabaseClient) CreateDatabaseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()
	err := d.client().Create(ctx, id, opts)
	require.NoError(t, err)
	database, err := d.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return database, d.DropDatabaseFunc(t, id)
}

func (d *DatabaseClient) DropDatabaseFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := d.client().Drop(ctx, id, &sdk.DropDatabaseOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
		err = d.context.client.Sessions.UseSchema(ctx, sdk.NewDatabaseObjectIdentifier(d.context.database, d.context.schema))
		require.NoError(t, err)
	}
}

func (d *DatabaseClient) CreateSecondaryDatabaseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, externalId sdk.ExternalObjectIdentifier, opts *sdk.CreateSecondaryDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes creating secondary db right after primary creation resulted in error
	time.Sleep(1 * time.Second)

	err := d.client().CreateSecondary(ctx, id, externalId, opts)
	require.NoError(t, err)

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes secondary database is not shown as SHOW REPLICATION DATABASES results right after creation
	time.Sleep(1 * time.Second)

	database, err := d.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return database, func() {
		err := d.client().Drop(ctx, id, nil)
		require.NoError(t, err)

		// TODO [926148]: make this wait better with tests stabilization
		// waiting because sometimes dropping primary db right after dropping the secondary resulted in error
		time.Sleep(1 * time.Second)
		err = d.context.client.Sessions.UseSchema(ctx, sdk.NewDatabaseObjectIdentifier(d.context.database, d.context.schema))
		require.NoError(t, err)
	}
}

func (d *DatabaseClient) UpdateDataRetentionTime(t *testing.T, id sdk.AccountObjectIdentifier, days int) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := d.client().Alter(ctx, id, &sdk.AlterDatabaseOptions{
			Set: &sdk.DatabaseSet{
				DataRetentionTimeInDays: sdk.Int(days),
			},
		})
		require.NoError(t, err)
	}
}
