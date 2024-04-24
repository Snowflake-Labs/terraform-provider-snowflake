package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type PipeClient struct {
	context *TestClientContext
}

func NewPipeClient(context *TestClientContext) *PipeClient {
	return &PipeClient{
		context: context,
	}
}

func (c *PipeClient) client() sdk.Pipes {
	return c.context.client.Pipes
}
