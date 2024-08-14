package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type FunctionClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewFunctionClient(context *TestClientContext, idsGenerator *IdsGenerator) *FunctionClient {
	return &FunctionClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *FunctionClient) client() sdk.Functions {
	return c.context.client.Functions
}

func (c *FunctionClient) Create(t *testing.T, arguments ...sdk.DataType) *sdk.Function {
	t.Helper()
	return c.CreateWithIdentifier(t, c.ids.RandomSchemaObjectIdentifierWithArguments(arguments...))
}

func (c *FunctionClient) CreateWithIdentifier(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) *sdk.Function {
	t.Helper()
	ctx := context.Background()
	argumentRequests := make([]sdk.FunctionArgumentRequest, len(id.ArgumentDataTypes()))
	for i, argumentDataType := range id.ArgumentDataTypes() {
		argumentRequests[i] = *sdk.NewFunctionArgumentRequest(c.ids.Alpha(), argumentDataType)
	}
	err := c.client().CreateForSQL(ctx,
		sdk.NewCreateForSQLFunctionRequest(
			id.SchemaObjectId(),
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeInt)),
			"SELECT 1",
		).WithArguments(argumentRequests),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, c.context.client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id).WithIfExists(true)))
	})

	function, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return function
}
