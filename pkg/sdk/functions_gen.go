package sdk

import "context"

type Functions interface {
	CreateFunctionForJava(ctx context.Context, request *CreateFunctionForJavaFunctionRequest) error
	CreateFunctionForJavascript(ctx context.Context, request *CreateFunctionForJavascriptFunctionRequest) error
	CreateFunctionForPython(ctx context.Context, request *CreateFunctionForPythonFunctionRequest) error
	CreateFunctionForScala(ctx context.Context, request *CreateFunctionForScalaFunctionRequest) error
	CreateFunctionForSQL(ctx context.Context, request *CreateFunctionForSQLFunctionRequest) error
	Alter(ctx context.Context, request *AlterFunctionRequest) error
	Drop(ctx context.Context, request *DropFunctionRequest) error
	Show(ctx context.Context, request *ShowFunctionRequest) ([]Function, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Function, error)
	Describe(ctx context.Context, request *DescribeFunctionRequest) ([]FunctionDetail, error)
}

type FunctionNullInputBehavior string

const (
	FunctionNullInputBehaviorCalledOnNullInput FunctionNullInputBehavior = "CALLED ON NULL INPUT"
	FunctionNullInputBehaviorReturnNullInput   FunctionNullInputBehavior = "RETURN NULL ON NULL INPUT"
	FunctionNullInputBehaviorStrict            FunctionNullInputBehavior = "STRICT"
)

type FunctionReturnResultsBehavior string

const (
	FunctionReturnResultsBehaviorVolatile  FunctionReturnResultsBehavior = "VOLATILE"
	FunctionReturnResultsBehaviorImmutable FunctionReturnResultsBehavior = "IMMUTABLE"
)

type FunctionReturnNullValues string

const (
	FunctionReturnNullValuesNull    FunctionReturnNullValues = "NULL"
	FunctionReturnNullValuesNotNull FunctionReturnNullValues = "NOT NULL"
)

// CreateFunctionForJavaFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function.
type CreateFunctionForJavaFunctionOptions struct {
	create                     bool                           `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                          `ddl:"keyword" sql:"OR REPLACE"`
	Temporary                  *bool                          `ddl:"keyword" sql:"TEMPORARY"`
	Secure                     *bool                          `ddl:"keyword" sql:"SECURE"`
	function                   bool                           `ddl:"static" sql:"FUNCTION"`
	IfNotExists                *bool                          `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier         `ddl:"identifier"`
	Arguments                  []FunctionArgument             `ddl:"parameter,parentheses,no_equals"`
	CopyGrants                 *bool                          `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    *FunctionReturns               `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues           *FunctionReturnNullValues      `ddl:"keyword"`
	languageJava               bool                           `ddl:"static" sql:"LANGUAGE JAVA"`
	NullInputBehavior          *FunctionNullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior      *FunctionReturnResultsBehavior `ddl:"keyword"`
	RuntimeVersion             *string                        `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment                    *string                        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports                    []FunctionImports              `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages                   []FunctionPackages             `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler                    *string                        `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier      `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []FunctionSecret               `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                        `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	FunctionDefinition         *string                        `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type FunctionArgument struct {
	ArgName     string   `ddl:"keyword,no_quotes"`
	ArgDataType DataType `ddl:"keyword,no_quotes"`
}

type FunctionReturns struct {
	ResultDataType *DataType             `ddl:"keyword"`
	Table          *FunctionReturnsTable `ddl:"keyword" sql:"TABLE"`
}

type FunctionReturnsTable struct {
	Columns []FunctionColumn `ddl:"parameter,parentheses,no_equals"`
}

type FunctionColumn struct {
	ColumnName     string   `ddl:"keyword,no_quotes"`
	ColumnDataType DataType `ddl:"keyword,no_quotes"`
}

type FunctionImports struct {
	Import string `ddl:"keyword,single_quotes"`
}

type FunctionPackages struct {
	Package string `ddl:"keyword,single_quotes"`
}

type FunctionSecret struct {
	SecretVariableName string `ddl:"keyword,single_quotes"`
	SecretName         string `ddl:"parameter,no_quotes"`
}

// CreateFunctionForJavascriptFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function.
type CreateFunctionForJavascriptFunctionOptions struct {
	create                bool                           `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                          `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                          `ddl:"keyword" sql:"TEMPORARY"`
	Secure                *bool                          `ddl:"keyword" sql:"SECURE"`
	function              bool                           `ddl:"static" sql:"FUNCTION"`
	IfNotExists           *bool                          `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier         `ddl:"identifier"`
	Arguments             []FunctionArgument             `ddl:"parameter,parentheses,no_equals"`
	CopyGrants            *bool                          `ddl:"keyword" sql:"COPY GRANTS"`
	Returns               *FunctionReturns               `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues      *FunctionReturnNullValues      `ddl:"keyword"`
	languageJavascript    bool                           `ddl:"static" sql:"LANGUAGE JAVASCRIPT"`
	NullInputBehavior     *FunctionNullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior *FunctionReturnResultsBehavior `ddl:"keyword"`
	Comment               *string                        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	FunctionDefinition    *string                        `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateFunctionForPythonFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function.
type CreateFunctionForPythonFunctionOptions struct {
	create                     bool                           `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                          `ddl:"keyword" sql:"OR REPLACE"`
	Temporary                  *bool                          `ddl:"keyword" sql:"TEMPORARY"`
	Secure                     *bool                          `ddl:"keyword" sql:"SECURE"`
	function                   bool                           `ddl:"static" sql:"FUNCTION"`
	IfNotExists                *bool                          `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier         `ddl:"identifier"`
	Arguments                  []FunctionArgument             `ddl:"parameter,parentheses,no_equals"`
	CopyGrants                 *bool                          `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    *FunctionReturns               `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues           *FunctionReturnNullValues      `ddl:"keyword"`
	languagePython             bool                           `ddl:"static" sql:"LANGUAGE PYTHON"`
	NullInputBehavior          *FunctionNullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior      *FunctionReturnResultsBehavior `ddl:"keyword"`
	RuntimeVersion             *string                        `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment                    *string                        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports                    []FunctionImports              `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages                   []FunctionPackages             `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler                    *string                        `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier      `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []FunctionSecret               `ddl:"parameter,parentheses" sql:"SECRETS"`
	FunctionDefinition         *string                        `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateFunctionForScalaFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function.
type CreateFunctionForScalaFunctionOptions struct {
	create                bool                           `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                          `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                          `ddl:"keyword" sql:"TEMPORARY"`
	Secure                *bool                          `ddl:"keyword" sql:"SECURE"`
	function              bool                           `ddl:"static" sql:"FUNCTION"`
	IfNotExists           *bool                          `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier         `ddl:"identifier"`
	Arguments             []FunctionArgument             `ddl:"parameter,parentheses,no_equals"`
	CopyGrants            *bool                          `ddl:"keyword" sql:"COPY GRANTS"`
	Returns               *FunctionReturns               `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues      *FunctionReturnNullValues      `ddl:"keyword"`
	languageScala         bool                           `ddl:"static" sql:"LANGUAGE SCALA"`
	NullInputBehavior     *FunctionNullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior *FunctionReturnResultsBehavior `ddl:"keyword"`
	RuntimeVersion        *string                        `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment               *string                        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports               []FunctionImports              `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages              []FunctionPackages             `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler               *string                        `ddl:"parameter,single_quotes" sql:"HANDLER"`
	TargetPath            *string                        `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	FunctionDefinition    *string                        `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateFunctionForSQLFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function.
type CreateFunctionForSQLFunctionOptions struct {
	create                bool                           `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                          `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                          `ddl:"keyword" sql:"TEMPORARY"`
	Secure                *bool                          `ddl:"keyword" sql:"SECURE"`
	function              bool                           `ddl:"static" sql:"FUNCTION"`
	IfNotExists           *bool                          `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier         `ddl:"identifier"`
	Arguments             []FunctionArgument             `ddl:"parameter,parentheses,no_equals"`
	CopyGrants            *bool                          `ddl:"keyword" sql:"COPY GRANTS"`
	Returns               *FunctionReturns               `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues      *FunctionReturnNullValues      `ddl:"keyword"`
	ReturnResultsBehavior *FunctionReturnResultsBehavior `ddl:"keyword"`
	Memoizable            *bool                          `ddl:"keyword" sql:"MEMOIZABLE"`
	Comment               *string                        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	FunctionDefinition    *string                        `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// AlterFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-function.
type AlterFunctionOptions struct {
	alter         bool                    `ddl:"static" sql:"ALTER"`
	function      bool                    `ddl:"static" sql:"FUNCTION"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier  `ddl:"identifier"`
	ArgumentTypes []FunctionArgumentType  `ddl:"parameter,parentheses,no_equals"`
	Set           *FunctionSet            `ddl:"keyword" sql:"SET"`
	Unset         *FunctionUnset          `ddl:"keyword" sql:"UNSET"`
	RenameTo      *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SetTags       []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags     []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

type FunctionArgumentType struct {
	ArgDataType DataType `ddl:"keyword,no_quotes"`
}

type FunctionSet struct {
	LogLevel   *string `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	TraceLevel *string `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	Comment    *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Secure     *bool   `ddl:"keyword" sql:"SECURE"`
}

type FunctionUnset struct {
	Secure     *bool `ddl:"keyword" sql:"SECURE"`
	Comment    *bool `ddl:"keyword" sql:"COMMENT"`
	LogLevel   *bool `ddl:"keyword" sql:"LOG_LEVEL"`
	TraceLevel *bool `ddl:"keyword" sql:"TRACE_LEVEL"`
}

// DropFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-function.
type DropFunctionOptions struct {
	drop          bool                   `ddl:"static" sql:"DROP"`
	function      bool                   `ddl:"static" sql:"FUNCTION"`
	IfExists      *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	ArgumentTypes []FunctionArgumentType `ddl:"parameter,parentheses,no_equals"`
}

// ShowFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-functions.
type ShowFunctionOptions struct {
	show          bool  `ddl:"static" sql:"SHOW"`
	userFunctions bool  `ddl:"static" sql:"USER FUNCTIONS"`
	Like          *Like `ddl:"keyword" sql:"LIKE"`
	In            *In   `ddl:"keyword" sql:"IN"`
}

type functionRow struct {
	CreatedOn          string `db:"created_on"`
	Name               string `db:"name"`
	SchemaName         string `db:"schema_name"`
	MinNumArguments    int    `db:"min_num_arguments"`
	MaxNumArguments    int    `db:"max_num_arguments"`
	Arguments          string `db:"arguments"`
	IsTableFunction    string `db:"is_table_function"`
	IsSecure           string `db:"is_secure"`
	IsExternalFunction string `db:"is_external_function"`
	Language           string `db:"language"`
	IsMemoizable       string `db:"is_memoizable"`
}

type Function struct {
	CreatedOn          string
	Name               string
	SchemaName         string
	MinNumArguments    int
	MaxNumArguments    int
	Arguments          string
	IsTableFunction    bool
	IsSecure           bool
	IsExternalFunction bool
	Language           string
	IsMemoizable       bool
}

// DescribeFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/describe-function.
type DescribeFunctionOptions struct {
	describe      bool                   `ddl:"static" sql:"DESCRIBE"`
	function      bool                   `ddl:"static" sql:"FUNCTION"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	ArgumentTypes []FunctionArgumentType `ddl:"parameter,parentheses,no_equals"`
}

type functionDetailRow struct {
	Property string `db:"property"`
	Value    string `db:"value"`
}

type FunctionDetail struct {
	Property string
	Value    string
}
