package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ViewClient struct {
	context *TestClientContext
}

func NewViewClient(context *TestClientContext) *ViewClient {
	return &ViewClient{
		context: context,
	}
}

func (c *ViewClient) client() sdk.Views {
	return c.context.client.Views
}
