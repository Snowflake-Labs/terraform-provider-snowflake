package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type FailoverGroupClient struct {
	context *TestClientContext
}

func NewFailoverGroupClient(context *TestClientContext) *FailoverGroupClient {
	return &FailoverGroupClient{
		context: context,
	}
}

func (c *FailoverGroupClient) client() sdk.FailoverGroups {
	return c.context.client.FailoverGroups
}
