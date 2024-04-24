package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ShareClient struct {
	context *TestClientContext
}

func NewShareClient(context *TestClientContext) *ShareClient {
	return &ShareClient{
		context: context,
	}
}

func (c *ShareClient) client() sdk.Shares {
	return c.context.client.Shares
}
