package helpers

import (
	"context"
	"testing"

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
	return database, func() {
		err := d.client().Drop(ctx, id, &sdk.DropDatabaseOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
		err = d.context.client.Sessions.UseSchema(ctx, sdk.NewDatabaseObjectIdentifier(d.context.database, d.context.schema))
		require.NoError(t, err)
	}
}
