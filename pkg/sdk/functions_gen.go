package sdk

// imports edited manually
import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
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
	ShowByID(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Function, error)
	Describe(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]FunctionDetail, error)

	// DescribeDetails is added manually; it returns aggregated describe results for the given function.
	DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*FunctionDetails, error)
}

// CreateForJavaFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#java-handler.
// TODO [SNOW-1348103 - this PR]: test secure (for each type, with owner and underprivileged role), read https://docs.snowflake.com/en/developer-guide/secure-udf-procedure
// TODO [SNOW-1348103 - this PR]: test setting the paths for all types (like imports, target paths)
// TODO [SNOW-1348103 - this PR]: test weird names for arg name
// TODO [SNOW-1348103 - this PR]: test two types of creation for each func
// TODO [SNOW-1348103 - next PRs]: check data type mappings https://docs.snowflake.com/en/sql-reference/sql/create-function#all-languages (signature + returns)
// TODO [SNOW-1348103 - this PR]: setting RUNTIME_VERSION (only 11.x, 17.x supported, 11.x being the default)
// TODO [SNOW-1348103 - this PR]: packages: package_name:version_number; do we validate? - check SELECT * FROM INFORMATION_SCHEMA.PACKAGES WHERE LANGUAGE = 'java';
// TODO [SNOW-1348103 - next PRs]: add to the resource docs https://docs.snowflake.com/en/sql-reference/sql/create-function#access-control-requirements
// TODO [SNOW-1348103 - this PR]: what delimiter do we use for <function_definition>: ' versus $$? - we use $$ as tasks
// TODO [SNOW-1348103 - this PR]: escaping single quotes test - don't have to do this with $$
// TODO [SNOW-1348103 - this PR]: validation of JAR (check https://docs.snowflake.com/en/sql-reference/sql/create-function#id6)
// TODO [SNOW-1348103 - next PRs]: active warehouse vs validations
// TODO [SNOW-1348103 - this PR]: check creation of all functions (using examples and more)
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
	Secrets                    []SecretReference         `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                   `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	FunctionDefinition         *string                   `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type FunctionArgument struct {
	ArgName        string             `ddl:"keyword,no_quotes"`
	ArgDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ArgDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
	DefaultValue   *string            `ddl:"parameter,no_equals" sql:"DEFAULT"`
}

type FunctionReturns struct {
	ResultDataType *FunctionReturnsResultDataType `ddl:"keyword"`
	Table          *FunctionReturnsTable          `ddl:"keyword" sql:"TABLE"`
}

type FunctionReturnsResultDataType struct {
	ResultDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ResultDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
}

type FunctionReturnsTable struct {
	Columns []FunctionColumn `ddl:"parameter,parentheses,no_equals"`
}

type FunctionColumn struct {
	ColumnName        string             `ddl:"keyword,no_quotes"`
	ColumnDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ColumnDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
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
// TODO [SNOW-1348103 - this PR]: test aggregate func creation
// TODO [SNOW-1348103 - this PR]: what about [==<version>] - SDK level or resource level? check also: SELECT * FROM INFORMATION_SCHEMA.PACKAGES WHERE LANGUAGE = 'python';
// TODO [SNOW-1348103 - this PR]: what about preview feature >= ?
// TODO [SNOW-1348103 - this PR]: what about '<module_file_name>.<function_name>' for non-inline functions?
// TODO [SNOW-1348103 - this PR]: setting RUNTIME_VERSION (only 3.8, 3.9, 3.10, 3.11 supported, which one is a default?)
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
	Secrets                    []SecretReference         `ddl:"parameter,parentheses" sql:"SECRETS"`
	FunctionDefinition         *string                   `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateForScalaFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#scala-handler.
// TODO [SNOW-1348103 - this PR]: setting RUNTIME_VERSION (only 2.12 supported, which is the default)
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
	returns               bool                   `ddl:"static" sql:"RETURNS"`
	ResultDataTypeOld     DataType               `ddl:"parameter,no_equals"`
	ResultDataType        datatypes.DataType     `ddl:"parameter,no_quotes,no_equals"`
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
// TODO [this PR]: can we run multiple sets/unsets?
// TODO [this PR]: add setting EXTERNAL_ACCESS_INTEGRATIONS/SECRETS
// TODO [this PR]: unset EXTERNAL_ACCESS_INTEGRATIONS or SECRETS?
// TODO [this PR]: EXTERNAL_ACCESS_INTEGRATIONS or SECRETS in Javascript or SQL
type AlterFunctionOptions struct {
	alter           bool                                `ddl:"static" sql:"ALTER"`
	function        bool                                `ddl:"static" sql:"FUNCTION"`
	IfExists        *bool                               `ddl:"keyword" sql:"IF EXISTS"`
	name            SchemaObjectIdentifierWithArguments `ddl:"identifier"`
	RenameTo        *SchemaObjectIdentifier             `ddl:"identifier" sql:"RENAME TO"`
	SetComment      *string                             `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	SetLogLevel     *string                             `ddl:"parameter,single_quotes" sql:"SET LOG_LEVEL"`
	SetTraceLevel   *string                             `ddl:"parameter,single_quotes" sql:"SET TRACE_LEVEL"`
	SetSecure       *bool                               `ddl:"keyword" sql:"SET SECURE"`
	UnsetSecure     *bool                               `ddl:"keyword" sql:"UNSET SECURE"`
	UnsetLogLevel   *bool                               `ddl:"keyword" sql:"UNSET LOG_LEVEL"`
	UnsetTraceLevel *bool                               `ddl:"keyword" sql:"UNSET TRACE_LEVEL"`
	UnsetComment    *bool                               `ddl:"keyword" sql:"UNSET COMMENT"`
	SetTags         []TagAssociation                    `ddl:"keyword" sql:"SET TAG"`
	UnsetTags       []ObjectIdentifier                  `ddl:"keyword" sql:"UNSET TAG"`
}

// DropFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-function.
type DropFunctionOptions struct {
	drop     bool                                `ddl:"static" sql:"DROP"`
	function bool                                `ddl:"static" sql:"FUNCTION"`
	IfExists *bool                               `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifierWithArguments `ddl:"identifier"`
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
	ArgumentsOld       []DataType
	ArgumentsRaw       string
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
// TODO [this PR]: create details struct similar to the one in user
// TODO [this PR]: list properties for all types of functions
type DescribeFunctionOptions struct {
	describe bool                                `ddl:"static" sql:"DESCRIBE"`
	function bool                                `ddl:"static" sql:"FUNCTION"`
	name     SchemaObjectIdentifierWithArguments `ddl:"identifier"`
}

type functionDetailRow struct {
	Property string         `db:"property"`
	Value    sql.NullString `db:"value"`
}

type FunctionDetail struct {
	Property string
	Value    string
}
