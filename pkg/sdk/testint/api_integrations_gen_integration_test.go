package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ApiIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// TODO [JIRA]: replace with real values?
	const awsAllowedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	const awsBlockedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/prod/"
	const azureAllowedPrefix = "https://apim-hello-world.azure-api.net/dev"
	const azureBlockedPrefix = "https://apim-hello-world.azure-api.net/prod"
	const googleAllowedPrefix = "https://gateway-id-123456.uc.gateway.dev/prod"
	const googleBlockedPrefix = "https://gateway-id-123456.uc.gateway.dev/dev"
	const apiAwsRoleArn = "arn:aws:iam::000000000001:/role/test"
	const azureTenantId = "00000000-0000-0000-0000-000000000000"
	const azureAdApplicationId = "11111111-1111-1111-1111-111111111111"
	const googleAudience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"

	prefixes := func(prefix string) []sdk.ApiIntegrationEndpointPrefix {
		return []sdk.ApiIntegrationEndpointPrefix{{Path: prefix}}
	}
	assertApiIntegration := func(t *testing.T, s *sdk.ApiIntegration, name sdk.AccountObjectIdentifier, comment string) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.Equal(t, true, s.Enabled)
		assert.Equal(t, "EXTERNAL_API", s.ApiType)
		assert.Equal(t, "API", s.Category)
		assert.Equal(t, comment, s.Comment)
	}

	cleanupApiIntegrationProvider := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
			require.NoError(t, err)
		}
	}

	createApiIntegrationAwsRequest := func(t *testing.T) *sdk.CreateApiIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		return sdk.NewCreateApiIntegrationRequest(id, prefixes(awsAllowedPrefix), true).
			WithAwsApiProviderParams(sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiGateway, apiAwsRoleArn))
	}

	createApiIntegrationAzureRequest := func(t *testing.T) *sdk.CreateApiIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		return sdk.NewCreateApiIntegrationRequest(id, prefixes(azureAllowedPrefix), true).
			WithAzureApiProviderParams(sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId))
	}

	createApiIntegrationGoogleRequest := func(t *testing.T) *sdk.CreateApiIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		return sdk.NewCreateApiIntegrationRequest(id, prefixes(googleAllowedPrefix), true).
			WithGoogleApiProviderParams(sdk.NewGoogleApiParamsRequest(googleAudience))
	}

	createApiIntegrationWithRequest := func(t *testing.T, request *sdk.CreateApiIntegrationRequest) *sdk.ApiIntegration {
		t.Helper()
		id := request.GetName()

		err := client.ApiIntegrations.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApiIntegrationProvider(id))

		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration
	}

	createAwsApiIntegration := func(t *testing.T) *sdk.ApiIntegration {
		t.Helper()
		return createApiIntegrationWithRequest(t, createApiIntegrationAwsRequest(t))
	}

	createAzureApiIntegration := func(t *testing.T) *sdk.ApiIntegration {
		t.Helper()
		return createApiIntegrationWithRequest(t, createApiIntegrationAzureRequest(t))
	}

	createGoogleApiIntegration := func(t *testing.T) *sdk.ApiIntegration {
		t.Helper()
		return createApiIntegrationWithRequest(t, createApiIntegrationGoogleRequest(t))
	}

	t.Run("create api integration: aws basic", func(t *testing.T) {
		request := createApiIntegrationAwsRequest(t)

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "")
	})

	t.Run("create api integration: azure basic", func(t *testing.T) {
		request := createApiIntegrationAzureRequest(t)

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "")
	})

	t.Run("create api integration: google basic", func(t *testing.T) {
		request := createApiIntegrationGoogleRequest(t)

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "")
	})

	t.Run("create api integration: aws more options", func(t *testing.T) {
		request := createApiIntegrationAwsRequest(t)

		request = request.
			WithAwsApiProviderParams(request.AwsApiProviderParams.WithApiKey(sdk.String("key"))).
			WithApiBlockedPrefixes(prefixes(awsBlockedPrefix)).
			WithComment(sdk.String("comment"))

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "comment")
	})

	t.Run("create api integration: azure more options", func(t *testing.T) {
		request := createApiIntegrationAzureRequest(t)

		request = request.
			WithAzureApiProviderParams(request.AzureApiProviderParams.WithApiKey(sdk.String("key"))).
			WithApiBlockedPrefixes(prefixes(azureBlockedPrefix)).
			WithComment(sdk.String("comment"))

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "comment")
	})

	t.Run("create api integration: google more options", func(t *testing.T) {
		request := createApiIntegrationGoogleRequest(t).
			WithApiBlockedPrefixes(prefixes(googleBlockedPrefix)).
			WithComment(sdk.String("comment"))

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "comment")
	})

	t.Run("alter api integration: aws", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter api integration: azure", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter api integration: google", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop api integration: existing", func(t *testing.T) {
		request := createApiIntegrationAwsRequest(t)
		id := request.GetName()

		err := client.ApiIntegrations.Create(ctx, request)
		require.NoError(t, err)

		err = client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
		require.NoError(t, err)

		_, err = client.ApiIntegrations.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop api integration: non-existing", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")

		err := client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show api integration: default", func(t *testing.T) {
		integrationAws := createAwsApiIntegration(t)
		integrationAzure := createAzureApiIntegration(t)
		integrationGoogle := createGoogleApiIntegration(t)

		showRequest := sdk.NewShowApiIntegrationRequest()
		returnedIntegrations, err := client.ApiIntegrations.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAws)
		assert.Contains(t, returnedIntegrations, *integrationAzure)
		assert.Contains(t, returnedIntegrations, *integrationGoogle)
	})

	t.Run("show api integration: with options", func(t *testing.T) {
		integrationAws := createAwsApiIntegration(t)
		integrationAzure := createAzureApiIntegration(t)

		showRequest := sdk.NewShowApiIntegrationRequest().
			WithLike(&sdk.Like{Pattern: &integrationAws.Name})
		returnedIntegrations, err := client.ApiIntegrations.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAws)
		assert.NotContains(t, returnedIntegrations, *integrationAzure)
	})

	t.Run("describe api integration: aws", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe api integration: azure", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe api integration: google", func(t *testing.T) {
		// TODO: fill me
	})
}
