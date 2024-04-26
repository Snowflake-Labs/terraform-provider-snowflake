package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ParameterClient struct {
	context *TestClientContext
}

func NewParameterClient(context *TestClientContext) *ParameterClient {
	return &ParameterClient{
		context: context,
	}
}

func (c *ParameterClient) client() sdk.Parameters {
	return c.context.client.Parameters
}

func (c *ParameterClient) UpdateAccountParameterTemporarily(t *testing.T, parameter sdk.AccountParameter, newValue string) func() {
	t.Helper()
	ctx := context.Background()

	param, err := c.client().ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)
	oldValue := param.Value

	err = c.client().SetAccountParameter(ctx, parameter, newValue)
	require.NoError(t, err)

	return func() {
		err = c.client().SetAccountParameter(ctx, parameter, oldValue)
		require.NoError(t, err)
	}
}
