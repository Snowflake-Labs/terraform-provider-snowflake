package sdk

import "context"

type Secrets interface {
	CreateWithOAuthClientCredentialsFlow(ctx context.Context, request *CreateWithOAuthClientCredentialsFlowSecretRequest) error
	CreateWithOAuthAuthorizationCodeFlow(ctx context.Context, request *CreateWithOAuthAuthorizationCodeFlowSecretRequest) error
	CreateWithBasicAuthentication(ctx context.Context, request *CreateWithBasicAuthenticationSecretRequest) error
	CreateWithGenericString(ctx context.Context, request *CreateWithGenericStringSecretRequest) error
	Alter(ctx context.Context, request *AlterSecretRequest) error
}

// CreateWithOAuthClientCredentialsFlowSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-secret.
type CreateWithOAuthClientCredentialsFlowSecretOptions struct {
	create              bool                       `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                      `ddl:"keyword" sql:"OR REPLACE"`
	secret              bool                       `ddl:"static" sql:"SECRET"`
	IfNotExists         *bool                      `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                SchemaObjectIdentifier     `ddl:"identifier"`
	Type                string                     `ddl:"static" sql:"TYPE = OAUTH2"`
	SecurityIntegration AccountObjectIdentifier    `ddl:"identifier,equals" sql:"API_AUTHENTICATION"`
	OauthScopes         []SecurityIntegrationScope `ddl:"parameter,parentheses" sql:"OAUTH_SCOPES"`
	Comment             *string                    `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type SecurityIntegrationScope struct {
	Scope string `ddl:"keyword,single_quotes"`
}

// CreateWithOAuthAuthorizationCodeFlowSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-secret.
type CreateWithOAuthAuthorizationCodeFlowSecretOptions struct {
	create                      bool                    `ddl:"static" sql:"CREATE"`
	OrReplace                   *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	secret                      bool                    `ddl:"static" sql:"SECRET"`
	IfNotExists                 *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                        SchemaObjectIdentifier  `ddl:"identifier"`
	Type                        string                  `ddl:"static" sql:"TYPE = OAUTH2"`
	OauthRefreshToken           string                  `ddl:"parameter,single_quotes" sql:"OAUTH_REFRESH_TOKEN"`
	OauthRefreshTokenExpiryTime string                  `ddl:"parameter,single_quotes" sql:"OAUTH_REFRESH_TOKEN_EXPIRY_TIME"`
	SecurityIntegration         AccountObjectIdentifier `ddl:"identifier,equals" sql:"API_AUTHENTICATION"`
	Comment                     *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateWithBasicAuthenticationSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-secret.
type CreateWithBasicAuthenticationSecretOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	secret      bool                   `ddl:"static" sql:"SECRET"`
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	Type        string                 `ddl:"static" sql:"TYPE = PASSWORD"`
	Username    string                 `ddl:"parameter,single_quotes" sql:"USERNAME"`
	Password    string                 `ddl:"parameter,single_quotes" sql:"PASSWORD"`
	Comment     *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateWithGenericStringSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-secret.
type CreateWithGenericStringSecretOptions struct {
	create       bool                   `ddl:"static" sql:"CREATE"`
	OrReplace    *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	secret       bool                   `ddl:"static" sql:"SECRET"`
	IfNotExists  *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
	Type         string                 `ddl:"static" sql:"TYPE = GENERIC_STRING"`
	SecretString string                 `ddl:"parameter,single_quotes" sql:"SECRET_STRING"`
	Comment      *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-secret.
type AlterSecretOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"`
	secret   bool                   `ddl:"static" sql:"SECRET"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
	Set      *SecretSet             `ddl:"keyword" sql:"SET"`
	Unset    *SecretUnset           `ddl:"keyword"`
}
type SecretSet struct {
	Comment                          *string                           `ddl:"parameter,single_quotes" sql:"COMMENT"`
	SetForOAuthClientCredentialsFlow *SetForOAuthClientCredentialsFlow `ddl:"keyword"`
	SetForOAuthAuthorizationFlow     *SetForOAuthAuthorizationFlow     `ddl:"keyword"`
	SetForBasicAuthentication        *SetForBasicAuthentication        `ddl:"keyword"`
	SetForGenericString              *SetForGenericString              `ddl:"keyword"`
}
type SetForOAuthClientCredentialsFlow struct {
	OauthScopes []SecurityIntegrationScope `ddl:"parameter,parentheses" sql:"OAUTH_SCOPES"`
}
type SetForOAuthAuthorizationFlow struct {
	OauthRefreshToken           *string `ddl:"parameter,single_quotes" sql:"OAUTH_REFRESH_TOKEN"`
	OauthRefreshTokenExpiryTime *string `ddl:"parameter,single_quotes" sql:"OAUTH_REFRESH_TOKEN_EXPIRY_TIME"`
}
type SetForBasicAuthentication struct {
	Username *string `ddl:"parameter,single_quotes" sql:"USERNAME"`
	Password *string `ddl:"parameter,single_quotes" sql:"PASSWORD"`
}
type SetForGenericString struct {
	SecretString *string `ddl:"parameter,single_quotes" sql:"SECRET_STRING"`
}
type SecretUnset struct {
	Comment *bool `ddl:"keyword" sql:"SET COMMENT = NULL"`
}
