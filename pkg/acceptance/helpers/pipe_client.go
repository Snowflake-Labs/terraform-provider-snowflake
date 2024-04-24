package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
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

func (c *PipeClient) CreatePipe(t *testing.T, copyStatement string) (*sdk.Pipe, func()) {
	t.Helper()
	return c.CreatePipeInSchema(t, sdk.NewDatabaseObjectIdentifier(c.context.database, c.context.schema), copyStatement)
}

func (c *PipeClient) CreatePipeInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, copyStatement string) (*sdk.Pipe, func()) {
	t.Helper()
	ctx := context.Background()

	id := sdk.NewSchemaObjectIdentifier(schemaId.DatabaseName(), schemaId.Name(), random.AlphaN(20))

	err := c.client().Create(ctx, id, copyStatement, &sdk.CreatePipeOptions{})
	require.NoError(t, err)

	pipe, err := c.client().Describe(ctx, id)
	require.NoError(t, err)

	return pipe, c.DropPipeFunc(t, id)
}

func (c *PipeClient) DropPipeFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropPipeOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}
