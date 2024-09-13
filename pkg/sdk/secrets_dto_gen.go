package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateWithOAuthClientCredentialsFlowSecretOptions] = new(CreateWithOAuthClientCredentialsFlowSecretRequest)
	_ optionsProvider[CreateWithOAuthAuthorizationCodeFlowSecretOptions] = new(CreateWithOAuthAuthorizationCodeFlowSecretRequest)
	_ optionsProvider[CreateWithBasicAuthenticationSecretOptions]        = new(CreateWithBasicAuthenticationSecretRequest)
	_ optionsProvider[CreateWithGenericStringSecretOptions]              = new(CreateWithGenericStringSecretRequest)
	_ optionsProvider[AlterSecretOptions]                                = new(AlterSecretRequest)
)

type CreateWithOAuthClientCredentialsFlowSecretRequest struct {
	OrReplace           *bool
	IfNotExists         *bool
	name                SchemaObjectIdentifier  // required
	SecurityIntegration AccountObjectIdentifier // required
	OauthScopes         []SecurityIntegrationScope
	Comment             *string
}

type CreateWithOAuthAuthorizationCodeFlowSecretRequest struct {
	OrReplace                   *bool
	IfNotExists                 *bool
	name                        SchemaObjectIdentifier  // required
	OauthRefreshToken           string                  // required
	OauthRefreshTokenExpiryTime string                  // required
	SecurityIntegration         AccountObjectIdentifier // required
	Comment                     *string
}

type CreateWithBasicAuthenticationSecretRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        SchemaObjectIdentifier // required
	Username    string                 // required
	Password    string                 // required
	Comment     *string
}

type CreateWithGenericStringSecretRequest struct {
	OrReplace    *bool
	IfNotExists  *bool
	name         SchemaObjectIdentifier // required
	SecretString string                 // required
	Comment      *string
}

type AlterSecretRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Set      *SecretSetRequest
	Unset    *SecretUnsetRequest
}

type SecretSetRequest struct {
	Comment                     *string
	OAuthScopes                 *OAuthScopesRequest
	OauthRefreshToken           *string
	OauthRefreshTokenExpiryTime *string
	Username                    *string
	Password                    *string
	SecretString                *string
}

type OAuthScopesRequest struct {
	OAuthScopes []SecurityIntegrationScope
}

type SecretUnsetRequest struct {
	UnsetComment *bool
}
