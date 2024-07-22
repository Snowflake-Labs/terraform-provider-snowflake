package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type StreamlitClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewStreamlitClient(context *TestClientContext, idsGenerator *IdsGenerator) *StreamlitClient {
	return &StreamlitClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *StreamlitClient) client() sdk.Streamlits {
	return c.context.client.Streamlits
}

func (c *StreamlitClient) Update(t *testing.T, request *sdk.AlterStreamlitRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}
