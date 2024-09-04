package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
	"testing"
)

type BcrBundlesClient struct {
	context *TestClientContext
}

func NewBcrBundlesClient(context *TestClientContext) *BcrBundlesClient {
	return &BcrBundlesClient{
		context: context,
	}
}

func (c *BcrBundlesClient) client() sdk.SystemFunctions {
	return c.context.client.SystemFunctions
}

func (c *BcrBundlesClient) EnableBcrBundle(t *testing.T, name string) {
	t.Helper()

	err := c.client().EnableBehaviorChangeBundle(name)
	require.NoError(t, err)

	t.Cleanup(c.DisableBcrBundleFunc(t, name))
}

func (c *BcrBundlesClient) DisableBcrBundleFunc(t *testing.T, name string) func() {
	t.Helper()

	return func() {
		err := c.client().DisableBehaviorChangeBundle(name)
		require.NoError(t, err)
	}
}
