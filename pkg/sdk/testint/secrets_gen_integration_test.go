package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInt_Secrets(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	integrationId := testClientHelper().Ids.RandomAccountObjectIdentifier()

	cleanupIntegration := func(t *testing.T, integrationId sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(integrationId).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	integrationRequest := sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo").
		WithOauthAllowedScopes([]sdk.AllowedScope{{"foo"}, {"bar"}})
	err := client.SecurityIntegrations.CreateApiAuthenticationWithClientCredentialsFlow(ctx, integrationRequest)
	require.NoError(t, err)
	t.Cleanup(cleanupIntegration(t, integrationId))

	cleanupSecret := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	assertSecret := func(t *testing.T, s *sdk.Secret, expectedId sdk.SchemaObjectIdentifier, expectedSecretType, expectedComment string) {
		t.Helper()
		assert.Equal(t, expectedId.Name(), s.Name)
		assert.Equal(t, expectedSecretType, s.SecretType)
		assert.Equal(t, expectedComment, s.Comment)
		assert.NotEmpty(t, s.CreatedOn)
		assert.NotEmpty(t, s.DatabaseName)
		assert.NotEmpty(t, s.SchemaName)
		assert.NotEmpty(t, s.OwnerRoleType)
		assert.NotEmpty(t, s.Owner)
	}

	createSecretWithOAuthClientCredentialsFlow := func(t *testing.T, integrationId sdk.AccountObjectIdentifier, scopes []sdk.SecurityIntegrationScope, with func(*sdk.CreateWithOAuthClientCredentialsFlowSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, integrationId, scopes)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSecret(id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithOAuthAuthorizationCodeFlow := func(t *testing.T, refreshToken, refreshTokenExpiryTime string, integrationId sdk.AccountObjectIdentifier, with func(*sdk.CreateWithOAuthAuthorizationCodeFlowSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthAuthorizationCodeFlowSecretRequest(id, refreshToken, refreshTokenExpiryTime, integrationId)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithOAuthAuthorizationCodeFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSecret(id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithBasicAuthorization := func(t *testing.T, username, password string, with func(*sdk.CreateWithBasicAuthenticationSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, username, password)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithBasicAuthentication(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSecret(id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	type secretDetails struct {
		Name                        string
		Comment                     string
		SecretType                  string
		Username                    string
		OauthAccessTokenExpiryTime  time.Time
		OauthRefreshTokenExpiryTime time.Time
		OauthScopes                 string
		IntegrationName             string
	}

	assertSecretDetails := func(actual *sdk.SecretDetails, expected secretDetails) {
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Comment, actual.Comment.String)
		assert.Equal(t, expected.SecretType, actual.SecretType)
		assert.Equal(t, expected.Username, actual.Username.String)
		assert.Equal(t, expected.OauthAccessTokenExpiryTime.String(), actual.OauthAccessTokenExpiryTime.String)
		assert.Equal(t, expected.OauthRefreshTokenExpiryTime.String(), actual.OauthRefreshTokenExpiryTime.String)
		assert.Equal(t, expected.OauthScopes, actual.OauthScopes.String)
		assert.Equal(t, expected.IntegrationName, actual.IntegrationName.String)
	}

	t.Run("Create secret with OAuth Client Credentials Flow", func(t *testing.T) {
		scopes := []sdk.SecurityIntegrationScope{{"foo"}, {"bar"}}
		secret, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, scopes, func(req *sdk.CreateWithOAuthClientCredentialsFlowSecretRequest) {
			req.WithComment("a").
				WithIfNotExists(true)
		})
		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			Comment:         "a",
			SecretType:      "OAUTH2",
			OauthScopes:     "[foo, bar]",
			IntegrationName: integrationId.Name(),
		})

		assertSecret(t, secret, id, "OAUTH2", "a")
	})

	// regarding the https://docs.snowflake.com/en/sql-reference/sql/create-secret secret should inherit security_integration scopes, but it does not
	t.Run("CreateSecretWithOAuthClientCredentialsFlow: No Scopes Specified", func(t *testing.T) {
		secret, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{}, nil)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			SecretType:      "OAUTH2",
			IntegrationName: integrationId.Name(),
		})

		assertSecret(t, secret, id, "OAUTH2", "")
	})

	t.Run("CreateWithOAuthAuthorizationCodeFlow", func(t *testing.T) {
		secret, id := createSecretWithOAuthAuthorizationCodeFlow(t, "foo", "2024-09-20", integrationId, func(req *sdk.CreateWithOAuthAuthorizationCodeFlowSecretRequest) {
			req.WithComment("a").
				WithIfNotExists(true)
		})

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  "OAUTH2",
			Comment:                     "a",
			OauthRefreshTokenExpiryTime: time.Time{}, //"2024-09-20",
			IntegrationName:             integrationId.Name(),
		})

		assertSecret(t, secret, id, "OAUTH2", "")
	})

	/*
		//require.EqualError(t, err, "Invalid data/time format string")
		t.Run("CreateWithOAuthAuthorizationCodeFlow: Invalid date/time format string", func(t *testing.T) {
			secret, id := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, func(req *sdk.CreateWithOAuthAuthorizationCodeFlowSecretRequest) {
				req.WithComment("a").
					WithIfNotExists(true)
			})

			details, err := client.Secrets.Describe(ctx, id)
			require.NoError(t, err)

			assertSecretDetails(details, secretDetails{
				Name:                        id.Name(),
				SecretType:                  "OAUTH2",
				OauthAccessTokenExpiryTime:  "foo",
				OauthRefreshTokenExpiryTime: "foo",
				IntegrationName:             integrationId.Name(),
			})

			assertSecret(t, secret, id, "OAUTH2", "")
		})
	*/

	t.Run("CreateWithBasicAuthentication", func(t *testing.T) {
		comment := random.Comment()
		secret, id := createSecretWithBasicAuthorization(t, "foo", "foo", func(req *sdk.CreateWithBasicAuthenticationSecretRequest) {
			req.WithComment(comment).
				WithIfNotExists(true)
		})
		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			Comment:    comment,
			SecretType: "PASSWORD",
			Username:   "foo",
		})

		assertSecret(t, secret, id, "PASSWORD", comment)
	})

	t.Run("CreateWithBasicAuthentication: Empty Username and Password", func(t *testing.T) {
		comment := random.Comment()
		secret, id := createSecretWithBasicAuthorization(t, "", "", func(req *sdk.CreateWithBasicAuthenticationSecretRequest) {
			req.WithComment(comment).
				WithIfNotExists(true)
		})
		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			Comment:    comment,
			SecretType: "PASSWORD",
		})

		assertSecret(t, secret, id, "PASSWORD", comment)
	})

	t.Run("CreateWithGenericString", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Alter", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Drop", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Show", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("ShowByID", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Describe", func(t *testing.T) {
		// TODO: fill me
	})
}
