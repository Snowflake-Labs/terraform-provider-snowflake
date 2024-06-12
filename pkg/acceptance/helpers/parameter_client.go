package helpers

import (
	"context"
	"fmt"
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

	param := c.ShowAccountParameter(t, parameter)
	oldValue := param.Value
	oldLevel := param.Level

	err := c.client().SetAccountParameter(ctx, parameter, newValue)
	require.NoError(t, err)

	return func() {
		if oldLevel == "" {
			c.UnsetAccountParameter(t, parameter)
		} else {
			err := c.client().SetAccountParameter(ctx, parameter, oldValue)
			require.NoError(t, err)
		}
	}
}

func (c *ParameterClient) ShowAccountParameter(t *testing.T, parameter sdk.AccountParameter) *sdk.Parameter {
	t.Helper()
	ctx := context.Background()

	param, err := c.client().ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)

	return param
}

// TODO [SNOW-1473408]: add unset account parameter to sdk.Parameters
func (c *ParameterClient) UnsetAccountParameter(t *testing.T, parameter sdk.AccountParameter) {
	t.Helper()
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf("ALTER ACCOUNT UNSET %s", parameter))
	require.NoError(t, err)
}
