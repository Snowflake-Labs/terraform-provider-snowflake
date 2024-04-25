package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ApplicationClient struct {
	context *TestClientContext
}

func NewApplicationClient(context *TestClientContext) *ApplicationClient {
	return &ApplicationClient{
		context: context,
	}
}

func (c *ApplicationClient) client() sdk.Applications {
	return c.context.client.Applications
}

func (c *ApplicationClient) CreateApplication(t *testing.T, packageId sdk.AccountObjectIdentifier, version string) (*sdk.Application, func()) {
	t.Helper()
	ctx := context.Background()

	id := sdk.NewAccountObjectIdentifier(random.AlphaN(8))
	err := c.client().Create(ctx, sdk.NewCreateApplicationRequest(id, packageId).WithVersion(sdk.NewApplicationVersionRequest().WithVersionAndPatch(sdk.NewVersionAndPatchRequest(version, nil))))
	require.NoError(t, err)

	application, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return application, c.DropApplicationFunc(t, id)
}

func (c *ApplicationClient) DropApplicationFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropApplicationRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}
