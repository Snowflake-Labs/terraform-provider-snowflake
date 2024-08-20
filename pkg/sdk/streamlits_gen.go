package sdk

import (
	"context"
	"database/sql"
)

type Streamlits interface {
	Create(ctx context.Context, request *CreateStreamlitRequest) error
	Alter(ctx context.Context, request *AlterStreamlitRequest) error
	Drop(ctx context.Context, request *DropStreamlitRequest) error
	Show(ctx context.Context, request *ShowStreamlitRequest) ([]Streamlit, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Streamlit, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*StreamlitDetail, error)
}

// CreateStreamlitOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-streamlit.
type CreateStreamlitOptions struct {
	create                     bool                        `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                       `ddl:"keyword" sql:"OR REPLACE"`
	streamlit                  bool                        `ddl:"static" sql:"STREAMLIT"`
	IfNotExists                *bool                       `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier      `ddl:"identifier"`
	RootLocation               string                      `ddl:"parameter,single_quotes" sql:"ROOT_LOCATION"`
	MainFile                   string                      `ddl:"parameter,single_quotes" sql:"MAIN_FILE"`
	QueryWarehouse             *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"QUERY_WAREHOUSE"`
	ExternalAccessIntegrations *ExternalAccessIntegrations `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Title                      *string                     `ddl:"parameter,single_quotes" sql:"TITLE"`
	Comment                    *string                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type ExternalAccessIntegrations struct {
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"list,must_parentheses"`
}

// AlterStreamlitOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-streamlit.
type AlterStreamlitOptions struct {
	alter     bool                    `ddl:"static" sql:"ALTER"`
	streamlit bool                    `ddl:"static" sql:"STREAMLIT"`
	IfExists  *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name      SchemaObjectIdentifier  `ddl:"identifier"`
	Set       *StreamlitSet           `ddl:"keyword" sql:"SET"`
	Unset     *StreamlitUnset         `ddl:"list,no_parentheses" sql:"UNSET"`
	RenameTo  *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
}
type StreamlitSet struct {
	RootLocation               *string                     `ddl:"parameter,single_quotes" sql:"ROOT_LOCATION"`
	MainFile                   *string                     `ddl:"parameter,single_quotes" sql:"MAIN_FILE"`
	QueryWarehouse             *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"QUERY_WAREHOUSE"`
	ExternalAccessIntegrations *ExternalAccessIntegrations `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Comment                    *string                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Title                      *string                     `ddl:"parameter,single_quotes" sql:"TITLE"`
}
type StreamlitUnset struct {
	QueryWarehouse *bool `ddl:"keyword" sql:"QUERY_WAREHOUSE"`
	Comment        *bool `ddl:"keyword" sql:"COMMENT"`
	Title          *bool `ddl:"keyword" sql:"TITLE"`
}

// DropStreamlitOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-streamlit.
type DropStreamlitOptions struct {
	drop      bool                   `ddl:"static" sql:"DROP"`
	streamlit bool                   `ddl:"static" sql:"STREAMLIT"`
	IfExists  *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name      SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowStreamlitOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-streamlits.
type ShowStreamlitOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	streamlits bool       `ddl:"static" sql:"STREAMLITS"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *In        `ddl:"keyword" sql:"IN"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}
type streamlitsRow struct {
	CreatedOn      string         `db:"created_on"`
	Name           string         `db:"name"`
	DatabaseName   string         `db:"database_name"`
	SchemaName     string         `db:"schema_name"`
	Title          sql.NullString `db:"title"`
	Owner          string         `db:"owner"`
	Comment        sql.NullString `db:"comment"`
	QueryWarehouse sql.NullString `db:"query_warehouse"`
	UrlId          string         `db:"url_id"`
	OwnerRoleType  string         `db:"owner_role_type"`
}
type Streamlit struct {
	CreatedOn      string
	Name           string
	DatabaseName   string
	SchemaName     string
	Title          string
	Owner          string
	Comment        string
	QueryWarehouse string
	UrlId          string
	OwnerRoleType  string
}

// DescribeStreamlitOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-streamlit.
type DescribeStreamlitOptions struct {
	describe  bool                   `ddl:"static" sql:"DESCRIBE"`
	streamlit bool                   `ddl:"static" sql:"STREAMLIT"`
	name      SchemaObjectIdentifier `ddl:"identifier"`
}
type streamlitsDetailRow struct {
	Name                       string         `db:"name"`
	Title                      sql.NullString `db:"title"`
	RootLocation               string         `db:"root_location"`
	MainFile                   string         `db:"main_file"`
	QueryWarehouse             sql.NullString `db:"query_warehouse"`
	UrlId                      string         `db:"url_id"`
	DefaultPackages            string         `db:"default_packages"`
	UserPackages               string         `db:"user_packages"`
	ImportUrls                 string         `db:"import_urls"`
	ExternalAccessIntegrations string         `db:"external_access_integrations"`
	ExternalAccessSecrets      string         `db:"external_access_secrets"`
}
type StreamlitDetail struct {
	Name                       string
	Title                      string
	RootLocation               string
	MainFile                   string
	QueryWarehouse             string
	UrlId                      string
	DefaultPackages            string
	UserPackages               []string
	ImportUrls                 []string
	ExternalAccessIntegrations []string
	ExternalAccessSecrets      string
}

func (s *Streamlit) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(s.DatabaseName, s.SchemaName, s.Name)
}
