package testint

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const Day = 24 * time.Hour

func TestInt_Secrets(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	integrationId := testClientHelper().Ids.RandomAccountObjectIdentifier()

	refreshTokenExpiryTime := time.Now().Add(Day).Format(time.DateOnly)

	stringDateToSnowflakeTimeFormat := func(inputLayout, date string) *time.Time {
		parsedTime, err := time.Parse(inputLayout, date)
		require.NoError(t, err)

		loc, err := time.LoadLocation("America/Los_Angeles")
		require.NoError(t, err)

		adjustedTime := parsedTime.In(loc)
		return &adjustedTime
	}

	cleanupIntegration := func(t *testing.T, integrationId sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(integrationId).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	err := client.SecurityIntegrations.CreateApiAuthenticationWithClientCredentialsFlow(
		ctx,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo").
			WithOauthAllowedScopes([]sdk.AllowedScope{{"foo"}, {"bar"}}),
	)
	require.NoError(t, err)
	t.Cleanup(cleanupIntegration(t, integrationId))

	cleanupSecret := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
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
		t.Cleanup(cleanupSecret(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithOAuthAuthorizationCodeFlow := func(t *testing.T, integrationId sdk.AccountObjectIdentifier, refreshToken, refreshTokenExpiryTime string, with func(*sdk.CreateWithOAuthAuthorizationCodeFlowSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthAuthorizationCodeFlowSecretRequest(id, refreshToken, refreshTokenExpiryTime, integrationId)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithOAuthAuthorizationCodeFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSecret(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithBasicAuthentication := func(t *testing.T, username, password string, with func(*sdk.CreateWithBasicAuthenticationSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, username, password)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithBasicAuthentication(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSecret(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithGenericString := func(t *testing.T, secretString string, with func(options *sdk.CreateWithGenericStringSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithGenericStringSecretRequest(id, secretString)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithGenericString(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSecret(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
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

	type secretDetails struct {
		Name                        string
		Comment                     string
		SecretType                  string
		Username                    string
		OauthAccessTokenExpiryTime  *time.Time
		OauthRefreshTokenExpiryTime *time.Time
		OauthScopes                 string
		IntegrationName             string
	}

	assertSecretDetails := func(actual *sdk.SecretDetails, expected secretDetails) {
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Comment, actual.Comment.String)
		assert.Equal(t, expected.SecretType, actual.SecretType)
		assert.Equal(t, expected.Username, actual.Username.String)
		assert.Equal(t, expected.OauthAccessTokenExpiryTime, actual.OauthAccessTokenExpiryTime)
		assert.Equal(t, expected.OauthRefreshTokenExpiryTime, actual.OauthRefreshTokenExpiryTime)
		assert.Equal(t, expected.OauthScopes, actual.OauthScopes.String)
		assert.Equal(t, expected.IntegrationName, actual.IntegrationName.String)
	}

	t.Run("Create: secretWithOAuthClientCredentialsFlow", func(t *testing.T) {
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

	// It is possible to create secret without specifying both refresh token properties and scopes
	// Scopes are not being inherited from the security_integration what is tested further
	t.Run("Create: secretWithOAuth - minimal, without token and scopes", func(t *testing.T) {
		secret, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{}, nil)
		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  "OAUTH2",
			OauthScopes:                 "",
			OauthAccessTokenExpiryTime:  nil,
			OauthRefreshTokenExpiryTime: nil,
			IntegrationName:             integrationId.Name(),
		})

		assertSecret(t, secret, id, "OAUTH2", "")
	})

	// regarding the https://docs.snowflake.com/en/sql-reference/sql/create-secret secret should inherit security_integration scopes, but it does not do so
	t.Run("Create: SecretWithOAuthClientCredentialsFlow - No Scopes Specified", func(t *testing.T) {
		_, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{}, nil)
		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		securityIntegrationProperties, _ := client.SecurityIntegrations.Describe(ctx, integrationId)
		assert.Contains(t, securityIntegrationProperties, sdk.SecurityIntegrationProperty{Name: "OAUTH_ALLOWED_SCOPES", Type: "List", Value: "[foo, bar]", Default: "[]"})

		assert.NotEqual(t, details.OauthScopes, securityIntegrationProperties)
	})

	t.Run("Create: SecretWithOAuthAuthorizationCodeFlow", func(t *testing.T) {
		secret, id := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, func(req *sdk.CreateWithOAuthAuthorizationCodeFlowSecretRequest) {
			req.WithComment("a").
				WithIfNotExists(true)
		})

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  "OAUTH2",
			Comment:                     "a",
			OauthRefreshTokenExpiryTime: stringDateToSnowflakeTimeFormat(time.DateOnly, refreshTokenExpiryTime),
			IntegrationName:             integrationId.Name(),
		})

		assertSecret(t, secret, id, "OAUTH2", "a")
	})

	t.Run("Create: WithBasicAuthentication", func(t *testing.T) {
		comment := random.Comment()

		secret, id := createSecretWithBasicAuthentication(t, "foo", "foo", func(req *sdk.CreateWithBasicAuthenticationSecretRequest) {
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

	t.Run("Create: WithBasicAuthentication - Empty Username and Password", func(t *testing.T) {
		comment := random.Comment()
		secret, id := createSecretWithBasicAuthentication(t, "", "", func(req *sdk.CreateWithBasicAuthenticationSecretRequest) {
			req.WithComment(comment).
				WithIfNotExists(true)
		})
		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			Comment:    comment,
			SecretType: "PASSWORD",
			Username:   "",
		})

		assertSecret(t, secret, id, "PASSWORD", comment)
	})

	t.Run("Create: WithGenericString", func(t *testing.T) {
		comment := random.Comment()
		secret, id := createSecretWithGenericString(t, "foo", func(req *sdk.CreateWithGenericStringSecretRequest) {
			req.WithComment(comment).
				WithIfNotExists(true)
		})

		assertSecret(t, secret, id, "GENERIC_STRING", comment)
	})

	t.Run("Create: WithGenericString - empty secret_string", func(t *testing.T) {
		secret, id := createSecretWithGenericString(t, "", nil)
		require.NoError(t, err)

		assertSecret(t, secret, id, "GENERIC_STRING", "")
	})

	t.Run("Alter: SecretWithOAuthClientCredentials", func(t *testing.T) {
		comment := random.Comment()
		_, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{{"foo"}}, nil)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForOAuthClientCredentialsFlow(
						*sdk.NewSetForOAuthClientCredentialsFlowRequest(
							[]sdk.SecurityIntegrationScope{{"foo"}, {"bar"}},
						),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			SecretType:      "OAUTH2",
			Comment:         comment,
			OauthScopes:     "[foo, bar]",
			IntegrationName: integrationId.Name(),
		})

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, details.Comment, sql.NullString{"", false})
	})

	t.Run("Alter: SecretWithOAuthAuthorizationCode", func(t *testing.T) {
		comment := random.Comment()
		alteredRefreshTokenExpiryTime := time.Now().Add(4 * Day).Format(time.DateOnly)

		_, id := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, nil)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForOAuthAuthorizationFlow(
						*sdk.NewSetForOAuthAuthorizationFlowRequest().
							WithOauthRefreshToken("bar").
							WithOauthRefreshTokenExpiryTime(alteredRefreshTokenExpiryTime),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  "OAUTH2",
			Comment:                     comment,
			OauthRefreshTokenExpiryTime: stringDateToSnowflakeTimeFormat(time.DateOnly, alteredRefreshTokenExpiryTime),
			IntegrationName:             integrationId.Name(),
		})

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, details.Comment, sql.NullString{"", false})
	})

	t.Run("Alter: SecretWithBasicAuthorization", func(t *testing.T) {
		comment := random.Comment()

		_, id := createSecretWithBasicAuthentication(t, "foo", "foo", nil)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForBasicAuthentication(
						*sdk.NewSetForBasicAuthenticationRequest().
							WithUsername("bar").
							WithPassword("bar"),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		// Cannot check password since show and describe on secret do not have access to password property
		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			SecretType: "PASSWORD",
			Comment:    comment,
			Username:   "bar",
		})

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, details.Comment, sql.NullString{"", false})
	})

	t.Run("Alter: SecretWithGenericString", func(t *testing.T) {
		comment := random.Comment()
		_, id := createSecretWithGenericString(t, "foo", nil)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForGenericString(
						*sdk.NewSetForGenericStringRequest().
							WithSecretString("bar"),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		/*
			// Cannot check secret_string since show and describe on secret do not have access to secret_string property
				assertSecretDetails(details, secretDetails{
					Name:       id.Name(),
					SecretType: "PASSWORD",
					Comment:    comment,
				})
		*/

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, details.Comment, sql.NullString{"", false})
	})

	t.Run("Drop", func(t *testing.T) {
		_, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{{"foo"}}, nil)

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NotNil(t, secret)
		require.NoError(t, err)

		err = client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id))
		require.NoError(t, err)

		secret, err = client.Secrets.ShowByID(ctx, id)
		require.Nil(t, secret)
		require.Error(t, err)
	})

	t.Run("Drop: non-existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		secretOAuthClientCredentials, _ := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{{"foo"}}, nil)
		secretOAuthAuthorizationCode, _ := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, nil)
		secretBasicAuthentication, _ := createSecretWithBasicAuthentication(t, "foo", "bar", nil)
		secretGenericString, _ := createSecretWithGenericString(t, "foo", nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest())
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)
	})

	t.Run("Show: SecretWithOAuthClientCredentialsFlow with Like", func(t *testing.T) {
		secret1, id1 := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{{"foo"}}, nil)
		secret2, _ := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{{"bar"}}, nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithOAuthAuthorization with Like", func(t *testing.T) {
		secret1, id1 := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo_1", refreshTokenExpiryTime, nil)
		secret2, _ := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo_2", refreshTokenExpiryTime, nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithBasicAuthentication with Like", func(t *testing.T) {
		secret1, id1 := createSecretWithBasicAuthentication(t, "foo_1", "bar_1", nil)
		secret2, _ := createSecretWithBasicAuthentication(t, "foo_2", "bar_2", nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithGenericString with Like", func(t *testing.T) {
		secret1, id1 := createSecretWithGenericString(t, "foo_1", nil)
		secret2, _ := createSecretWithGenericString(t, "foo_2", nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.Pointer(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithOAuthClientCredentialsFlow with In", func(t *testing.T) {
		secret, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{{"foo"}}, nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Account: sdk.Pointer(true)}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Schema: id.SchemaId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)
	})

	t.Run("Show: SecretWithOAuthAuthorizationCodeFlow with In", func(t *testing.T) {
		secret, id := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Account: sdk.Pointer(true)}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Schema: id.SchemaId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)
	})

	t.Run("Show: with In", func(t *testing.T) {
		secretOAuthClientCredentials, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.SecurityIntegrationScope{{"foo"}}, nil)
		secretOAuthAuthorizationCode, _ := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, nil)
		secretBasicAuthentication, _ := createSecretWithBasicAuthentication(t, "foo", "bar", nil)
		secretGenericString, _ := createSecretWithGenericString(t, "foo", nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Account: sdk.Pointer(true)}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Schema: id.SchemaId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)
	})

	t.Run("Show: SecretWithGenericString with In", func(t *testing.T) {
		secret, id := createSecretWithGenericString(t, "foo", nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Account: sdk.Pointer(true)}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Schema: id.SchemaId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)
	})

	t.Run("ShowByID", func(t *testing.T) {
		_, id := createSecretWithGenericString(t, "foo", nil)

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSecret(t, secret, id, "GENERIC_STRING", "")
	})

	t.Run("Describe", func(t *testing.T) {
		_, id := createSecretWithGenericString(t, "foo", func(req *sdk.CreateWithGenericStringSecretRequest) {
			req.WithComment("Lorem ipsum")
		})

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			Comment:    "Lorem ipsum",
			SecretType: "GENERIC_STRING",
		})
	})
}

func TestInt_SecretsShowWithIn(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupSecret := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	createSecretWithGenericString := func(t *testing.T, id sdk.SchemaObjectIdentifier, secretString string) *sdk.Secret {
		t.Helper()
		request := sdk.NewCreateWithGenericStringSecretRequest(id, secretString)
		err := client.Secrets.CreateWithGenericString(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSecret(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret
	}

	t.Run("Show with In - same id in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		secret1 := createSecretWithGenericString(t, id1, "foo")
		secret2 := createSecretWithGenericString(t, id2, "bar")

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Schema: id1.SchemaId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.In{Database: id1.DatabaseId()}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.Contains(t, returnedSecrets, *secret2)
	})
}
