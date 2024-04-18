package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
