package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TagClient struct {
	context *TestClientContext
}

func NewTagClient(context *TestClientContext) *TagClient {
	return &TagClient{
		context: context,
	}
}

func (c *TagClient) client() sdk.Tags {
	return c.context.client.Tags
}

func (c *TagClient) CreateTag(t *testing.T) (*sdk.Tag, func()) {
	t.Helper()
	return c.CreateTagInSchema(t, sdk.NewDatabaseObjectIdentifier(c.context.database, c.context.schema))
}

func (c *TagClient) CreateTagInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Tag, func()) {
	t.Helper()
	ctx := context.Background()

	id := sdk.NewSchemaObjectIdentifier(schemaId.DatabaseName(), schemaId.Name(), random.AlphanumericN(12))

	err := c.client().Create(ctx, sdk.NewCreateTagRequest(id))
	require.NoError(t, err)

	tag, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return tag, c.DropTagFunc(t, id)
}

func (c *TagClient) DropTagFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropTagRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
