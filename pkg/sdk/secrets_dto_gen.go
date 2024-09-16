package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateWithOAuthClientCredentialsFlowSecretOptions] = new(CreateWithOAuthClientCredentialsFlowSecretRequest)
	_ optionsProvider[CreateWithOAuthAuthorizationCodeFlowSecretOptions] = new(CreateWithOAuthAuthorizationCodeFlowSecretRequest)
	_ optionsProvider[CreateWithBasicAuthenticationSecretOptions]        = new(CreateWithBasicAuthenticationSecretRequest)
	_ optionsProvider[CreateWithGenericStringSecretOptions]              = new(CreateWithGenericStringSecretRequest)
	_ optionsProvider[AlterSecretOptions]                                = new(AlterSecretRequest)
	_ optionsProvider[DropSecretOptions]                                 = new(DropSecretRequest)
	_ optionsProvider[ShowSecretOptions]                                 = new(ShowSecretRequest)
	_ optionsProvider[DescribeSecretOptions]                             = new(DescribeSecretRequest)
)

type CreateWithOAuthClientCredentialsFlowSecretRequest struct {
	OrReplace           *bool
	IfNotExists         *bool
	name                SchemaObjectIdentifier     // required
	SecurityIntegration AccountObjectIdentifier    // required
	OauthScopes         []SecurityIntegrationScope // required
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
	Comment                          *string
	SetForOAuthClientCredentialsFlow *SetForOAuthClientCredentialsFlowRequest
	SetForOAuthAuthorizationFlow     *SetForOAuthAuthorizationFlowRequest
	SetForBasicAuthentication        *SetForBasicAuthenticationRequest
	SetForGenericString              *SetForGenericStringRequest
}

type SetForOAuthClientCredentialsFlowRequest struct {
	OauthScopes []SecurityIntegrationScope // required
}

type SetForOAuthAuthorizationFlowRequest struct {
	OauthRefreshToken           *string
	OauthRefreshTokenExpiryTime *string
}

type SetForBasicAuthenticationRequest struct {
	Username *string
	Password *string
}

type SetForGenericStringRequest struct {
	SecretString *string
}

type SecretUnsetRequest struct {
	Comment *bool
}

type DropSecretRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowSecretRequest struct {
	Like *Like
	In   *In
}

type DescribeSecretRequest struct {
	name SchemaObjectIdentifier // required
}
