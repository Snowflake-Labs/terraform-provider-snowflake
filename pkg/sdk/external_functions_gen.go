package sdk

import (
	"context"
	"database/sql"
)

type ExternalFunctions interface {
	Create(ctx context.Context, request *CreateExternalFunctionRequest) error
	Alter(ctx context.Context, request *AlterExternalFunctionRequest) error
	Show(ctx context.Context, request *ShowExternalFunctionRequest) ([]ExternalFunction, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier, arguments []DataType) (*ExternalFunction, error)
	Describe(ctx context.Context, request *DescribeExternalFunctionRequest) ([]ExternalFunctionProperty, error)
}

// CreateExternalFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-external-function.
type CreateExternalFunctionOptions struct {
	create                bool                            `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                           `ddl:"keyword" sql:"OR REPLACE"`
	Secure                *bool                           `ddl:"keyword" sql:"SECURE"`
	externalFunction      bool                            `ddl:"static" sql:"EXTERNAL FUNCTION"`
	name                  SchemaObjectIdentifier          `ddl:"identifier"`
	Arguments             []ExternalFunctionArgument      `ddl:"list,must_parentheses"`
	ResultDataType        DataType                        `ddl:"parameter,no_equals" sql:"RETURNS"`
	ReturnNullValues      *ReturnNullValues               `ddl:"keyword"`
	NullInputBehavior     *NullInputBehavior              `ddl:"keyword"`
	ReturnResultsBehavior *ReturnResultsBehavior          `ddl:"keyword"`
	Comment               *string                         `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ApiIntegration        *AccountObjectIdentifier        `ddl:"identifier" sql:"API_INTEGRATION ="`
	Headers               []ExternalFunctionHeader        `ddl:"parameter,parentheses" sql:"HEADERS"`
	ContextHeaders        []ExternalFunctionContextHeader `ddl:"parameter,parentheses" sql:"CONTEXT_HEADERS"`
	MaxBatchRows          *int                            `ddl:"parameter" sql:"MAX_BATCH_ROWS"`
	Compression           *string                         `ddl:"parameter" sql:"COMPRESSION"`
	RequestTranslator     *SchemaObjectIdentifier         `ddl:"identifier" sql:"REQUEST_TRANSLATOR ="`
	ResponseTranslator    *SchemaObjectIdentifier         `ddl:"identifier" sql:"RESPONSE_TRANSLATOR ="`
	As                    string                          `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type ExternalFunctionArgument struct {
	ArgName     string   `ddl:"keyword,no_quotes"`
	ArgDataType DataType `ddl:"keyword,no_quotes"`
}

type ExternalFunctionHeader struct {
	Name  string `ddl:"keyword,single_quotes"`
	Value string `ddl:"parameter,single_quotes"`
}

type ExternalFunctionContextHeader struct {
	ContextFunction string `ddl:"keyword,no_quotes"`
}

// AlterExternalFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-function.
type AlterExternalFunctionOptions struct {
	alter             bool                   `ddl:"static" sql:"ALTER"`
	function          bool                   `ddl:"static" sql:"FUNCTION"`
	IfExists          *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	ArgumentDataTypes []DataType             `ddl:"keyword,must_parentheses"`
	Set               *ExternalFunctionSet   `ddl:"keyword" sql:"SET"`
	Unset             *ExternalFunctionUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type ExternalFunctionSet struct {
	ApiIntegration     *AccountObjectIdentifier        `ddl:"identifier" sql:"API_INTEGRATION ="`
	Headers            []ExternalFunctionHeader        `ddl:"parameter,parentheses" sql:"HEADERS"`
	ContextHeaders     []ExternalFunctionContextHeader `ddl:"parameter,parentheses" sql:"CONTEXT_HEADERS"`
	MaxBatchRows       *int                            `ddl:"parameter" sql:"MAX_BATCH_ROWS"`
	Compression        *string                         `ddl:"parameter" sql:"COMPRESSION"`
	RequestTranslator  *SchemaObjectIdentifier         `ddl:"identifier" sql:"REQUEST_TRANSLATOR ="`
	ResponseTranslator *SchemaObjectIdentifier         `ddl:"identifier" sql:"RESPONSE_TRANSLATOR ="`
}

type ExternalFunctionUnset struct {
	Comment            *bool `ddl:"keyword" sql:"COMMENT"`
	Headers            *bool `ddl:"keyword" sql:"HEADERS"`
	ContextHeaders     *bool `ddl:"keyword" sql:"CONTEXT_HEADERS"`
	MaxBatchRows       *bool `ddl:"keyword" sql:"MAX_BATCH_ROWS"`
	Compression        *bool `ddl:"keyword" sql:"COMPRESSION"`
	Secure             *bool `ddl:"keyword" sql:"SECURE"`
	RequestTranslator  *bool `ddl:"keyword" sql:"REQUEST_TRANSLATOR"`
	ResponseTranslator *bool `ddl:"keyword" sql:"RESPONSE_TRANSLATOR"`
}

// ShowExternalFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-external-functions.
type ShowExternalFunctionOptions struct {
	show              bool  `ddl:"static" sql:"SHOW"`
	externalFunctions bool  `ddl:"static" sql:"EXTERNAL FUNCTIONS"`
	Like              *Like `ddl:"keyword" sql:"LIKE"`
}

type externalFunctionRow struct {
	CreatedOn          string         `db:"created_on"`
	Name               string         `db:"name"`
	SchemaName         string         `db:"schema_name"`
	IsBuiltin          string         `db:"is_builtin"`
	IsAggregate        string         `db:"is_aggregate"`
	IsAnsi             string         `db:"is_ansi"`
	MinNumArguments    int            `db:"min_num_arguments"`
	MaxNumArguments    int            `db:"max_num_arguments"`
	Arguments          string         `db:"arguments"`
	Description        string         `db:"description"`
	CatalogName        string         `db:"catalog_name"`
	IsTableFunction    string         `db:"is_table_function"`
	ValidForClustering string         `db:"valid_for_clustering"`
	IsSecure           sql.NullString `db:"is_secure"`
	IsExternalFunction string         `db:"is_external_function"`
	Language           string         `db:"language"`
	IsMemoizable       sql.NullString `db:"is_memoizable"`
	IsDataMetric       sql.NullString `db:"is_data_metric"`
}

type ExternalFunction struct {
	CreatedOn          string
	Name               string
	SchemaName         string
	IsBuiltin          bool
	IsAggregate        bool
	IsAnsi             bool
	MinNumArguments    int
	MaxNumArguments    int
	Arguments          string
	Description        string
	CatalogName        string
	IsTableFunction    bool
	ValidForClustering bool
	IsSecure           bool
	IsExternalFunction bool
	Language           string
	IsMemoizable       bool
	IsDataMetric       bool
}

// DescribeExternalFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-function.
type DescribeExternalFunctionOptions struct {
	describe          bool                   `ddl:"static" sql:"DESCRIBE"`
	function          bool                   `ddl:"static" sql:"FUNCTION"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	ArgumentDataTypes []DataType             `ddl:"keyword,must_parentheses"`
}

type externalFunctionPropertyRow struct {
	Property string `db:"property"`
	Value    string `db:"value"`
}

type ExternalFunctionProperty struct {
	Property string
	Value    string
}
