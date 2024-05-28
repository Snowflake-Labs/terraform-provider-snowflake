package sdk

import (
	"context"
	"database/sql"
	"time"
)

type SecurityIntegrations interface {
	CreateApiAuthenticationClientCredentialsFlow(ctx context.Context, request *CreateApiAuthenticationClientCredentialsFlowSecurityIntegrationRequest) error
	CreateApiAuthenticationAuthorizationCodeGrantFlow(ctx context.Context, request *CreateApiAuthenticationAuthorizationCodeGrantFlowSecurityIntegrationRequest) error
	CreateApiAuthenticationJwtBearerFlow(ctx context.Context, request *CreateApiAuthenticationJwtBearerFlowSecurityIntegrationRequest) error
	CreateExternalOauth(ctx context.Context, request *CreateExternalOauthSecurityIntegrationRequest) error
	CreateOauthForPartnerApplications(ctx context.Context, request *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) error
	CreateOauthForCustomClients(ctx context.Context, request *CreateOauthForCustomClientsSecurityIntegrationRequest) error
	CreateSaml2(ctx context.Context, request *CreateSaml2SecurityIntegrationRequest) error
	CreateScim(ctx context.Context, request *CreateScimSecurityIntegrationRequest) error
	AlterApiAuthenticationClientCredentialsFlow(ctx context.Context, request *AlterApiAuthenticationClientCredentialsFlowSecurityIntegrationRequest) error
	AlterApiAuthenticationAuthorizationCodeGrantFlow(ctx context.Context, request *AlterApiAuthenticationAuthorizationCodeGrantFlowSecurityIntegrationRequest) error
	AlterApiAuthenticationJwtBearerFlow(ctx context.Context, request *AlterApiAuthenticationJwtBearerFlowSecurityIntegrationRequest) error
	AlterExternalOauth(ctx context.Context, request *AlterExternalOauthSecurityIntegrationRequest) error
	AlterOauthForPartnerApplications(ctx context.Context, request *AlterOauthForPartnerApplicationsSecurityIntegrationRequest) error
	AlterOauthForCustomClients(ctx context.Context, request *AlterOauthForCustomClientsSecurityIntegrationRequest) error
	AlterSaml2(ctx context.Context, request *AlterSaml2SecurityIntegrationRequest) error
	AlterScim(ctx context.Context, request *AlterScimSecurityIntegrationRequest) error
	Drop(ctx context.Context, request *DropSecurityIntegrationRequest) error
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]SecurityIntegrationProperty, error)
	Show(ctx context.Context, request *ShowSecurityIntegrationRequest) ([]SecurityIntegration, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*SecurityIntegration, error)
}

// CreateApiAuthenticationClientCredentialsFlowSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth.
type CreateApiAuthenticationClientCredentialsFlowSecurityIntegrationOptions struct {
	create                    bool                                                             `ddl:"static" sql:"CREATE"`
	OrReplace                 *bool                                                            `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration       bool                                                             `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists               *bool                                                            `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                      AccountObjectIdentifier                                          `ddl:"identifier"`
	integrationType           string                                                           `ddl:"static" sql:"TYPE = API_AUTHENTICATION"`
	authType                  string                                                           `ddl:"static" sql:"AUTH_TYPE = OAUTH2"`
	Enabled                   bool                                                             `ddl:"parameter" sql:"ENABLED"`
	OauthTokenEndpoint        *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_TOKEN_ENDPOINT"`
	OauthClientAuthMethod     *ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption `ddl:"parameter" sql:"OAUTH_CLIENT_AUTH_METHOD"`
	OauthClientId             string                                                           `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_ID"`
	OauthClientSecret         string                                                           `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_SECRET"`
	OauthGrant                *ApiAuthenticationSecurityIntegrationOauthGrantOption            `ddl:"parameter" sql:"OAUTH_GRANT"`
	OauthAccessTokenValidity  *int                                                             `ddl:"parameter" sql:"OAUTH_ACCESS_TOKEN_VALIDITY"`
	OauthRefreshTokenValidity *int                                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	OauthAllowedScopes        []AllowedScope                                                   `ddl:"parameter,parentheses" sql:"OAUTH_ALLOWED_SCOPES"`
	Comment                   *string                                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AllowedScope struct {
	Scope string `ddl:"keyword,single_quotes"`
}

// CreateApiAuthenticationAuthorizationCodeGrantFlowSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth.
type CreateApiAuthenticationAuthorizationCodeGrantFlowSecurityIntegrationOptions struct {
	create                     bool                                                             `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                                                            `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration        bool                                                             `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists                *bool                                                            `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       AccountObjectIdentifier                                          `ddl:"identifier"`
	integrationType            string                                                           `ddl:"static" sql:"TYPE = API_AUTHENTICATION"`
	authType                   string                                                           `ddl:"static" sql:"AUTH_TYPE = OAUTH2"`
	Enabled                    bool                                                             `ddl:"parameter" sql:"ENABLED"`
	OauthAuthorizationEndpoint *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_AUTHORIZATION_ENDPOINT"`
	OauthTokenEndpoint         *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_TOKEN_ENDPOINT"`
	OauthClientAuthMethod      *ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption `ddl:"parameter" sql:"OAUTH_CLIENT_AUTH_METHOD"`
	OauthClientId              string                                                           `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_ID"`
	OauthClientSecret          string                                                           `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_SECRET"`
	OauthGrant                 *ApiAuthenticationSecurityIntegrationOauthGrantOption            `ddl:"parameter" sql:"OAUTH_GRANT"`
	OauthAccessTokenValidity   *int                                                             `ddl:"parameter" sql:"OAUTH_ACCESS_TOKEN_VALIDITY"`
	OauthRefreshTokenValidity  *int                                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	Comment                    *string                                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateApiAuthenticationJwtBearerFlowSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth.
type CreateApiAuthenticationJwtBearerFlowSecurityIntegrationOptions struct {
	create                     bool                                                             `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                                                            `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration        bool                                                             `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists                *bool                                                            `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       AccountObjectIdentifier                                          `ddl:"identifier"`
	integrationType            string                                                           `ddl:"static" sql:"TYPE = API_AUTHENTICATION"`
	authType                   string                                                           `ddl:"static" sql:"AUTH_TYPE = OAUTH2"`
	Enabled                    bool                                                             `ddl:"parameter" sql:"ENABLED"`
	OauthAssertionIssuer       string                                                           `ddl:"parameter,single_quotes" sql:"OAUTH_ASSERTION_ISSUER"`
	OauthAuthorizationEndpoint *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_AUTHORIZATION_ENDPOINT"`
	OauthTokenEndpoint         *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_TOKEN_ENDPOINT"`
	OauthClientAuthMethod      *ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption `ddl:"parameter" sql:"OAUTH_CLIENT_AUTH_METHOD"`
	OauthClientId              string                                                           `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_ID"`
	OauthClientSecret          string                                                           `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_SECRET"`
	OauthGrant                 *ApiAuthenticationSecurityIntegrationOauthGrantOption            `ddl:"parameter" sql:"OAUTH_GRANT"`
	OauthAccessTokenValidity   *int                                                             `ddl:"parameter" sql:"OAUTH_ACCESS_TOKEN_VALIDITY"`
	OauthRefreshTokenValidity  *int                                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	Comment                    *string                                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateExternalOauthSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-external.
type CreateExternalOauthSecurityIntegrationOptions struct {
	create                                     bool                                                                `ddl:"static" sql:"CREATE"`
	OrReplace                                  *bool                                                               `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration                        bool                                                                `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists                                *bool                                                               `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                                       AccountObjectIdentifier                                             `ddl:"identifier"`
	integrationType                            string                                                              `ddl:"static" sql:"TYPE = EXTERNAL_OAUTH"`
	Enabled                                    bool                                                                `ddl:"parameter" sql:"ENABLED"`
	ExternalOauthType                          ExternalOauthSecurityIntegrationTypeOption                          `ddl:"parameter" sql:"EXTERNAL_OAUTH_TYPE"`
	ExternalOauthIssuer                        string                                                              `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_ISSUER"`
	ExternalOauthTokenUserMappingClaim         []TokenUserMappingClaim                                             `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM"`
	ExternalOauthSnowflakeUserMappingAttribute ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE"`
	ExternalOauthJwsKeysUrl                    []JwsKeysUrl                                                        `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_JWS_KEYS_URL"`
	ExternalOauthBlockedRolesList              *BlockedRolesList                                                   `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_BLOCKED_ROLES_LIST"`
	ExternalOauthAllowedRolesList              *AllowedRolesList                                                   `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_ALLOWED_ROLES_LIST"`
	ExternalOauthRsaPublicKey                  *string                                                             `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_RSA_PUBLIC_KEY"`
	ExternalOauthRsaPublicKey2                 *string                                                             `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2"`
	ExternalOauthAudienceList                  *AudienceList                                                       `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_AUDIENCE_LIST"`
	ExternalOauthAnyRoleMode                   *ExternalOauthSecurityIntegrationAnyRoleModeOption                  `ddl:"parameter" sql:"EXTERNAL_OAUTH_ANY_ROLE_MODE"`
	ExternalOauthScopeDelimiter                *string                                                             `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_SCOPE_DELIMITER"`
	ExternalOauthScopeMappingAttribute         *string                                                             `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE"`
	Comment                                    *string                                                             `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AllowedRolesList struct {
	AllowedRolesList []AccountObjectIdentifier `ddl:"list,must_parentheses"`
}

type BlockedRolesList struct {
	BlockedRolesList []AccountObjectIdentifier `ddl:"list,must_parentheses"`
}

type JwsKeysUrl struct {
	JwsKeyUrl string `ddl:"keyword,single_quotes"`
}

type AudienceList struct {
	AudienceList []AudienceListItem `ddl:"list,must_parentheses"`
}

type AudienceListItem struct {
	Item string `ddl:"keyword,single_quotes"`
}

type TokenUserMappingClaim struct {
	Claim string `ddl:"keyword,single_quotes"`
}

// CreateOauthForPartnerApplicationsSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake.
type CreateOauthForPartnerApplicationsSecurityIntegrationOptions struct {
	create                    bool                                             `ddl:"static" sql:"CREATE"`
	OrReplace                 *bool                                            `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration       bool                                             `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists               *bool                                            `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                      AccountObjectIdentifier                          `ddl:"identifier"`
	integrationType           string                                           `ddl:"static" sql:"TYPE = OAUTH"`
	OauthClient               OauthSecurityIntegrationClientOption             `ddl:"parameter" sql:"OAUTH_CLIENT"`
	OauthRedirectUri          *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_REDIRECT_URI"`
	Enabled                   *bool                                            `ddl:"parameter" sql:"ENABLED"`
	OauthIssueRefreshTokens   *bool                                            `ddl:"parameter" sql:"OAUTH_ISSUE_REFRESH_TOKENS"`
	OauthRefreshTokenValidity *int                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	OauthUseSecondaryRoles    *OauthSecurityIntegrationUseSecondaryRolesOption `ddl:"parameter" sql:"OAUTH_USE_SECONDARY_ROLES"`
	BlockedRolesList          *BlockedRolesList                                `ddl:"parameter,parentheses" sql:"BLOCKED_ROLES_LIST"`
	Comment                   *string                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type PreAuthorizedRolesList struct {
	PreAuthorizedRolesList []AccountObjectIdentifier `ddl:"list,must_parentheses"`
}

// CreateOauthForCustomClientsSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake.
type CreateOauthForCustomClientsSecurityIntegrationOptions struct {
	create                      bool                                             `ddl:"static" sql:"CREATE"`
	OrReplace                   *bool                                            `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration         bool                                             `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists                 *bool                                            `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                        AccountObjectIdentifier                          `ddl:"identifier"`
	integrationType             string                                           `ddl:"static" sql:"TYPE = OAUTH"`
	oauthClient                 string                                           `ddl:"static" sql:"OAUTH_CLIENT = CUSTOM"`
	OauthClientType             OauthSecurityIntegrationClientTypeOption         `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_TYPE"`
	OauthRedirectUri            string                                           `ddl:"parameter,single_quotes" sql:"OAUTH_REDIRECT_URI"`
	Enabled                     *bool                                            `ddl:"parameter" sql:"ENABLED"`
	OauthAllowNonTlsRedirectUri *bool                                            `ddl:"parameter" sql:"OAUTH_ALLOW_NON_TLS_REDIRECT_URI"`
	OauthEnforcePkce            *bool                                            `ddl:"parameter" sql:"OAUTH_ENFORCE_PKCE"`
	OauthUseSecondaryRoles      *OauthSecurityIntegrationUseSecondaryRolesOption `ddl:"parameter" sql:"OAUTH_USE_SECONDARY_ROLES"`
	PreAuthorizedRolesList      *PreAuthorizedRolesList                          `ddl:"parameter,parentheses" sql:"PRE_AUTHORIZED_ROLES_LIST"`
	BlockedRolesList            *BlockedRolesList                                `ddl:"parameter,parentheses" sql:"BLOCKED_ROLES_LIST"`
	OauthIssueRefreshTokens     *bool                                            `ddl:"parameter" sql:"OAUTH_ISSUE_REFRESH_TOKENS"`
	OauthRefreshTokenValidity   *int                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	NetworkPolicy               *AccountObjectIdentifier                         `ddl:"identifier,equals" sql:"NETWORK_POLICY"`
	OauthClientRsaPublicKey     *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_RSA_PUBLIC_KEY"`
	OauthClientRsaPublicKey2    *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_RSA_PUBLIC_KEY_2"`
	Comment                     *string                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateSaml2SecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-saml2.
type CreateSaml2SecurityIntegrationOptions struct {
	create                         bool                    `ddl:"static" sql:"CREATE"`
	OrReplace                      *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration            bool                    `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists                    *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                           AccountObjectIdentifier `ddl:"identifier"`
	integrationType                string                  `ddl:"static" sql:"TYPE = SAML2"`
	Enabled                        bool                    `ddl:"parameter" sql:"ENABLED"`
	Saml2Issuer                    string                  `ddl:"parameter,single_quotes" sql:"SAML2_ISSUER"`
	Saml2SsoUrl                    string                  `ddl:"parameter,single_quotes" sql:"SAML2_SSO_URL"`
	Saml2Provider                  string                  `ddl:"parameter,single_quotes" sql:"SAML2_PROVIDER"`
	Saml2X509Cert                  string                  `ddl:"parameter,single_quotes" sql:"SAML2_X509_CERT"`
	AllowedUserDomains             []UserDomain            `ddl:"parameter,parentheses" sql:"ALLOWED_USER_DOMAINS"`
	AllowedEmailPatterns           []EmailPattern          `ddl:"parameter,parentheses" sql:"ALLOWED_EMAIL_PATTERNS"`
	Saml2SpInitiatedLoginPageLabel *string                 `ddl:"parameter,single_quotes" sql:"SAML2_SP_INITIATED_LOGIN_PAGE_LABEL"`
	Saml2EnableSpInitiated         *bool                   `ddl:"parameter" sql:"SAML2_ENABLE_SP_INITIATED"`
	Saml2SnowflakeX509Cert         *string                 `ddl:"parameter,single_quotes" sql:"SAML2_SNOWFLAKE_X509_CERT"`
	Saml2SignRequest               *bool                   `ddl:"parameter" sql:"SAML2_SIGN_REQUEST"`
	Saml2RequestedNameidFormat     *string                 `ddl:"parameter,single_quotes" sql:"SAML2_REQUESTED_NAMEID_FORMAT"`
	Saml2PostLogoutRedirectUrl     *string                 `ddl:"parameter,single_quotes" sql:"SAML2_POST_LOGOUT_REDIRECT_URL"`
	Saml2ForceAuthn                *bool                   `ddl:"parameter" sql:"SAML2_FORCE_AUTHN"`
	Saml2SnowflakeIssuerUrl        *string                 `ddl:"parameter,single_quotes" sql:"SAML2_SNOWFLAKE_ISSUER_URL"`
	Saml2SnowflakeAcsUrl           *string                 `ddl:"parameter,single_quotes" sql:"SAML2_SNOWFLAKE_ACS_URL"`
	Comment                        *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type UserDomain struct {
	Domain string `ddl:"keyword,single_quotes"`
}

type EmailPattern struct {
	Pattern string `ddl:"keyword,single_quotes"`
}

// CreateScimSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-scim.
type CreateScimSecurityIntegrationOptions struct {
	create              bool                                    `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                                   `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration bool                                    `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists         *bool                                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                AccountObjectIdentifier                 `ddl:"identifier"`
	integrationType     string                                  `ddl:"static" sql:"TYPE = SCIM"`
	Enabled             bool                                    `ddl:"parameter" sql:"ENABLED"`
	ScimClient          ScimSecurityIntegrationScimClientOption `ddl:"parameter,single_quotes" sql:"SCIM_CLIENT"`
	RunAsRole           ScimSecurityIntegrationRunAsRoleOption  `ddl:"parameter,single_quotes" sql:"RUN_AS_ROLE"`
	NetworkPolicy       *AccountObjectIdentifier                `ddl:"identifier,equals" sql:"NETWORK_POLICY"`
	SyncPassword        *bool                                   `ddl:"parameter" sql:"SYNC_PASSWORD"`
	Comment             *string                                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterApiAuthenticationClientCredentialsFlowSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-api-auth.
type AlterApiAuthenticationClientCredentialsFlowSecurityIntegrationOptions struct {
	alter               bool                                                    `ddl:"static" sql:"ALTER"`
	securityIntegration bool                                                    `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                                                   `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier                                 `ddl:"identifier"`
	SetTags             []TagAssociation                                        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier                                      `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *ApiAuthenticationClientCredentialsFlowIntegrationSet   `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *ApiAuthenticationClientCredentialsFlowIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type ApiAuthenticationClientCredentialsFlowIntegrationSet struct {
	Enabled                   *bool                                                            `ddl:"parameter" sql:"ENABLED"`
	OauthTokenEndpoint        *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_TOKEN_ENDPOINT"`
	OauthClientAuthMethod     *ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption `ddl:"parameter" sql:"OAUTH_CLIENT_AUTH_METHOD"`
	OauthClientId             *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_ID"`
	OauthClientSecret         *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_SECRET"`
	OauthGrant                *ApiAuthenticationSecurityIntegrationOauthGrantOption            `ddl:"parameter" sql:"OAUTH_GRANT"`
	OauthAccessTokenValidity  *int                                                             `ddl:"parameter" sql:"OAUTH_ACCESS_TOKEN_VALIDITY"`
	OauthRefreshTokenValidity *int                                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	OauthAllowedScopes        []AllowedScope                                                   `ddl:"parameter,parentheses" sql:"OAUTH_ALLOWED_SCOPES"`
	Comment                   *string                                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ApiAuthenticationClientCredentialsFlowIntegrationUnset struct {
	Enabled *bool `ddl:"keyword" sql:"ENABLED"`
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// AlterApiAuthenticationAuthorizationCodeGrantFlowSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-api-auth.
type AlterApiAuthenticationAuthorizationCodeGrantFlowSecurityIntegrationOptions struct {
	alter               bool                                                         `ddl:"static" sql:"ALTER"`
	securityIntegration bool                                                         `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                                                        `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier                                      `ddl:"identifier"`
	SetTags             []TagAssociation                                             `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier                                           `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *ApiAuthenticationAuthorizationCodeGrantFlowIntegrationSet   `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *ApiAuthenticationAuthorizationCodeGrantFlowIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type ApiAuthenticationAuthorizationCodeGrantFlowIntegrationSet struct {
	Enabled                    *bool                                                            `ddl:"parameter" sql:"ENABLED"`
	OauthAuthorizationEndpoint *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_AUTHORIZATION_ENDPOINT"`
	OauthTokenEndpoint         *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_TOKEN_ENDPOINT"`
	OauthClientAuthMethod      *ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption `ddl:"parameter" sql:"OAUTH_CLIENT_AUTH_METHOD"`
	OauthClientId              *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_ID"`
	OauthClientSecret          *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_SECRET"`
	OauthGrant                 *ApiAuthenticationSecurityIntegrationOauthGrantOption            `ddl:"parameter" sql:"OAUTH_GRANT"`
	OauthAccessTokenValidity   *int                                                             `ddl:"parameter" sql:"OAUTH_ACCESS_TOKEN_VALIDITY"`
	OauthRefreshTokenValidity  *int                                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	Comment                    *string                                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ApiAuthenticationAuthorizationCodeGrantFlowIntegrationUnset struct {
	Enabled *bool `ddl:"keyword" sql:"ENABLED"`
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// AlterApiAuthenticationJwtBearerFlowSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-api-auth.
type AlterApiAuthenticationJwtBearerFlowSecurityIntegrationOptions struct {
	alter               bool                                            `ddl:"static" sql:"ALTER"`
	securityIntegration bool                                            `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                                           `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier                         `ddl:"identifier"`
	SetTags             []TagAssociation                                `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier                              `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *ApiAuthenticationJwtBearerFlowIntegrationSet   `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *ApiAuthenticationJwtBearerFlowIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type ApiAuthenticationJwtBearerFlowIntegrationSet struct {
	Enabled                    *bool                                                            `ddl:"parameter" sql:"ENABLED"`
	OauthAuthorizationEndpoint *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_AUTHORIZATION_ENDPOINT"`
	OauthTokenEndpoint         *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_TOKEN_ENDPOINT"`
	OauthClientAuthMethod      *ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption `ddl:"parameter" sql:"OAUTH_CLIENT_AUTH_METHOD"`
	OauthClientId              *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_ID"`
	OauthClientSecret          *string                                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_SECRET"`
	OauthGrant                 *ApiAuthenticationSecurityIntegrationOauthGrantOption            `ddl:"parameter" sql:"OAUTH_GRANT"`
	OauthAccessTokenValidity   *int                                                             `ddl:"parameter" sql:"OAUTH_ACCESS_TOKEN_VALIDITY"`
	OauthRefreshTokenValidity  *int                                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	Comment                    *string                                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ApiAuthenticationJwtBearerFlowIntegrationUnset struct {
	Enabled *bool `ddl:"keyword" sql:"ENABLED"`
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// AlterExternalOauthSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-external.
type AlterExternalOauthSecurityIntegrationOptions struct {
	alter               bool                           `ddl:"static" sql:"ALTER"`
	securityIntegration bool                           `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                          `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier        `ddl:"identifier"`
	SetTags             []TagAssociation               `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier             `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *ExternalOauthIntegrationSet   `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *ExternalOauthIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type ExternalOauthIntegrationSet struct {
	Enabled                                    *bool                                                                `ddl:"parameter" sql:"ENABLED"`
	ExternalOauthType                          *ExternalOauthSecurityIntegrationTypeOption                          `ddl:"parameter" sql:"EXTERNAL_OAUTH_TYPE"`
	ExternalOauthIssuer                        *string                                                              `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_ISSUER"`
	ExternalOauthTokenUserMappingClaim         []TokenUserMappingClaim                                              `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM"`
	ExternalOauthSnowflakeUserMappingAttribute *ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE"`
	ExternalOauthJwsKeysUrl                    []JwsKeysUrl                                                         `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_JWS_KEYS_URL"`
	ExternalOauthBlockedRolesList              *BlockedRolesList                                                    `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_BLOCKED_ROLES_LIST"`
	ExternalOauthAllowedRolesList              *AllowedRolesList                                                    `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_ALLOWED_ROLES_LIST"`
	ExternalOauthRsaPublicKey                  *string                                                              `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_RSA_PUBLIC_KEY"`
	ExternalOauthRsaPublicKey2                 *string                                                              `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2"`
	ExternalOauthAudienceList                  *AudienceList                                                        `ddl:"parameter,parentheses" sql:"EXTERNAL_OAUTH_AUDIENCE_LIST"`
	ExternalOauthAnyRoleMode                   *ExternalOauthSecurityIntegrationAnyRoleModeOption                   `ddl:"parameter" sql:"EXTERNAL_OAUTH_ANY_ROLE_MODE"`
	ExternalOauthScopeDelimiter                *string                                                              `ddl:"parameter,single_quotes" sql:"EXTERNAL_OAUTH_SCOPE_DELIMITER"`
	Comment                                    *string                                                              `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ExternalOauthIntegrationUnset struct {
	Enabled                   *bool `ddl:"keyword" sql:"ENABLED"`
	ExternalOauthAudienceList *bool `ddl:"keyword" sql:"EXTERNAL_OAUTH_AUDIENCE_LIST"`
}

// AlterOauthForPartnerApplicationsSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake.
type AlterOauthForPartnerApplicationsSecurityIntegrationOptions struct {
	alter               bool                                         `ddl:"static" sql:"ALTER"`
	securityIntegration bool                                         `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                                        `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier                      `ddl:"identifier"`
	SetTags             []TagAssociation                             `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier                           `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *OauthForPartnerApplicationsIntegrationSet   `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *OauthForPartnerApplicationsIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type OauthForPartnerApplicationsIntegrationSet struct {
	Enabled                   *bool                                            `ddl:"parameter" sql:"ENABLED"`
	OauthIssueRefreshTokens   *bool                                            `ddl:"parameter" sql:"OAUTH_ISSUE_REFRESH_TOKENS"`
	OauthRedirectUri          *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_REDIRECT_URI"`
	OauthRefreshTokenValidity *int                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	OauthUseSecondaryRoles    *OauthSecurityIntegrationUseSecondaryRolesOption `ddl:"parameter" sql:"OAUTH_USE_SECONDARY_ROLES"`
	BlockedRolesList          *BlockedRolesList                                `ddl:"parameter,parentheses" sql:"BLOCKED_ROLES_LIST"`
	Comment                   *string                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type OauthForPartnerApplicationsIntegrationUnset struct {
	Enabled                *bool `ddl:"keyword" sql:"ENABLED"`
	OauthUseSecondaryRoles *bool `ddl:"keyword" sql:"OAUTH_USE_SECONDARY_ROLES"`
}

// AlterOauthForCustomClientsSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake.
type AlterOauthForCustomClientsSecurityIntegrationOptions struct {
	alter               bool                                   `ddl:"static" sql:"ALTER"`
	securityIntegration bool                                   `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                                  `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier                `ddl:"identifier"`
	SetTags             []TagAssociation                       `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier                     `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *OauthForCustomClientsIntegrationSet   `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *OauthForCustomClientsIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type OauthForCustomClientsIntegrationSet struct {
	Enabled                     *bool                                            `ddl:"parameter" sql:"ENABLED"`
	OauthRedirectUri            *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_REDIRECT_URI"`
	OauthAllowNonTlsRedirectUri *bool                                            `ddl:"parameter" sql:"OAUTH_ALLOW_NON_TLS_REDIRECT_URI"`
	OauthEnforcePkce            *bool                                            `ddl:"parameter" sql:"OAUTH_ENFORCE_PKCE"`
	PreAuthorizedRolesList      *PreAuthorizedRolesList                          `ddl:"parameter,parentheses" sql:"PRE_AUTHORIZED_ROLES_LIST"`
	BlockedRolesList            *BlockedRolesList                                `ddl:"parameter,parentheses" sql:"BLOCKED_ROLES_LIST"`
	OauthIssueRefreshTokens     *bool                                            `ddl:"parameter" sql:"OAUTH_ISSUE_REFRESH_TOKENS"`
	OauthRefreshTokenValidity   *int                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	OauthUseSecondaryRoles      *OauthSecurityIntegrationUseSecondaryRolesOption `ddl:"parameter" sql:"OAUTH_USE_SECONDARY_ROLES"`
	NetworkPolicy               *AccountObjectIdentifier                         `ddl:"identifier,equals" sql:"NETWORK_POLICY"`
	OauthClientRsaPublicKey     *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_RSA_PUBLIC_KEY"`
	OauthClientRsaPublicKey2    *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_CLIENT_RSA_PUBLIC_KEY_2"`
	Comment                     *string                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type OauthForCustomClientsIntegrationUnset struct {
	Enabled                  *bool `ddl:"keyword" sql:"ENABLED"`
	NetworkPolicy            *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
	OauthClientRsaPublicKey  *bool `ddl:"keyword" sql:"OAUTH_CLIENT_RSA_PUBLIC_KEY"`
	OauthClientRsaPublicKey2 *bool `ddl:"keyword" sql:"OAUTH_CLIENT_RSA_PUBLIC_KEY_2"`
	OauthUseSecondaryRoles   *bool `ddl:"keyword" sql:"OAUTH_USE_SECONDARY_ROLES"`
}

// AlterSaml2SecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-saml2.
type AlterSaml2SecurityIntegrationOptions struct {
	alter                           bool                    `ddl:"static" sql:"ALTER"`
	securityIntegration             bool                    `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists                        *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                            AccountObjectIdentifier `ddl:"identifier"`
	SetTags                         []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags                       []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
	Set                             *Saml2IntegrationSet    `ddl:"list,no_parentheses" sql:"SET"`
	Unset                           *Saml2IntegrationUnset  `ddl:"list,no_parentheses" sql:"UNSET"`
	RefreshSaml2SnowflakePrivateKey *bool                   `ddl:"keyword" sql:"REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY"`
}

type Saml2IntegrationSet struct {
	Enabled                        *bool          `ddl:"parameter" sql:"ENABLED"`
	Saml2Issuer                    *string        `ddl:"parameter,single_quotes" sql:"SAML2_ISSUER"`
	Saml2SsoUrl                    *string        `ddl:"parameter,single_quotes" sql:"SAML2_SSO_URL"`
	Saml2Provider                  *string        `ddl:"parameter,single_quotes" sql:"SAML2_PROVIDER"`
	Saml2X509Cert                  *string        `ddl:"parameter,single_quotes" sql:"SAML2_X509_CERT"`
	AllowedUserDomains             []UserDomain   `ddl:"parameter,parentheses" sql:"ALLOWED_USER_DOMAINS"`
	AllowedEmailPatterns           []EmailPattern `ddl:"parameter,parentheses" sql:"ALLOWED_EMAIL_PATTERNS"`
	Saml2SpInitiatedLoginPageLabel *string        `ddl:"parameter,single_quotes" sql:"SAML2_SP_INITIATED_LOGIN_PAGE_LABEL"`
	Saml2EnableSpInitiated         *bool          `ddl:"parameter" sql:"SAML2_ENABLE_SP_INITIATED"`
	Saml2SnowflakeX509Cert         *string        `ddl:"parameter,single_quotes" sql:"SAML2_SNOWFLAKE_X509_CERT"`
	Saml2SignRequest               *bool          `ddl:"parameter" sql:"SAML2_SIGN_REQUEST"`
	Saml2RequestedNameidFormat     *string        `ddl:"parameter,single_quotes" sql:"SAML2_REQUESTED_NAMEID_FORMAT"`
	Saml2PostLogoutRedirectUrl     *string        `ddl:"parameter,single_quotes" sql:"SAML2_POST_LOGOUT_REDIRECT_URL"`
	Saml2ForceAuthn                *bool          `ddl:"parameter" sql:"SAML2_FORCE_AUTHN"`
	Saml2SnowflakeIssuerUrl        *string        `ddl:"parameter,single_quotes" sql:"SAML2_SNOWFLAKE_ISSUER_URL"`
	Saml2SnowflakeAcsUrl           *string        `ddl:"parameter,single_quotes" sql:"SAML2_SNOWFLAKE_ACS_URL"`
	Comment                        *string        `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type Saml2IntegrationUnset struct {
	Saml2ForceAuthn            *bool `ddl:"keyword" sql:"SAML2_FORCE_AUTHN"`
	Saml2RequestedNameidFormat *bool `ddl:"keyword" sql:"SAML2_REQUESTED_NAMEID_FORMAT"`
	Saml2PostLogoutRedirectUrl *bool `ddl:"keyword" sql:"SAML2_POST_LOGOUT_REDIRECT_URL"`
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
}

// AlterScimSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-scim.
type AlterScimSecurityIntegrationOptions struct {
	alter               bool                    `ddl:"static" sql:"ALTER"`
	securityIntegration bool                    `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier `ddl:"identifier"`
	SetTags             []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *ScimIntegrationSet     `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *ScimIntegrationUnset   `ddl:"list,no_parentheses" sql:"UNSET"`
}

type ScimIntegrationSet struct {
	Enabled       *bool                    `ddl:"parameter" sql:"ENABLED"`
	NetworkPolicy *AccountObjectIdentifier `ddl:"identifier,equals" sql:"NETWORK_POLICY"`
	SyncPassword  *bool                    `ddl:"parameter" sql:"SYNC_PASSWORD"`
	Comment       *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ScimIntegrationUnset struct {
	Enabled       *bool `ddl:"keyword" sql:"ENABLED"`
	NetworkPolicy *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
	SyncPassword  *bool `ddl:"keyword" sql:"SYNC_PASSWORD"`
	Comment       *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-integration.
type DropSecurityIntegrationOptions struct {
	drop                bool                    `ddl:"static" sql:"DROP"`
	securityIntegration bool                    `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier `ddl:"identifier"`
}

// DescribeSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-integration.
type DescribeSecurityIntegrationOptions struct {
	describe            bool                    `ddl:"static" sql:"DESCRIBE"`
	securityIntegration bool                    `ddl:"static" sql:"SECURITY INTEGRATION"`
	name                AccountObjectIdentifier `ddl:"identifier"`
}

type securityIntegrationDescRow struct {
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

type SecurityIntegrationProperty struct {
	Name    string
	Type    string
	Value   string
	Default string
}

// ShowSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-integrations.
type ShowSecurityIntegrationOptions struct {
	show                 bool  `ddl:"static" sql:"SHOW"`
	securityIntegrations bool  `ddl:"static" sql:"SECURITY INTEGRATIONS"`
	Like                 *Like `ddl:"keyword" sql:"LIKE"`
}

type securityIntegrationShowRow struct {
	Name      string         `db:"name"`
	Type      string         `db:"type"`
	Category  string         `db:"category"`
	Enabled   bool           `db:"enabled"`
	Comment   sql.NullString `db:"comment"`
	CreatedOn time.Time      `db:"created_on"`
}

type SecurityIntegration struct {
	Name            string
	IntegrationType string
	Category        string
	Enabled         bool
	Comment         string
	CreatedOn       time.Time
}
