package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type WarehouseClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewWarehouseClient(context *TestClientContext, idsGenerator *IdsGenerator) *WarehouseClient {
	return &WarehouseClient{
		context: context,
		ids:     idsGenerator,
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
		err = c.context.client.Sessions.UseWarehouse(ctx, c.ids.WarehouseId())
		require.NoError(t, err)
	}
}

func (c *WarehouseClient) CreateWarehouse(t *testing.T) (*sdk.Warehouse, func()) {
	t.Helper()
	return c.CreateWarehouseWithOptions(t, c.ids.RandomAccountObjectIdentifier(), &sdk.CreateWarehouseOptions{})
}

func (c *WarehouseClient) CreateWarehouseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateWarehouseOptions) (*sdk.Warehouse, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)
	return &sdk.Warehouse{Name: id.Name()}, c.DropWarehouseFunc(t, id)
}

func (c *WarehouseClient) DropWarehouseFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropWarehouseOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
		err = c.context.client.Sessions.UseWarehouse(ctx, c.ids.WarehouseId())
		require.NoError(t, err)
	}
}

func (c *WarehouseClient) UpdateMaxConcurrencyLevel(t *testing.T, id sdk.AccountObjectIdentifier, level int) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{MaxConcurrencyLevel: sdk.Int(level)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateWarehouseSize(t *testing.T, id sdk.AccountObjectIdentifier, newSize sdk.WarehouseSize) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{WarehouseSize: sdk.Pointer(newSize)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Warehouse, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
