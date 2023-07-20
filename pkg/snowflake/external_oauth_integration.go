package snowflake

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

type ExternalOauthType string

const (
	Okta         ExternalOauthType = "OKTA"
	Azure        ExternalOauthType = "AZURE"
	PingFederate ExternalOauthType = "PING_FEDERATE"
	Custom       ExternalOauthType = "CUSTOM"
)

type SFUserMappingAttribute string

const (
	LoginName    SFUserMappingAttribute = "LOGIN_NAME"
	EmailAddress SFUserMappingAttribute = "EMAIL_ADDRESS"
)

type AnyRoleMode string

const (
	Disable            AnyRoleMode = "DISABLE"
	Enable             AnyRoleMode = "ENABLE"
	EnableForPrivilege AnyRoleMode = "ENABLE_FOR_PRIVILEGE"
)

type ExternalOauthIntegration3 struct {
	TopLevelIdentifier

	Type                                         string `pos:"parameter" db:"type"`
	TypeOk                                       bool
	Enabled                                      bool `pos:"parameter" db:"enabled"`
	EnabledOk                                    bool
	ExternalOauthType                            ExternalOauthType `pos:"parameter" db:"EXTERNAL_OAUTH_TYPE"`
	ExternalOauthTypeOk                          bool
	ExternalOauthIssuer                          string `pos:"parameter" db:"EXTERNAL_OAUTH_ISSUER"`
	ExternalOauthIssuerOk                        bool
	ExternalOauthTokenUserMappingClaim           []string `pos:"parameter" db:"EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM"`
	ExternalOauthTokenUserMappingClaimOk         bool
	ExternalOauthSnowflakeUserMappingAttribute   SFUserMappingAttribute `pos:"parameter" db:"EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE"`
	ExternalOauthSnowflakeUserMappingAttributeOk bool
	ExternalOauthJwsKeysURL                      []string `pos:"parameter" db:"EXTERNAL_OAUTH_JWS_KEYS_URL"`
	ExternalOauthJwsKeysURLOk                    bool
	ExternalOauthBlockedRolesList                []string `pos:"parameter" db:"EXTERNAL_OAUTH_BLOCKED_ROLES_LIST"`
	ExternalOauthBlockedRolesListOk              bool
	ExternalOauthAllowedRolesList                []string `pos:"parameter" db:"EXTERNAL_OAUTH_ALLOWED_ROLES_LIST"`
	ExternalOauthAllowedRolesListOk              bool
	ExternalOauthRsaPublicKey                    string `pos:"parameter" db:"EXTERNAL_OAUTH_RSA_PUBLIC_KEY"`
	ExternalOauthRsaPublicKeyOk                  bool
	ExternalOauthRsaPublicKey2                   string `pos:"parameter" db:"EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2"`
	ExternalOauthRsaPublicKey2Ok                 bool
	ExternalOauthAudienceList                    []string `pos:"parameter" db:"EXTERNAL_OAUTH_AUDIENCE_LIST"`
	ExternalOauthAudienceListOk                  bool
	ExternalOauthAnyRoleMode                     AnyRoleMode `pos:"parameter" db:"EXTERNAL_OAUTH_ANY_ROLE_MODE"`
	ExternalOauthAnyRoleModeOk                   bool
	ExternalOauthScopeDelimiter                  string `pos:"parameter" db:"EXTERNAL_OAUTH_SCOPE_DELIMITER"`
	ExternalOauthScopeDelimiterOk                bool
	ExternalOauthScopeMappingAttribute           string `pos:"parameter" db:"EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE"`
	ExternalOauthScopeMappingAttributeOk         bool

	Comment   sql.NullString `pos:"parameter" db:"comment"`
	CommentOk bool
}

type ExternalOauthIntegration3Manager struct {
	BaseManager
}

func NewExternalOauthIntegration3Manager() (*ExternalOauthIntegration3Manager, error) {
	sqlBuilder, err := newSQLBuilder(
		"SECURITY INTEGRATION",
		"SECURITY INTEGRATIONS",
		reflect.TypeOf(ExternalOauthIntegration3CreateInput{}),
		reflect.TypeOf(ExternalOauthIntegration3UpdateInput{}),
		reflect.TypeOf(ExternalOauthIntegration3UpdateInput{}),
		reflect.TypeOf(ExternalOauthIntegration3DeleteInput{}),
		reflect.TypeOf(ExternalOauthIntegration3ReadOutput{}),
	)
	if err != nil {
		return nil, err
	}

	return &ExternalOauthIntegration3Manager{
		BaseManager: BaseManager{
			sqlBuilder: *sqlBuilder,
		},
	}, nil
}

type ExternalOauthIntegration3CreateInput struct {
	ExternalOauthIntegration3

	OrReplace     bool `pos:"beforeObjectType" value:"OR REPLACE"`
	OrReplaceOk   bool
	IfNotExists   bool `pos:"afterObjectType" value:"IF NOT EXISTS"`
	IfNotExistsOk bool
}

func (m *ExternalOauthIntegration3Manager) Create(x *ExternalOauthIntegration3CreateInput) (string, error) {
	return m.sqlBuilder.Create(x)
}

type (
	ExternalOauthIntegration3ReadInput  = TopLevelIdentifier
	ExternalOauthIntegration3ReadOutput = ExternalOauthIntegration3
)

func (m *ExternalOauthIntegration3Manager) ReadDescribe(x *ExternalOauthIntegration3ReadInput) (string, error) {
	return m.sqlBuilder.Describe(x)
}

func (m *ExternalOauthIntegration3Manager) ParseDescribe(rows *sql.Rows) (*ExternalOauthIntegration3ReadOutput, error) {
	output := &ExternalOauthIntegration3ReadOutput{}
	err := m.sqlBuilder.ParseDescribe(rows, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (m *ExternalOauthIntegration3Manager) ReadShow(x *ExternalOauthIntegration3ReadInput) (string, error) {
	return m.sqlBuilder.ShowLike(x)
}

func (m *ExternalOauthIntegration3Manager) ParseShow(row *sqlx.Row) (*ExternalOauthIntegration3ReadOutput, error) {
	result := &ExternalOauthIntegration3{}
	if err := row.StructScan(result); err != nil {
		return nil, fmt.Errorf("error scanning result: %w", err)
	}
	return result, nil
}

type ExternalOauthIntegration3UpdateInput struct {
	ExternalOauthIntegration3

	IfExists   bool `pos:"afterObjectType" value:"IF EXISTS"`
	IfExistsOk bool
}

func (m *ExternalOauthIntegration3Manager) Update(x *ExternalOauthIntegration3UpdateInput) (string, error) {
	return m.sqlBuilder.Alter(x)
}

func (m *ExternalOauthIntegration3Manager) Unset(x *ExternalOauthIntegration3UpdateInput) (string, error) {
	return m.sqlBuilder.Unset(x)
}

type ExternalOauthIntegration3DeleteInput struct {
	TopLevelIdentifier

	IfExists   bool `pos:"afterObjectType" value:"IF EXISTS"`
	IfExistsOk bool
}

func (m *ExternalOauthIntegration3Manager) Delete(x *ExternalOauthIntegration3DeleteInput) (string, error) {
	return m.sqlBuilder.Drop(x)
}
