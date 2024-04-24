package helpers

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
	"testing"
)

type RowAccessPolicyClient struct {
	context *TestClientContext
}

func NewRowAccessPolicyClient(context *TestClientContext) *RowAccessPolicyClient {
	return &RowAccessPolicyClient{
		context: context,
	}
}

func (c *RowAccessPolicyClient) client() sdk.RowAccessPolicies {
	return c.context.client.RowAccessPolicies
}

func (c *RowAccessPolicyClient) createRowAccessPolicy(t *testing.T) (*sdk.RowAccessPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.context.newSchemaObjectIdentifier(random.AlphanumericN(12))
	arg := sdk.NewCreateRowAccessPolicyArgsRequest("A", sdk.DataTypeNumber)
	body := "true"
	createRequest := sdk.NewCreateRowAccessPolicyRequest(id, []sdk.CreateRowAccessPolicyArgsRequest{*arg}, body)

	err := c.client().Create(ctx, createRequest)
	require.NoError(t, err)

	rowAccessPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return rowAccessPolicy, c.DropRowAccessPolicyFunc(t, id)
}

func (c *RowAccessPolicyClient) DropRowAccessPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}
