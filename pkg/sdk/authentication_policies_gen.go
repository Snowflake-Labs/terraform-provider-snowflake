package sdk

import "context"

type AuthenticationPolicies interface {
	Create(ctx context.Context, request *CreateAuthenticationPolicyRequest) error
	Alter(ctx context.Context, request *AlterAuthenticationPolicyRequest) error
	Drop(ctx context.Context, request *DropAuthenticationPolicyRequest) error
	Show(ctx context.Context, request *ShowAuthenticationPolicyRequest) ([]AuthenticationPolicy, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*AuthenticationPolicy, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]AuthenticationPolicyDescription, error)
}

// CreateAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy.
type CreateAuthenticationPolicyOptions struct {
	create                   bool                         `ddl:"static" sql:"CREATE"`
	OrReplace                *bool                        `ddl:"keyword" sql:"OR REPLACE"`
	authenticationPolicy     bool                         `ddl:"static" sql:"AUTHENTICATION POLICY"`
	name                     SchemaObjectIdentifier       `ddl:"identifier"`
	AuthenticationMethods    []AuthenticationMethods      `ddl:"parameter,parentheses" sql:"AUTHENTICATION_METHODS"`
	MfaAuthenticationMethods []MfaAuthenticationMethods   `ddl:"parameter,parentheses" sql:"MFA_AUTHENTICATION_METHODS"`
	MfaEnrollment            *string                      `ddl:"parameter" sql:"MFA_ENROLLMENT"`
	ClientTypes              []ClientTypes                `ddl:"parameter,parentheses" sql:"CLIENT_TYPES"`
	SecurityIntegrations     []SecurityIntegrationsOption `ddl:"parameter,parentheses" sql:"SECURITY_INTEGRATIONS"`
	Comment                  *string                      `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type AuthenticationMethods struct {
	Method string `ddl:"keyword,single_quotes"`
}
type MfaAuthenticationMethods struct {
	Method string `ddl:"keyword,single_quotes"`
}
type ClientTypes struct {
	ClientType string `ddl:"keyword,single_quotes"`
}
type SecurityIntegrationsOption struct {
	Name string `ddl:"keyword,single_quotes"`
}

// AlterAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-authentication-policy.
type AlterAuthenticationPolicyOptions struct {
	alter                bool                       `ddl:"static" sql:"ALTER"`
	authenticationPolicy bool                       `ddl:"static" sql:"AUTHENTICATION POLICY"`
	IfExists             *bool                      `ddl:"keyword" sql:"IF EXISTS"`
	name                 SchemaObjectIdentifier     `ddl:"identifier"`
	Set                  *AuthenticationPolicySet   `ddl:"keyword" sql:"SET"`
	Unset                *AuthenticationPolicyUnset `ddl:"list,no_parentheses" sql:"UNSET"`
	RenameTo             *SchemaObjectIdentifier    `ddl:"identifier" sql:"RENAME TO"`
}
type AuthenticationPolicySet struct {
	AuthenticationMethods    []AuthenticationMethods      `ddl:"parameter,parentheses" sql:"AUTHENTICATION_METHODS"`
	MfaAuthenticationMethods []MfaAuthenticationMethods   `ddl:"parameter,parentheses" sql:"MFA_AUTHENTICATION_METHODS"`
	MfaEnrollment            *string                      `ddl:"parameter,single_quotes" sql:"MFA_ENROLLMENT"`
	ClientTypes              []ClientTypes                `ddl:"parameter,parentheses" sql:"CLIENT_TYPES"`
	SecurityIntegrations     []SecurityIntegrationsOption `ddl:"parameter,parentheses" sql:"SECURITY_INTEGRATIONS"`
	Comment                  *string                      `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type AuthenticationPolicyUnset struct {
	ClientTypes              *bool `ddl:"keyword" sql:"CLIENT_TYPES"`
	AuthenticationMethods    *bool `ddl:"keyword" sql:"AUTHENTICATION_METHODS"`
	SecurityIntegrations     *bool `ddl:"keyword" sql:"SECURITY_INTEGRATIONS"`
	MfaAuthenticationMethods *bool `ddl:"keyword" sql:"MFA_AUTHENTICATION_METHODS"`
	MfaEnrollment            *bool `ddl:"keyword" sql:"MFA_ENROLLMENT"`
	Comment                  *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-authentication-policy.
type DropAuthenticationPolicyOptions struct {
	drop                 bool                   `ddl:"static" sql:"DROP"`
	authenticationPolicy bool                   `ddl:"static" sql:"AUTHENTICATION POLICY"`
	IfExists             *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name                 SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-authentication-policies.
type ShowAuthenticationPolicyOptions struct {
	show                   bool       `ddl:"static" sql:"SHOW"`
	authenticationPolicies bool       `ddl:"static" sql:"AUTHENTICATION POLICIES"`
	Like                   *Like      `ddl:"keyword" sql:"LIKE"`
	In                     *In        `ddl:"keyword" sql:"IN"`
	StartsWith             *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit                  *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}
type showAuthenticationPolicyDBRow struct {
	CreatedOn     string `db:"created_on"`
	Name          string `db:"name"`
	Comment       string `db:"comment"`
	DatabaseName  string `db:"database_name"`
	SchemaName    string `db:"schema_name"`
	Owner         string `db:"owner"`
	OwnerRoleType string `db:"owner_role_type"`
	Options       string `db:"options"`
}
type AuthenticationPolicy struct {
	CreatedOn     string
	Name          string
	Comment       string
	DatabaseName  string
	SchemaName    string
	Owner         string
	OwnerRoleType string
	Options       string
}

// DescribeAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-authentication-policy.
type DescribeAuthenticationPolicyOptions struct {
	describe             bool                   `ddl:"static" sql:"DESCRIBE"`
	authenticationPolicy bool                   `ddl:"static" sql:"AUTHENTICATION POLICY"`
	name                 SchemaObjectIdentifier `ddl:"identifier"`
}
type describeAuthenticationPolicyDBRow struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}
type AuthenticationPolicyDescription struct {
	Name  string
	Value string
}
