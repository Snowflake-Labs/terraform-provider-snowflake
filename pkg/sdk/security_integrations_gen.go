package sdk

import (
	"context"
	"database/sql"
	"time"
)

type SecurityIntegrations interface {
	CreateOauthPartner(ctx context.Context, request *CreateOauthPartnerSecurityIntegrationRequest) error
	CreateOauthCustom(ctx context.Context, request *CreateOauthCustomSecurityIntegrationRequest) error
	CreateSaml2(ctx context.Context, request *CreateSaml2SecurityIntegrationRequest) error
	CreateScim(ctx context.Context, request *CreateScimSecurityIntegrationRequest) error
	AlterOauthPartner(ctx context.Context, request *AlterOauthPartnerSecurityIntegrationRequest) error
	AlterOauthCustom(ctx context.Context, request *AlterOauthCustomSecurityIntegrationRequest) error
	AlterSaml2(ctx context.Context, request *AlterSaml2SecurityIntegrationRequest) error
	AlterScim(ctx context.Context, request *AlterScimSecurityIntegrationRequest) error
	Drop(ctx context.Context, request *DropSecurityIntegrationRequest) error
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]SecurityIntegrationProperty, error)
	Show(ctx context.Context, request *ShowSecurityIntegrationRequest) ([]SecurityIntegration, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*SecurityIntegration, error)
}

// CreateOauthPartnerSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake.
type CreateOauthPartnerSecurityIntegrationOptions struct {
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

type BlockedRolesList struct {
	BlockedRolesList []AccountObjectIdentifier `ddl:"list,must_parentheses"`
}

// CreateOauthCustomSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake.
type CreateOauthCustomSecurityIntegrationOptions struct {
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

// AlterOauthPartnerSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake.
type AlterOauthPartnerSecurityIntegrationOptions struct {
	alter               bool                                   `ddl:"static" sql:"ALTER"`
	securityIntegration bool                                   `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                                  `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier                `ddl:"identifier"`
	SetTags             []TagAssociation                       `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier                     `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *OauthPartnerIntegrationSet   `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *OauthPartnerIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type OauthPartnerIntegrationSet struct {
	Enabled                   *bool                                            `ddl:"parameter" sql:"ENABLED"`
	OauthRedirectUri          *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_REDIRECT_URI"`
	OauthIssueRefreshTokens   *bool                                            `ddl:"parameter" sql:"OAUTH_ISSUE_REFRESH_TOKENS"`
	OauthRefreshTokenValidity *int                                             `ddl:"parameter" sql:"OAUTH_REFRESH_TOKEN_VALIDITY"`
	OauthUseSecondaryRoles    *OauthSecurityIntegrationUseSecondaryRolesOption `ddl:"parameter" sql:"OAUTH_USE_SECONDARY_ROLES"`
	BlockedRolesList          *BlockedRolesList                                `ddl:"parameter,parentheses" sql:"BLOCKED_ROLES_LIST"`
	Comment                   *string                                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type OauthPartnerIntegrationUnset struct {
	Enabled                *bool `ddl:"keyword" sql:"ENABLED"`
	OauthUseSecondaryRoles *bool `ddl:"keyword" sql:"OAUTH_USE_SECONDARY_ROLES"`
}

// AlterOauthCustomSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake.
type AlterOauthCustomSecurityIntegrationOptions struct {
	alter               bool                                  `ddl:"static" sql:"ALTER"`
	securityIntegration bool                                  `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                                 `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier               `ddl:"identifier"`
	SetTags             []TagAssociation                      `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier                    `ddl:"keyword" sql:"UNSET TAG"`
	Set                 *OauthCustomIntegrationSet   `ddl:"list,no_parentheses" sql:"SET"`
	Unset               *OauthCustomIntegrationUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type OauthCustomIntegrationSet struct {
	Enabled                     *bool                                            `ddl:"parameter" sql:"ENABLED"`
	OauthRedirectUri            *string                                          `ddl:"parameter,single_quotes" sql:"OAUTH_REDIRECT_URI"`
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

type OauthCustomIntegrationUnset struct {
	Enabled                  *bool `ddl:"keyword" sql:"ENABLED"`
	OauthUseSecondaryRoles   *bool `ddl:"keyword" sql:"OAUTH_USE_SECONDARY_ROLES"`
	NetworkPolicy            *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
	OauthClientRsaPublicKey  *bool `ddl:"keyword" sql:"OAUTH_CLIENT_RSA_PUBLIC_KEY"`
	OauthClientRsaPublicKey2 *bool `ddl:"keyword" sql:"OAUTH_CLIENT_RSA_PUBLIC_KEY_2"`
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
