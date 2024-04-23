package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ApplicationClient struct {
	context *TestClientContext
}

func NewApplicationClient(context *TestClientContext) *ApplicationClient {
	return &ApplicationClient{
		context: context,
	}
}

func (c *ApplicationClient) client() sdk.Applications {
	return c.context.client.Applications
}
