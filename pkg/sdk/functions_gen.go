package sdk

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"log"
	"strings"
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
	CreatedOn       string
	Name            string
	SchemaName      string
	IsBuiltin       bool
	IsAggregate     bool
	IsAnsi          bool
	MinNumArguments int
	MaxNumArguments int
	// TODO(SNOW-function refactor): Remove raw arguments
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

// parseFunctionArgumentsFromDetails parses arguments from signature that is contained in function details
func parseFunctionArgumentsFromDetails(details []FunctionDetail) ([]DataType, error) {
	signatureProperty, err := collections.FindOne(details, func(detail FunctionDetail) bool { return detail.Property == "signature" })
	if err != nil {
		return nil, err
	}
	// signature has a structure of (<column name> <data type>, ...); column names can contain commas and other special characters,
	// and they're not escaped, meaning for names such as `"a,b.c|d e" NUMBER` the signature will contain `(a,b.c|d e NUMBER)`.
	arguments := make([]DataType, 0)
	// TODO(TODO - ticket number): Handle arguments with comma in the name (right now this could break for arguments containing dots)
	for _, arg := range strings.Split(strings.Trim(signatureProperty.Value, "()"), ",") {
		// single argument has a structure of <column name> <data type>
		argumentSignatureParts := strings.Split(strings.TrimSpace(arg), " ")
		arguments = append(arguments, DataType(argumentSignatureParts[len(argumentSignatureParts)-1]))
	}
	return arguments, nil
}

// Move to sdk/identifier_parsers.go
func parseFunctionArgumentsFromString(arguments string) ([]DataType, error) {
	dataTypes := make([]DataType, 0)

	stringBuffer := bytes.NewBufferString(arguments)
	for stringBuffer.Len() > 0 {

		// we use another buffer to peek into next data type
		peekBuffer := bytes.NewBufferString(stringBuffer.String())
		peekDataType, _ := peekBuffer.ReadString(',')
		peekDataType = strings.TrimSpace(peekDataType)

		// For function purposes only Vector needs special case
		switch {
		case strings.HasPrefix(peekDataType, "VECTOR"):
			vectorDataType, _ := stringBuffer.ReadString(')')
			vectorDataType = strings.TrimSpace(vectorDataType)
			if stringBuffer.Len() > 0 {
				commaByte, err := stringBuffer.ReadByte()
				if commaByte != ',' {
					return nil, fmt.Errorf("expected a comma delimited string but found %s", string(commaByte))
				}
				if err != nil {
					return nil, err
				}
			}
			log.Println("Adding vec:", vectorDataType)
			dataTypes = append(dataTypes, DataType(vectorDataType))
		default:
			dataType, err := stringBuffer.ReadString(',')
			if err == nil {
				dataType = dataType[:len(dataType)-1]
			}
			dataType = strings.TrimSpace(dataType)
			log.Println("Adding:", dataType)
			dataTypes = append(dataTypes, DataType(dataType))
		}
	}

	return dataTypes, nil
}

func (v *Function) ID(details []FunctionDetail) (SchemaObjectIdentifierWithArguments, error) {
	arguments, err := parseFunctionArgumentsFromDetails(details)
	if err != nil {
		return SchemaObjectIdentifierWithArguments{}, err
	}
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, arguments...), nil
}

// DescribeFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-function.
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
