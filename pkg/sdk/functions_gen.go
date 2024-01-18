package sdk

import (
	"context"
	"database/sql"
)

type Functions interface {
	CreateForJava(ctx context.Context, request *CreateForJavaFunctionRequest) error
	CreateForJavascript(ctx context.Context, request *CreateForJavascriptFunctionRequest) error
	CreateForPython(ctx context.Context, request *CreateForPythonFunctionRequest) error
	CreateForScala(ctx context.Context, request *CreateForScalaFunctionRequest) error
	CreateForSQL(ctx context.Context, request *CreateForSQLFunctionRequest) error
	Alter(ctx context.Context, request *AlterFunctionRequest) error
	Drop(ctx context.Context, request *DropFunctionRequest) error
	Show(ctx context.Context, request *ShowFunctionRequest) ([]Function, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Function, error)
	Describe(ctx context.Context, request *DescribeFunctionRequest) ([]FunctionDetail, error)
}

// CreateForJavaFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#java-handler.
type CreateForJavaFunctionOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Temporary                  *bool                     `ddl:"keyword" sql:"TEMPORARY"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	function                   bool                      `ddl:"static" sql:"FUNCTION"`
	IfNotExists                *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []FunctionArgument        `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    FunctionReturns           `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues           *ReturnNullValues         `ddl:"keyword"`
	languageJava               bool                      `ddl:"static" sql:"LANGUAGE JAVA"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	ReturnResultsBehavior      *ReturnResultsBehavior    `ddl:"keyword"`
	RuntimeVersion             *string                   `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports                    []FunctionImport          `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages                   []FunctionPackage         `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []Secret                  `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                   `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	FunctionDefinition         *string                   `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type FunctionArgument struct {
	ArgName      string   `ddl:"keyword,no_quotes"`
	ArgDataType  DataType `ddl:"keyword,no_quotes"`
	DefaultValue *string  `ddl:"parameter,no_equals" sql:"DEFAULT"`
}

type FunctionReturns struct {
	ResultDataType *FunctionReturnsResultDataType `ddl:"keyword"`
	Table          *FunctionReturnsTable          `ddl:"keyword" sql:"TABLE"`
}

type FunctionReturnsResultDataType struct {
	ResultDataType DataType `ddl:"keyword,no_quotes"`
}

type FunctionReturnsTable struct {
	Columns []FunctionColumn `ddl:"parameter,parentheses,no_equals"`
}

type FunctionColumn struct {
	ColumnName     string   `ddl:"keyword,no_quotes"`
	ColumnDataType DataType `ddl:"keyword,no_quotes"`
}

type FunctionImport struct {
	Import string `ddl:"keyword,single_quotes"`
}

type FunctionPackage struct {
	Package string `ddl:"keyword,single_quotes"`
}

// CreateForJavascriptFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#javascript-handler.
type CreateForJavascriptFunctionOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	Secure                *bool                  `ddl:"keyword" sql:"SECURE"`
	function              bool                   `ddl:"static" sql:"FUNCTION"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	Arguments             []FunctionArgument     `ddl:"list,must_parentheses"`
	CopyGrants            *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Returns               FunctionReturns        `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues      *ReturnNullValues      `ddl:"keyword"`
	languageJavascript    bool                   `ddl:"static" sql:"LANGUAGE JAVASCRIPT"`
	NullInputBehavior     *NullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior *ReturnResultsBehavior `ddl:"keyword"`
	Comment               *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	FunctionDefinition    string                 `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateForPythonFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#python-handler.
type CreateForPythonFunctionOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Temporary                  *bool                     `ddl:"keyword" sql:"TEMPORARY"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	function                   bool                      `ddl:"static" sql:"FUNCTION"`
	IfNotExists                *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []FunctionArgument        `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    FunctionReturns           `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues           *ReturnNullValues         `ddl:"keyword"`
	languagePython             bool                      `ddl:"static" sql:"LANGUAGE PYTHON"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	ReturnResultsBehavior      *ReturnResultsBehavior    `ddl:"keyword"`
	RuntimeVersion             string                    `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports                    []FunctionImport          `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages                   []FunctionPackage         `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []Secret                  `ddl:"parameter,parentheses" sql:"SECRETS"`
	FunctionDefinition         *string                   `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateForScalaFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#scala-handler.
type CreateForScalaFunctionOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	Secure                *bool                  `ddl:"keyword" sql:"SECURE"`
	function              bool                   `ddl:"static" sql:"FUNCTION"`
	IfNotExists           *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	Arguments             []FunctionArgument     `ddl:"list,must_parentheses"`
	CopyGrants            *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	ResultDataType        DataType               `ddl:"parameter,no_equals" sql:"RETURNS"`
	ReturnNullValues      *ReturnNullValues      `ddl:"keyword"`
	languageScala         bool                   `ddl:"static" sql:"LANGUAGE SCALA"`
	NullInputBehavior     *NullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior *ReturnResultsBehavior `ddl:"keyword"`
	RuntimeVersion        *string                `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment               *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports               []FunctionImport       `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages              []FunctionPackage      `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler               string                 `ddl:"parameter,single_quotes" sql:"HANDLER"`
	TargetPath            *string                `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	FunctionDefinition    *string                `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateForSQLFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#sql-handler.
type CreateForSQLFunctionOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	Secure                *bool                  `ddl:"keyword" sql:"SECURE"`
	function              bool                   `ddl:"static" sql:"FUNCTION"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	Arguments             []FunctionArgument     `ddl:"list,must_parentheses"`
	CopyGrants            *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Returns               FunctionReturns        `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues      *ReturnNullValues      `ddl:"keyword"`
	ReturnResultsBehavior *ReturnResultsBehavior `ddl:"keyword"`
	Memoizable            *bool                  `ddl:"keyword" sql:"MEMOIZABLE"`
	Comment               *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	FunctionDefinition    string                 `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// AlterFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-function.
type AlterFunctionOptions struct {
	alter             bool                    `ddl:"static" sql:"ALTER"`
	function          bool                    `ddl:"static" sql:"FUNCTION"`
	IfExists          *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name              SchemaObjectIdentifier  `ddl:"identifier"`
	ArgumentDataTypes []DataType              `ddl:"keyword,must_parentheses"`
	RenameTo          *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SetComment        *string                 `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	SetLogLevel       *string                 `ddl:"parameter,single_quotes" sql:"SET LOG_LEVEL"`
	SetTraceLevel     *string                 `ddl:"parameter,single_quotes" sql:"SET TRACE_LEVEL"`
	SetSecure         *bool                   `ddl:"keyword" sql:"SET SECURE"`
	UnsetSecure       *bool                   `ddl:"keyword" sql:"UNSET SECURE"`
	UnsetLogLevel     *bool                   `ddl:"keyword" sql:"UNSET LOG_LEVEL"`
	UnsetTraceLevel   *bool                   `ddl:"keyword" sql:"UNSET TRACE_LEVEL"`
	UnsetComment      *bool                   `ddl:"keyword" sql:"UNSET COMMENT"`
	SetTags           []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags         []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

// DropFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-function.
type DropFunctionOptions struct {
	drop              bool                   `ddl:"static" sql:"DROP"`
	function          bool                   `ddl:"static" sql:"FUNCTION"`
	IfExists          *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	ArgumentDataTypes []DataType             `ddl:"keyword,must_parentheses"`
}

// ShowFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-user-functions.
type ShowFunctionOptions struct {
	show          bool  `ddl:"static" sql:"SHOW"`
	userFunctions bool  `ddl:"static" sql:"USER FUNCTIONS"`
	Like          *Like `ddl:"keyword" sql:"LIKE"`
	In            *In   `ddl:"keyword" sql:"IN"`
}

type functionRow struct {
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
}

type Function struct {
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
}

// DescribeFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-function.
type DescribeFunctionOptions struct {
	describe          bool                   `ddl:"static" sql:"DESCRIBE"`
	function          bool                   `ddl:"static" sql:"FUNCTION"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	ArgumentDataTypes []DataType             `ddl:"keyword,must_parentheses"`
}

type functionDetailRow struct {
	Property string `db:"property"`
	Value    string `db:"value"`
}

type FunctionDetail struct {
	Property string
	Value    string
}
