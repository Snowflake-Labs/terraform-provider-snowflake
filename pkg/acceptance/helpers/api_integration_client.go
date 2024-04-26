package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ApiIntegrationClient struct {
	context *TestClientContext
}

func NewApiIntegrationClient(context *TestClientContext) *ApiIntegrationClient {
	return &ApiIntegrationClient{
		context: context,
	}
}

func (c *ApiIntegrationClient) client() sdk.ApiIntegrations {
	return c.context.client.ApiIntegrations
}

func (c *ApiIntegrationClient) CreateApiIntegration(t *testing.T) (*sdk.ApiIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	id := sdk.NewAccountObjectIdentifier(random.AlphanumericN(12))
	apiAllowedPrefixes := []sdk.ApiIntegrationEndpointPrefix{{Path: "https://xyz.execute-api.us-west-2.amazonaws.com/production"}}
	req := sdk.NewCreateApiIntegrationRequest(id, apiAllowedPrefixes, true)
	req.WithAwsApiProviderParams(sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiGateway, "arn:aws:iam::123456789012:role/hello_cloud_account_role"))

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	apiIntegration, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return apiIntegration, c.DropApiIntegrationFunc(t, id)
}

func (c *ApiIntegrationClient) DropApiIntegrationFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropApiIntegrationRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}
