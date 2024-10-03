package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Secrets interface {
	CreateWithOAuthClientCredentialsFlow(ctx context.Context, request *CreateWithOAuthClientCredentialsFlowSecretRequest) error
	CreateWithOAuthAuthorizationCodeFlow(ctx context.Context, request *CreateWithOAuthAuthorizationCodeFlowSecretRequest) error
	CreateWithBasicAuthentication(ctx context.Context, request *CreateWithBasicAuthenticationSecretRequest) error
	CreateWithGenericString(ctx context.Context, request *CreateWithGenericStringSecretRequest) error
	Alter(ctx context.Context, request *AlterSecretRequest) error
	Drop(ctx context.Context, request *DropSecretRequest) error
	Show(ctx context.Context, request *ShowSecretRequest) ([]Secret, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Secret, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*SecretDetails, error)
}

// CreateWithOAuthClientCredentialsFlowSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-secret.
type CreateWithOAuthClientCredentialsFlowSecretOptions struct {
	create         bool                    `ddl:"static" sql:"CREATE"`
	OrReplace      *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	secret         bool                    `ddl:"static" sql:"SECRET"`
	IfNotExists    *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name           SchemaObjectIdentifier  `ddl:"identifier"`
	secretType     string                  `ddl:"static" sql:"TYPE = OAUTH2"`
	ApiIntegration AccountObjectIdentifier `ddl:"identifier,equals" sql:"API_AUTHENTICATION"`
	OauthScopes    *OauthScopesList        `ddl:"parameter,parentheses" sql:"OAUTH_SCOPES"`
	Comment        *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type ApiIntegrationScope struct {
	Scope string `ddl:"keyword,single_quotes"`
}
type OauthScopesList struct {
	OauthScopesList []ApiIntegrationScope `ddl:"list,must_parentheses"`
}

// CreateWithOAuthAuthorizationCodeFlowSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-secret.
type CreateWithOAuthAuthorizationCodeFlowSecretOptions struct {
	create                      bool                    `ddl:"static" sql:"CREATE"`
	OrReplace                   *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	secret                      bool                    `ddl:"static" sql:"SECRET"`
	IfNotExists                 *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                        SchemaObjectIdentifier  `ddl:"identifier"`
	secretType                  string                  `ddl:"static" sql:"TYPE = OAUTH2"`
	OauthRefreshToken           string                  `ddl:"parameter,single_quotes" sql:"OAUTH_REFRESH_TOKEN"`
	OauthRefreshTokenExpiryTime string                  `ddl:"parameter,single_quotes" sql:"OAUTH_REFRESH_TOKEN_EXPIRY_TIME"`
	ApiIntegration              AccountObjectIdentifier `ddl:"identifier,equals" sql:"API_AUTHENTICATION"`
	Comment                     *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateWithBasicAuthenticationSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-secret.
type CreateWithBasicAuthenticationSecretOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	secret      bool                   `ddl:"static" sql:"SECRET"`
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	secretType  string                 `ddl:"static" sql:"TYPE = PASSWORD"`
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
	secretType   string                 `ddl:"static" sql:"TYPE = GENERIC_STRING"`
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
	OauthScopes *OauthScopesList `ddl:"parameter,parentheses" sql:"OAUTH_SCOPES"`
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

// DropSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-secret.
type DropSecretOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	secret   bool                   `ddl:"static" sql:"SECRET"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-secrets.
type ShowSecretOptions struct {
	show    bool        `ddl:"static" sql:"SHOW"`
	secrets bool        `ddl:"static" sql:"SECRETS"`
	Like    *Like       `ddl:"keyword" sql:"LIKE"`
	In      *ExtendedIn `ddl:"keyword" sql:"IN"`
}
type secretDBRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	SchemaName    string         `db:"schema_name"`
	DatabaseName  string         `db:"database_name"`
	Owner         string         `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	SecretType    string         `db:"secret_type"`
	OauthScopes   sql.NullString `db:"oauth_scopes"`
	OwnerRoleType string         `db:"owner_role_type"`
}
type Secret struct {
	CreatedOn     time.Time
	Name          string
	SchemaName    string
	DatabaseName  string
	Owner         string
	Comment       *string
	SecretType    string
	OauthScopes   []string
	OwnerRoleType string
}

func (s *Secret) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(s.DatabaseName, s.SchemaName, s.Name)
}

func (s *Secret) ObjectType() ObjectType {
	return ObjectTypeSecret
}

// DescribeSecretOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-secret.
type DescribeSecretOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	secret   bool                   `ddl:"static" sql:"SECRET"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}
type secretDetailsDBRow struct {
	CreatedOn                   time.Time      `db:"created_on"`
	Name                        string         `db:"name"`
	SchemaName                  string         `db:"schema_name"`
	DatabaseName                string         `db:"database_name"`
	Owner                       string         `db:"owner"`
	Comment                     sql.NullString `db:"comment"`
	SecretType                  string         `db:"secret_type"`
	Username                    sql.NullString `db:"username"`
	OauthAccessTokenExpiryTime  *time.Time     `db:"oauth_access_token_expiry_time"`
	OauthRefreshTokenExpiryTime *time.Time     `db:"oauth_refresh_token_expiry_time"`
	OauthScopes                 sql.NullString `db:"oauth_scopes"`
	IntegrationName             sql.NullString `db:"integration_name"`
}
type SecretDetails struct {
	CreatedOn                   time.Time
	Name                        string
	SchemaName                  string
	DatabaseName                string
	Owner                       string
	Comment                     *string
	SecretType                  string
	Username                    *string
	OauthAccessTokenExpiryTime  *time.Time
	OauthRefreshTokenExpiryTime *time.Time
	OauthScopes                 []string
	IntegrationName             *string
}
