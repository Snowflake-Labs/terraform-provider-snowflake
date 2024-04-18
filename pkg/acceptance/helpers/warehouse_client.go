package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type WarehouseClient struct {
	context *TestClientContext
}

func NewWarehouseClient(context *TestClientContext) *WarehouseClient {
	return &WarehouseClient{
		context: context,
	}
}

func (c *WarehouseClient) client() sdk.Warehouses {
	return c.context.client.Warehouses
}

func (c *WarehouseClient) UseWarehouse(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	err := c.context.client.Sessions.UseWarehouse(ctx, id)
	require.NoError(t, err)
	return func() {
		err = c.context.client.Sessions.UseWarehouse(ctx, sdk.NewAccountObjectIdentifier(c.context.warehouse))
		require.NoError(t, err)
	}
}

func (c *WarehouseClient) CreateWarehouse(t *testing.T) (*sdk.Warehouse, func()) {
	t.Helper()
	return c.CreateWarehouseWithOptions(t, &sdk.CreateWarehouseOptions{})
}

func (c *WarehouseClient) CreateWarehouseWithOptions(t *testing.T, opts *sdk.CreateWarehouseOptions) (*sdk.Warehouse, func()) {
	t.Helper()
	name := random.StringRange(8, 28)
	id := sdk.NewAccountObjectIdentifier(name)
	ctx := context.Background()
	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)
	return &sdk.Warehouse{Name: name}, c.DropWarehouseFunc(t, id)
}

func (c *WarehouseClient) DropWarehouseFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropWarehouseOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
		err = c.context.client.Sessions.UseWarehouse(ctx, sdk.NewAccountObjectIdentifier(c.context.warehouse))
		require.NoError(t, err)
	}
}
