package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type StageClient struct {
	context *TestClientContext
}

func NewStageClient(context *TestClientContext) *StageClient {
	return &StageClient{
		context: context,
	}
}

func (c *StageClient) client() sdk.Stages {
	return c.context.client.Stages
}
