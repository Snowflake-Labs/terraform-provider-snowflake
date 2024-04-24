package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ShareClient struct {
	context *TestClientContext
}

func NewShareClient(context *TestClientContext) *ShareClient {
	return &ShareClient{
		context: context,
	}
}

func (c *ShareClient) client() sdk.Shares {
	return c.context.client.Shares
}

func (c *ShareClient) CreateShare(t *testing.T) (*sdk.Share, func()) {
	t.Helper()
	// TODO(SNOW-1058419): Try with identifier containing dot during identifiers rework
	id := sdk.RandomAlphanumericAccountObjectIdentifier()
	return c.CreateShareWithOptions(t, id, &sdk.CreateShareOptions{})
}

func (c *ShareClient) CreateShareWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateShareOptions) (*sdk.Share, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)

	share, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return share, c.DropShareFunc(t, id)
}

func (c *ShareClient) DropShareFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// there is no if exists in the docs https://docs.snowflake.com/en/sql-reference/sql/drop-share
		err := c.client().Drop(ctx, id)
		require.NoError(t, err)
	}
}
