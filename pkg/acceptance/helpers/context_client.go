package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ContextClient struct {
	context *TestClientContext
}

func NewContextClient(context *TestClientContext) *ContextClient {
	return &ContextClient{
		context: context,
	}
}

func (c *ContextClient) client() sdk.ContextFunctions {
	return c.context.client.ContextFunctions
}
