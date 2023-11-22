package sdk

import "context"

type Procedures interface {
	CreateProcedureForJava(ctx context.Context, request *CreateProcedureForJavaProcedureRequest) error
	CreateProcedureForJavaScript(ctx context.Context, request *CreateProcedureForJavaScriptProcedureRequest) error
	CreateProcedureForPython(ctx context.Context, request *CreateProcedureForPythonProcedureRequest) error
	CreateProcedureForScala(ctx context.Context, request *CreateProcedureForScalaProcedureRequest) error
	CreateProcedureForSQL(ctx context.Context, request *CreateProcedureForSQLProcedureRequest) error
	Alter(ctx context.Context, request *AlterProcedureRequest) error
	Drop(ctx context.Context, request *DropProcedureRequest) error
	Show(ctx context.Context, request *ShowProcedureRequest) ([]Procedure, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Procedure, error)
	Describe(ctx context.Context, request *DescribeProcedureRequest) ([]ProcedureDetail, error)
}

type ProcedureNullInputBehavior string

const (
	ProcedureNullInputBehaviorCalledOnNullInput ProcedureNullInputBehavior = "CALLED ON NULL INPUT"
	ProcedureNullInputBehaviorReturnNullInput   ProcedureNullInputBehavior = "RETURN NULL ON NULL INPUT"
	ProcedureNullInputBehaviorStrict            ProcedureNullInputBehavior = "STRICT"
)

type ProcedureExecuteAs string

const (
	ProcedureExecuteAsCaller ProcedureExecuteAs = "EXECUTE AS CALLER"
	ProcedureExecuteAsOwner  ProcedureExecuteAs = "EXECUTE AS OWNER"
)

// CreateProcedureForJavaProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure.
type CreateProcedureForJavaProcedureOptions struct {
	create                     bool                        `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                       `ddl:"keyword" sql:"OR REPLACE"`
	Secure                     *bool                       `ddl:"keyword" sql:"SECURE"`
	procedure                  bool                        `ddl:"static" sql:"PROCEDURE"`
	name                       SchemaObjectIdentifier      `ddl:"identifier"`
	Arguments                  []ProcedureArgument         `ddl:"parameter,parentheses,no_equals"`
	CopyGrants                 *bool                       `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    *ProcedureReturns           `ddl:"keyword" sql:"RETURNS"`
	languageJava               bool                        `ddl:"static" sql:"LANGUAGE JAVA"`
	RuntimeVersion             *string                     `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages                   []ProcedurePackage          `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports                    []ProcedureImport           `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler                    string                      `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier   `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []ProcedureSecret           `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                     `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	NullInputBehavior          *ProcedureNullInputBehavior `ddl:"keyword"`
	Comment                    *string                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs                  *ProcedureExecuteAs         `ddl:"keyword"`
	ProcedureDefinition        *string                     `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type ProcedureArgument struct {
	ArgName     string   `ddl:"keyword,no_quotes"`
	ArgDataType DataType `ddl:"keyword,no_quotes"`
}

type ProcedureReturns struct {
	ResultDataType *ProcedureReturnsResultDataType `ddl:"keyword"`
	Table          *ProcedureReturnsTable          `ddl:"keyword" sql:"TABLE"`
}

type ProcedureReturnsResultDataType struct {
	ResultDataType DataType `ddl:"keyword"`
	Null           *bool    `ddl:"keyword" sql:"NULL"`
	NotNull        *bool    `ddl:"keyword" sql:"NOT NULL"`
}

type ProcedureReturnsTable struct {
	Columns []ProcedureColumn `ddl:"parameter,parentheses,no_equals"`
}

type ProcedureColumn struct {
	ColumnName     string   `ddl:"keyword,no_quotes"`
	ColumnDataType DataType `ddl:"keyword,no_quotes"`
}

type ProcedurePackage struct {
	Package string `ddl:"keyword,single_quotes"`
}

type ProcedureImport struct {
	Import string `ddl:"keyword,single_quotes"`
}

type ProcedureSecret struct {
	SecretVariableName string `ddl:"keyword,single_quotes"`
	SecretName         string `ddl:"parameter,no_quotes"`
}

// CreateProcedureForJavaScriptProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure.
type CreateProcedureForJavaScriptProcedureOptions struct {
	create              bool                        `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                       `ddl:"keyword" sql:"OR REPLACE"`
	Secure              *bool                       `ddl:"keyword" sql:"SECURE"`
	procedure           bool                        `ddl:"static" sql:"PROCEDURE"`
	name                SchemaObjectIdentifier      `ddl:"identifier"`
	Arguments           []ProcedureArgument         `ddl:"parameter,parentheses,no_equals"`
	CopyGrants          *bool                       `ddl:"keyword" sql:"COPY GRANTS"`
	Returns             *ProcedureReturns2          `ddl:"keyword" sql:"RETURNS"`
	languageJavascript  bool                        `ddl:"static" sql:"LANGUAGE JAVASCRIPT"`
	NullInputBehavior   *ProcedureNullInputBehavior `ddl:"keyword"`
	Comment             *string                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs           *ProcedureExecuteAs         `ddl:"keyword"`
	ProcedureDefinition *string                     `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type ProcedureReturns2 struct {
	ResultDataType DataType `ddl:"keyword"`
	NotNull        *bool    `ddl:"keyword" sql:"NOT NULL"`
}

// CreateProcedureForPythonProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure.
type CreateProcedureForPythonProcedureOptions struct {
	create                     bool                        `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                       `ddl:"keyword" sql:"OR REPLACE"`
	Secure                     *bool                       `ddl:"keyword" sql:"SECURE"`
	procedure                  bool                        `ddl:"static" sql:"PROCEDURE"`
	name                       SchemaObjectIdentifier      `ddl:"identifier"`
	Arguments                  []ProcedureArgument         `ddl:"parameter,parentheses,no_equals"`
	CopyGrants                 *bool                       `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    *ProcedureReturns           `ddl:"keyword" sql:"RETURNS"`
	languagePython             bool                        `ddl:"static" sql:"LANGUAGE PYTHON"`
	RuntimeVersion             *string                     `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages                   []ProcedurePackage          `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports                    []ProcedureImport           `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler                    string                      `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier   `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []ProcedureSecret           `ddl:"parameter,parentheses" sql:"SECRETS"`
	NullInputBehavior          *ProcedureNullInputBehavior `ddl:"keyword"`
	Comment                    *string                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs                  *ProcedureExecuteAs         `ddl:"keyword"`
	ProcedureDefinition        *string                     `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateProcedureForScalaProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure.
type CreateProcedureForScalaProcedureOptions struct {
	create              bool                        `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                       `ddl:"keyword" sql:"OR REPLACE"`
	Secure              *bool                       `ddl:"keyword" sql:"SECURE"`
	procedure           bool                        `ddl:"static" sql:"PROCEDURE"`
	name                SchemaObjectIdentifier      `ddl:"identifier"`
	Arguments           []ProcedureArgument         `ddl:"parameter,parentheses,no_equals"`
	CopyGrants          *bool                       `ddl:"keyword" sql:"COPY GRANTS"`
	Returns             *ProcedureReturns           `ddl:"keyword" sql:"RETURNS"`
	languageScala       bool                        `ddl:"static" sql:"LANGUAGE SCALA"`
	RuntimeVersion      *string                     `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages            []ProcedurePackage          `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports             []ProcedureImport           `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler             string                      `ddl:"parameter,single_quotes" sql:"HANDLER"`
	TargetPath          *string                     `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	NullInputBehavior   *ProcedureNullInputBehavior `ddl:"keyword"`
	Comment             *string                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs           *ProcedureExecuteAs         `ddl:"keyword"`
	ProcedureDefinition *string                     `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateProcedureForSQLProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure.
type CreateProcedureForSQLProcedureOptions struct {
	create              bool                        `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                       `ddl:"keyword" sql:"OR REPLACE"`
	Secure              *bool                       `ddl:"keyword" sql:"SECURE"`
	procedure           bool                        `ddl:"static" sql:"PROCEDURE"`
	name                SchemaObjectIdentifier      `ddl:"identifier"`
	Arguments           []ProcedureArgument         `ddl:"parameter,parentheses,no_equals"`
	CopyGrants          *bool                       `ddl:"keyword" sql:"COPY GRANTS"`
	Returns             *ProcedureReturns3          `ddl:"keyword" sql:"RETURNS"`
	languageSql         bool                        `ddl:"static" sql:"LANGUAGE SQL"`
	NullInputBehavior   *ProcedureNullInputBehavior `ddl:"keyword"`
	Comment             *string                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs           *ProcedureExecuteAs         `ddl:"keyword"`
	ProcedureDefinition *string                     `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type ProcedureReturns3 struct {
	ResultDataType *ProcedureReturnsResultDataType `ddl:"keyword"`
	Table          *ProcedureReturnsTable          `ddl:"keyword" sql:"TABLE"`
	NotNull        *bool                           `ddl:"keyword" sql:"NOT NULL"`
}

// AlterProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-procedure.
type AlterProcedureOptions struct {
	alter         bool                    `ddl:"static" sql:"ALTER"`
	procedure     bool                    `ddl:"static" sql:"PROCEDURE"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier  `ddl:"identifier"`
	ArgumentTypes []ProcedureArgumentType `ddl:"parameter,parentheses,no_equals"`
	Set           *ProcedureSet           `ddl:"keyword" sql:"SET"`
	Unset         *ProcedureUnset         `ddl:"keyword" sql:"UNSET"`
	ExecuteAs     *ProcedureExecuteAs     `ddl:"keyword"`
	RenameTo      *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SetTags       []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags     []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

type ProcedureArgumentType struct {
	ArgDataType DataType `ddl:"keyword,no_quotes"`
}

type ProcedureSet struct {
	LogLevel   *string `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	TraceLevel *string `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	Comment    *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ProcedureUnset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-procedure.
type DropProcedureOptions struct {
	drop          bool                    `ddl:"static" sql:"DROP"`
	procedure     bool                    `ddl:"static" sql:"PROCEDURE"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier  `ddl:"identifier"`
	ArgumentTypes []ProcedureArgumentType `ddl:"parameter,parentheses,no_equals"`
}

// ShowProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-procedures.
type ShowProcedureOptions struct {
	show       bool  `ddl:"static" sql:"SHOW"`
	procedures bool  `ddl:"static" sql:"PROCEDURES"`
	Like       *Like `ddl:"keyword" sql:"LIKE"`
	In         *In   `ddl:"keyword" sql:"IN"`
}

type procedureRow struct {
	CreatedOn       string `db:"created_on"`
	Name            string `db:"name"`
	SchemaName      string `db:"schema_name"`
	MinNumArguments int    `db:"min_num_arguments"`
	MaxNumArguments int    `db:"max_num_arguments"`
	Arguments       string `db:"arguments"`
	IsTableFunction string `db:"is_table_function"`
}

type Procedure struct {
	CreatedOn       string
	Name            string
	SchemaName      string
	MinNumArguments int
	MaxNumArguments int
	Arguments       string
	IsTableFunction string
}

// DescribeProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/describe-procedure.
type DescribeProcedureOptions struct {
	describe      bool                    `ddl:"static" sql:"DESCRIBE"`
	procedure     bool                    `ddl:"static" sql:"PROCEDURE"`
	name          SchemaObjectIdentifier  `ddl:"identifier"`
	ArgumentTypes []ProcedureArgumentType `ddl:"parameter,parentheses,no_equals"`
}

type procedureDetailRow struct {
	Property string `db:"property"`
	Value    string `db:"value"`
}

type ProcedureDetail struct {
	Property string
	Value    string
}
