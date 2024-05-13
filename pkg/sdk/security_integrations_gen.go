package sdk

import (
	"context"
	"database/sql"
	"time"
)

type SecurityIntegrations interface {
	CreateSCIM(ctx context.Context, request *CreateSCIMSecurityIntegrationRequest) error
	AlterSCIMIntegration(ctx context.Context, request *AlterSCIMIntegrationSecurityIntegrationRequest) error
	Drop(ctx context.Context, request *DropSecurityIntegrationRequest) error
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]SecurityIntegrationProperty, error)
	Show(ctx context.Context, request *ShowSecurityIntegrationRequest) ([]SecurityIntegration, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*SecurityIntegration, error)
}

// CreateSCIMSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-scim.
type CreateSCIMSecurityIntegrationOptions struct {
	create              bool                     `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	securityIntegration bool                     `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfNotExists         *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                AccountObjectIdentifier  `ddl:"identifier"`
	integrationType     string                   `ddl:"static" sql:"TYPE = SCIM"`
	Enabled             bool                     `ddl:"parameter" sql:"ENABLED"`
	ScimClient          string                   `ddl:"parameter,single_quotes" sql:"SCIM_CLIENT"`
	RunAsRole           string                   `ddl:"parameter,single_quotes" sql:"RUN_AS_ROLE"`
	NetworkPolicy       *AccountObjectIdentifier `ddl:"identifier,equals" sql:"NETWORK_POLICY"`
	SyncPassword        *bool                    `ddl:"parameter" sql:"SYNC_PASSWORD"`
	Comment             *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterSCIMIntegrationSecurityIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-scim.
type AlterSCIMIntegrationSecurityIntegrationOptions struct {
	alter               bool                    `ddl:"static" sql:"ALTER"`
	securityIntegration bool                    `ddl:"static" sql:"SECURITY INTEGRATION"`
	IfExists            *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier `ddl:"identifier"`
	Set                 *SCIMIntegrationSet     `ddl:"keyword" sql:"SET"`
	Unset               *SCIMIntegrationUnset   `ddl:"keyword" sql:"UNSET"`
	SetTag              []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTag            []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

type SCIMIntegrationSet struct {
	Enabled       *bool                    `ddl:"parameter" sql:"ENABLED"`
	NetworkPolicy *AccountObjectIdentifier `ddl:"identifier,equals" sql:"NETWORK_POLICY"`
	SyncPassword  *bool                    `ddl:"parameter" sql:"SYNC_PASSWORD"`
	Comment       *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type SCIMIntegrationUnset struct {
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
	Enabled   bool           `db:"enabled"`
	Comment   sql.NullString `db:"comment"`
	CreatedOn time.Time      `db:"created_on"`
}

type SecurityIntegration struct {
	Name            string
	IntegrationType string
	Enabled         bool
	Comment         string
	CreatedOn       time.Time
}
