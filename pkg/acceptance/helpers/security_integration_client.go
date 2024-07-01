package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type SecurityIntegrationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSecurityIntegrationClient(context *TestClientContext, idsGenerator *IdsGenerator) *SecurityIntegrationClient {
	return &SecurityIntegrationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SecurityIntegrationClient) client() sdk.SecurityIntegrations {
	return c.context.client.SecurityIntegrations
}

func (c *SecurityIntegrationClient) CreateSaml2(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.SecurityIntegration, func()) {
	t.Helper()
	return c.CreateSaml2WithRequest(t, sdk.NewCreateSaml2SecurityIntegrationRequest(id, false, c.ids.Alpha(), "https://example.com", "Custom", random.GenerateX509(t)))
}

func (c *SecurityIntegrationClient) CreateSaml2WithRequest(t *testing.T, request *sdk.CreateSaml2SecurityIntegrationRequest) (*sdk.SecurityIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateSaml2(ctx, request)
	require.NoError(t, err)

	si, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return si, c.DropSecurityIntegrationFunc(t, request.GetName())
}

func (c *SecurityIntegrationClient) CreateScim(t *testing.T) (*sdk.SecurityIntegration, func()) {
	t.Helper()
	return c.CreateScimWithRequest(t, sdk.NewCreateScimSecurityIntegrationRequest(c.ids.RandomAccountObjectIdentifier(), sdk.ScimSecurityIntegrationScimClientGeneric, sdk.ScimSecurityIntegrationRunAsRoleGenericScimProvisioner))
}

func (c *SecurityIntegrationClient) CreateScimWithRequest(t *testing.T, request *sdk.CreateScimSecurityIntegrationRequest) (*sdk.SecurityIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateScim(ctx, request)
	require.NoError(t, err)

	si, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return si, c.DropSecurityIntegrationFunc(t, request.GetName())
}

func (c *SecurityIntegrationClient) DropSecurityIntegrationFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
