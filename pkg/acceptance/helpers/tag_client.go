package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type TagClient struct {
	context *TestClientContext
}

func NewTagClient(context *TestClientContext) *TagClient {
	return &TagClient{
		context: context,
	}
}

func (c *TagClient) client() sdk.Tags {
	return c.context.client.Tags
}
