package testint

import (
	"database/sql"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInt_Secrets(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupSecret := func(id sdk.SchemaObjectIdentifier) {
		t.Cleanup(func() {
			err := client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))
			assert.NoError(t, err)
		})
	}

	assertSecret := func(t *testing.T, s *sdk.Secret, id sdk.SchemaObjectIdentifier, secretType, comment string) {
		t.Helper()
		assert.NotEmpty(t, s.CreatedOn)
		assert.NotEmpty(t, s.DatabaseName)
		assert.NotEmpty(t, s.SchemaName)
		assert.NotEmpty(t, s.OwnerRoleType)
		assert.NotEmpty(t, s.Owner)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, secretType, s.SecretType)
		assert.Equal(t, comment, s.Comment)
	}

	type secretDetails struct {
		Name                        string
		Comment                     sql.NullString
		SecretType                  string
		Username                    sql.NullString
		OauthAccessTokenExpiryTime  sql.NullString
		OauthRefreshTokenExpiryTime sql.NullString
		OauthScopes                 sql.NullString
		IntegrationName             sql.NullString
	}

	assertBasicAuth := func(details *sdk.SecretDetails, d secretDetails) {
		assert.Equal(t, d.Name, details.Name)
		assert.Equal(t, d.Comment, details.Comment)
		assert.Equal(t, d.SecretType, details.SecretType)
		assert.Equal(t, d.Username, details.Username)
		assert.Equal(t, d.OauthAccessTokenExpiryTime, details.OauthAccessTokenExpiryTime)
		assert.Equal(t, d.OauthRefreshTokenExpiryTime, details.OauthRefreshTokenExpiryTime)
		assert.Equal(t, d.OauthScopes, details.OauthScopes)
		assert.Equal(t, d.IntegrationName, details.IntegrationName)
	}
	_ = assertBasicAuth

	/*
		createOAuthClientCredentialsFlowSecret := func(id sdk.SchemaObjectIdentifier, si sdk.CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest) {
			//TODO
		}

		createBasicAuthenticationSecret := func(t *testing.T, id sdk.SchemaObjectIdentifier, username, password string) *sdk.Secret {
			t.Helper()

			err := client.Secrets.CreateWithBasicAuthentication(ctx, sdk.NewCreateWithBasicAuthenticationSecretRequest(id, username, password))
			require.NoError(t, err)
			cleanupSecret(id)

			secret, err := client.Secrets.ShowByID(ctx, id)
			return secret
		}
	*/

	t.Run("Create secret with OAuth Client Credentials Flow", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("CreateWithOAuthAuthorizationCodeFlow", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Create With Basic Authentication", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, "foo", "foo").WithComment(comment).WithIfNotExists(true)
		// todo: secrets = nil
		err := client.Secrets.CreateWithBasicAuthentication(ctx, request)
		require.NoError(t, err)
		cleanupSecret(id)

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSecret(t, secret, id, "PASSWORD", comment)
		/*
			details, err := client.Secrets.Describe(ctx, id)

			assertBasicAuth(details, secretDetails{
				Name:                        id.FullyQualifiedName(),
				Comment:                     sql.NullString{String: comment},
				SecretType:                  "PASSWORD",
				Username:                    sql.NullString{String: "foo"},
				OauthAccessTokenExpiryTime:  sql.NullString{},
				OauthRefreshTokenExpiryTime: sql.NullString{},
				OauthScopes:                 sql.NullString{},
				IntegrationName:             sql.NullString{},
			})
		*/
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
