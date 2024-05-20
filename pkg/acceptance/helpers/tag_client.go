package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TagClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewTagClient(context *TestClientContext, idsGenerator *IdsGenerator) *TagClient {
	return &TagClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *TagClient) client() sdk.Tags {
	return c.context.client.Tags
}

func (c *TagClient) CreateTag(t *testing.T) (*sdk.Tag, func()) {
	t.Helper()
	return c.CreateTagInSchema(t, c.ids.SchemaId())
}

func (c *TagClient) CreateTagInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Tag, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)

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
