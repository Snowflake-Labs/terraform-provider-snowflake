package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

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

func (c *ParameterClient) ShowAccountParameters(t *testing.T) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Account: sdk.Bool(true),
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowDatabaseParameters(t *testing.T, id sdk.AccountObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Database: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowWarehouseParameters(t *testing.T, id sdk.AccountObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Warehouse: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowSchemaParameters(t *testing.T, id sdk.DatabaseObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Schema: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowUserParameters(t *testing.T, id sdk.AccountObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			User: id,
		},
	})
	require.NoError(t, err)
	return params
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

func FindParameter[T ~string](t *testing.T, parameters []*sdk.Parameter, parameter T) *sdk.Parameter {
	t.Helper()
	param, err := collections.FindFirst(parameters, func(p *sdk.Parameter) bool { return p.Key == string(parameter) })
	require.NoError(t, err)
	return *param
}
