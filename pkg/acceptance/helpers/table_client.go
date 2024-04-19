package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type TableClient struct {
	context *TestClientContext
}

func NewTableClient(context *TestClientContext) *TableClient {
	return &TableClient{
		context: context,
	}
}

func (c *TableClient) client() sdk.Tables {
	return c.context.client.Tables
}
