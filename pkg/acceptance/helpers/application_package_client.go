package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ApplicationPackageClient struct {
	context *TestClientContext
}

func NewApplicationPackageClient(context *TestClientContext) *ApplicationPackageClient {
	return &ApplicationPackageClient{
		context: context,
	}
}

func (c *ApplicationPackageClient) client() sdk.ApplicationPackages {
	return c.context.client.ApplicationPackages
}

func (c *ApplicationPackageClient) CreateApplicationPackage(t *testing.T) (*sdk.ApplicationPackage, func()) {
	t.Helper()
	ctx := context.Background()

	id := sdk.NewAccountObjectIdentifier(random.AlphaN(8))
	err := c.client().Create(ctx, sdk.NewCreateApplicationPackageRequest(id))
	require.NoError(t, err)

	applicationPackage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return applicationPackage, c.DropApplicationPackageFunc(t, id)
}

func (c *ApplicationPackageClient) DropApplicationPackageFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropApplicationPackageRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}

func (c *ApplicationPackageClient) AddApplicationPackageVersion(t *testing.T, id sdk.AccountObjectIdentifier, stageId sdk.SchemaObjectIdentifier, versionName string) {
	t.Helper()
	ctx := context.Background()

	using := "@" + stageId.FullyQualifiedName()

	err := c.client().Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithAddVersion(sdk.NewAddVersionRequest(using).WithVersionIdentifier(sdk.String(versionName))))
	require.NoError(t, err)
}
