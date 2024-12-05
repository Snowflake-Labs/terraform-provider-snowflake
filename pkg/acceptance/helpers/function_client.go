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

	return c.CreateWithRequest(t, id,
		sdk.NewCreateForSQLFunctionRequest(
			id.SchemaObjectId(),
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeInt)),
			"SELECT 1",
		),
	)
}

func (c *FunctionClient) CreateSecure(t *testing.T, arguments ...sdk.DataType) *sdk.Function {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifierWithArguments(arguments...)

	return c.CreateWithRequest(t, id,
		sdk.NewCreateForSQLFunctionRequest(
			id.SchemaObjectId(),
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeInt)),
			"SELECT 1",
		).WithSecure(true),
	)
}

func (c *FunctionClient) CreateWithRequest(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, req *sdk.CreateForSQLFunctionRequest) *sdk.Function {
	t.Helper()
	ctx := context.Background()
	argumentRequests := make([]sdk.FunctionArgumentRequest, len(id.ArgumentDataTypes()))
	for i, argumentDataType := range id.ArgumentDataTypes() {
		argumentRequests[i] = *sdk.NewFunctionArgumentRequest(c.ids.Alpha(), nil).WithArgDataTypeOld(argumentDataType)
	}
	err := c.client().CreateForSQL(ctx, req.WithArguments(argumentRequests))
	require.NoError(t, err)

	t.Cleanup(c.DropFunctionFunc(t, id))

	function, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return function
}

func (c *FunctionClient) DropFunctionFunc(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropFunctionRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *FunctionClient) Show(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Function, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *FunctionClient) SampleJavaDefinition(t *testing.T) string {
	t.Helper()

	return `
	class TestFunc {
		public static String echoVarchar(String x) {
			return x;
		}
	}
`
}

func (c *FunctionClient) SampleJavaScriptDefinition(t *testing.T) string {
	t.Helper()

	return `
	if (D <= 0) {
		return 1;
	} else {
		var result = 1;
		for (var i = 2; i <= D; i++) {
			result = result * i;
		}
		return result;
	}
`
}

func (c *FunctionClient) SamplePythonDefinition(t *testing.T) string {
	t.Helper()

	return `
	def dump(i):
	print("Hello World!")
`
}

func (c *FunctionClient) SampleScalaDefinition(t *testing.T) string {
	t.Helper()

	return `
	class Echo {
		def echoVarchar(x : String): String = {
			return x
		}
	}
`
}

func (c *FunctionClient) SampleSqlDefinition(t *testing.T) string {
	t.Helper()

	return "3.141592654::FLOAT"
}
