package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ViewClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewViewClient(context *TestClientContext, idsGenerator *IdsGenerator) *ViewClient {
	return &ViewClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ViewClient) client() sdk.Views {
	return c.context.client.Views
}

func (c *ViewClient) CreateView(t *testing.T, query string) (*sdk.View, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()

	err := c.client().Create(ctx, sdk.NewCreateViewRequest(id, query))
	require.NoError(t, err)

	view, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return view, c.DropViewFunc(t, id)
}

func (c *ViewClient) DropViewFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropViewRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}
