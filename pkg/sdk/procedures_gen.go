package sdk

import (
	"context"
	"database/sql"
)

type Procedures interface {
	CreateForJava(ctx context.Context, request *CreateForJavaProcedureRequest) error
	CreateForJavaScript(ctx context.Context, request *CreateForJavaScriptProcedureRequest) error
	CreateForPython(ctx context.Context, request *CreateForPythonProcedureRequest) error
	CreateForScala(ctx context.Context, request *CreateForScalaProcedureRequest) error
	CreateForSQL(ctx context.Context, request *CreateForSQLProcedureRequest) error
	Alter(ctx context.Context, request *AlterProcedureRequest) error
	Drop(ctx context.Context, request *DropProcedureRequest) error
	Show(ctx context.Context, request *ShowProcedureRequest) ([]Procedure, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Procedure, error)
	Describe(ctx context.Context, request *DescribeProcedureRequest) ([]ProcedureDetail, error)
}

// CreateForJavaProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#java-handler.
type CreateForJavaProcedureOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	procedure                  bool                      `ddl:"static" sql:"PROCEDURE"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []ProcedureArgument       `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    ProcedureReturns          `ddl:"keyword" sql:"RETURNS"`
	languageJava               bool                      `ddl:"static" sql:"LANGUAGE JAVA"`
	RuntimeVersion             string                    `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages                   []ProcedurePackage        `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports                    []ProcedureImport         `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []Secret                  `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                   `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs                  *ExecuteAs                `ddl:"keyword"`
	ProcedureDefinition        *string                   `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type ProcedureArgument struct {
	ArgName      string   `ddl:"keyword,no_quotes"`
	ArgDataType  DataType `ddl:"keyword,no_quotes"`
	DefaultValue *string  `ddl:"parameter,no_equals" sql:"DEFAULT"`
}

type ProcedureReturns struct {
	ResultDataType *ProcedureReturnsResultDataType `ddl:"keyword"`
	Table          *ProcedureReturnsTable          `ddl:"keyword" sql:"TABLE"`
}

type ProcedureReturnsResultDataType struct {
	ResultDataType DataType `ddl:"keyword,no_quotes"`
	Null           *bool    `ddl:"keyword" sql:"NULL"`
	NotNull        *bool    `ddl:"keyword" sql:"NOT NULL"`
}

type ProcedureReturnsTable struct {
	Columns []ProcedureColumn `ddl:"list,must_parentheses"`
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

// CreateForJavaScriptProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#javascript-handler.
type CreateForJavaScriptProcedureOptions struct {
	create              bool                   `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Secure              *bool                  `ddl:"keyword" sql:"SECURE"`
	procedure           bool                   `ddl:"static" sql:"PROCEDURE"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	Arguments           []ProcedureArgument    `ddl:"list,must_parentheses"`
	CopyGrants          *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	ResultDataType      DataType               `ddl:"parameter,no_equals" sql:"RETURNS"`
	NotNull             *bool                  `ddl:"keyword" sql:"NOT NULL"`
	languageJavascript  bool                   `ddl:"static" sql:"LANGUAGE JAVASCRIPT"`
	NullInputBehavior   *NullInputBehavior     `ddl:"keyword"`
	Comment             *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs           *ExecuteAs             `ddl:"keyword"`
	ProcedureDefinition string                 `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateForPythonProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#python-handler.
type CreateForPythonProcedureOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	procedure                  bool                      `ddl:"static" sql:"PROCEDURE"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []ProcedureArgument       `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    ProcedureReturns          `ddl:"keyword" sql:"RETURNS"`
	languagePython             bool                      `ddl:"static" sql:"LANGUAGE PYTHON"`
	RuntimeVersion             string                    `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages                   []ProcedurePackage        `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports                    []ProcedureImport         `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []Secret                  `ddl:"parameter,parentheses" sql:"SECRETS"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs                  *ExecuteAs                `ddl:"keyword"`
	ProcedureDefinition        *string                   `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateForScalaProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#scala-handler.
type CreateForScalaProcedureOptions struct {
	create              bool                   `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Secure              *bool                  `ddl:"keyword" sql:"SECURE"`
	procedure           bool                   `ddl:"static" sql:"PROCEDURE"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	Arguments           []ProcedureArgument    `ddl:"list,must_parentheses"`
	CopyGrants          *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Returns             ProcedureReturns       `ddl:"keyword" sql:"RETURNS"`
	languageScala       bool                   `ddl:"static" sql:"LANGUAGE SCALA"`
	RuntimeVersion      string                 `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages            []ProcedurePackage     `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports             []ProcedureImport      `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler             string                 `ddl:"parameter,single_quotes" sql:"HANDLER"`
	TargetPath          *string                `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	NullInputBehavior   *NullInputBehavior     `ddl:"keyword"`
	Comment             *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs           *ExecuteAs             `ddl:"keyword"`
	ProcedureDefinition *string                `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

// CreateForSQLProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#snowflake-scripting-handler.
type CreateForSQLProcedureOptions struct {
	create              bool                   `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Secure              *bool                  `ddl:"keyword" sql:"SECURE"`
	procedure           bool                   `ddl:"static" sql:"PROCEDURE"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	Arguments           []ProcedureArgument    `ddl:"list,must_parentheses"`
	CopyGrants          *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Returns             ProcedureSQLReturns    `ddl:"keyword" sql:"RETURNS"`
	languageSql         bool                   `ddl:"static" sql:"LANGUAGE SQL"`
	NullInputBehavior   *NullInputBehavior     `ddl:"keyword"`
	Comment             *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs           *ExecuteAs             `ddl:"keyword"`
	ProcedureDefinition string                 `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
}

type ProcedureSQLReturns struct {
	ResultDataType *ProcedureReturnsResultDataType `ddl:"keyword"`
	Table          *ProcedureReturnsTable          `ddl:"keyword" sql:"TABLE"`
	NotNull        *bool                           `ddl:"keyword" sql:"NOT NULL"`
}

// AlterProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-procedure.
type AlterProcedureOptions struct {
	alter             bool                    `ddl:"static" sql:"ALTER"`
	procedure         bool                    `ddl:"static" sql:"PROCEDURE"`
	IfExists          *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name              SchemaObjectIdentifier  `ddl:"identifier"`
	ArgumentDataTypes []DataType              `ddl:"keyword,parentheses"`
	RenameTo          *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SetComment        *string                 `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	SetLogLevel       *string                 `ddl:"parameter,single_quotes" sql:"SET LOG_LEVEL"`
	SetTraceLevel     *string                 `ddl:"parameter,single_quotes" sql:"SET TRACE_LEVEL"`
	UnsetComment      *bool                   `ddl:"keyword" sql:"UNSET COMMENT"`
	SetTags           []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags         []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
	ExecuteAs         *ExecuteAs              `ddl:"keyword"`
}

// DropProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-procedure.
type DropProcedureOptions struct {
	drop              bool                   `ddl:"static" sql:"DROP"`
	procedure         bool                   `ddl:"static" sql:"PROCEDURE"`
	IfExists          *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	ArgumentDataTypes []DataType             `ddl:"keyword,must_parentheses"`
}

// ShowProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-procedures.
type ShowProcedureOptions struct {
	show       bool  `ddl:"static" sql:"SHOW"`
	procedures bool  `ddl:"static" sql:"PROCEDURES"`
	Like       *Like `ddl:"keyword" sql:"LIKE"`
	In         *In   `ddl:"keyword" sql:"IN"`
}

type procedureRow struct {
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
}

type Procedure struct {
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
}

// DescribeProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-procedure.
type DescribeProcedureOptions struct {
	describe          bool                   `ddl:"static" sql:"DESCRIBE"`
	procedure         bool                   `ddl:"static" sql:"PROCEDURE"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	ArgumentDataTypes []DataType             `ddl:"keyword,parentheses"`
}

type procedureDetailRow struct {
	Property string `db:"property"`
	Value    string `db:"value"`
}

type ProcedureDetail struct {
	Property string
	Value    string
}
