package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateFunctionForJavaFunctionOptions]       = new(CreateFunctionForJavaFunctionRequest)
	_ optionsProvider[CreateFunctionForJavascriptFunctionOptions] = new(CreateFunctionForJavascriptFunctionRequest)
	_ optionsProvider[CreateFunctionForPythonFunctionOptions]     = new(CreateFunctionForPythonFunctionRequest)
	_ optionsProvider[CreateFunctionForScalaFunctionOptions]      = new(CreateFunctionForScalaFunctionRequest)
	_ optionsProvider[CreateFunctionForSQLFunctionOptions]        = new(CreateFunctionForSQLFunctionRequest)
	_ optionsProvider[AlterFunctionOptions]                       = new(AlterFunctionRequest)
	_ optionsProvider[DropFunctionOptions]                        = new(DropFunctionRequest)
	_ optionsProvider[ShowFunctionOptions]                        = new(ShowFunctionRequest)
	_ optionsProvider[DescribeFunctionOptions]                    = new(DescribeFunctionRequest)
)

type CreateFunctionForJavaFunctionRequest struct {
	OrReplace                  *bool
	Temporary                  *bool
	Secure                     *bool
	IfNotExists                *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []FunctionArgumentRequest
	CopyGrants                 *bool
	Returns                    *FunctionReturnsRequest
	ReturnNullValues           *FunctionReturnNullValues
	NullInputBehavior          *FunctionNullInputBehavior
	ReturnResultsBehavior      *FunctionReturnResultsBehavior
	RuntimeVersion             *string
	Comment                    *string
	Imports                    []FunctionImportsRequest
	Packages                   []FunctionPackagesRequest
	Handler                    string
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []FunctionSecretRequest
	TargetPath                 *string
	FunctionDefinition         string
}

type FunctionArgumentRequest struct {
	ArgName     string
	ArgDataType DataType
	Default     *string
}

type FunctionReturnsRequest struct {
	ResultDataType *DataType
	Table          *FunctionReturnsTableRequest
}

type FunctionReturnsTableRequest struct {
	Columns []FunctionColumnRequest
}

type FunctionColumnRequest struct {
	ColumnName     string
	ColumnDataType DataType
}

type FunctionImportsRequest struct {
	Import string
}

type FunctionPackagesRequest struct {
	Package string
}

type FunctionSecretRequest struct {
	SecretVariableName string
	SecretName         string
}

type CreateFunctionForJavascriptFunctionRequest struct {
	OrReplace             *bool
	Temporary             *bool
	Secure                *bool
	IfNotExists           *bool
	name                  SchemaObjectIdentifier // required
	Arguments             []FunctionArgumentRequest
	CopyGrants            *bool
	Returns               *FunctionReturnsRequest
	ReturnNullValues      *FunctionReturnNullValues
	NullInputBehavior     *FunctionNullInputBehavior
	ReturnResultsBehavior *FunctionReturnResultsBehavior
	Comment               *string
	FunctionDefinition    string
}

type CreateFunctionForPythonFunctionRequest struct {
	OrReplace                  *bool
	Temporary                  *bool
	Secure                     *bool
	IfNotExists                *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []FunctionArgumentRequest
	CopyGrants                 *bool
	Returns                    *FunctionReturnsRequest
	ReturnNullValues           *FunctionReturnNullValues
	NullInputBehavior          *FunctionNullInputBehavior
	ReturnResultsBehavior      *FunctionReturnResultsBehavior
	RuntimeVersion             string
	Comment                    *string
	Imports                    []FunctionImportsRequest
	Packages                   []FunctionPackagesRequest
	Handler                    string
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []FunctionSecretRequest
	FunctionDefinition         string
}

type CreateFunctionForScalaFunctionRequest struct {
	OrReplace             *bool
	Temporary             *bool
	Secure                *bool
	IfNotExists           *bool
	name                  SchemaObjectIdentifier // required
	Arguments             []FunctionArgumentRequest
	CopyGrants            *bool
	Returns               *FunctionReturnsRequest
	ReturnNullValues      *FunctionReturnNullValues
	NullInputBehavior     *FunctionNullInputBehavior
	ReturnResultsBehavior *FunctionReturnResultsBehavior
	RuntimeVersion        *string
	Comment               *string
	Imports               []FunctionImportsRequest
	Packages              []FunctionPackagesRequest
	Handler               string
	TargetPath            *string
	FunctionDefinition    string
}

type CreateFunctionForSQLFunctionRequest struct {
	OrReplace             *bool
	Temporary             *bool
	Secure                *bool
	IfNotExists           *bool
	name                  SchemaObjectIdentifier // required
	Arguments             []FunctionArgumentRequest
	CopyGrants            *bool
	Returns               *FunctionReturnsRequest
	ReturnNullValues      *FunctionReturnNullValues
	ReturnResultsBehavior *FunctionReturnResultsBehavior
	Memoizable            *bool
	Comment               *string
	FunctionDefinition    string
}

type AlterFunctionRequest struct {
	IfExists      *bool
	name          SchemaObjectIdentifier // required
	ArgumentTypes []FunctionArgumentTypeRequest
	Set           *FunctionSetRequest
	Unset         *FunctionUnsetRequest
	RenameTo      *SchemaObjectIdentifier
	SetTags       []TagAssociation
	UnsetTags     []ObjectIdentifier
}

type FunctionArgumentTypeRequest struct {
	ArgDataType DataType
}

type FunctionSetRequest struct {
	LogLevel   *string
	TraceLevel *string
	Comment    *string
	Secure     *bool
}

type FunctionUnsetRequest struct {
	Secure     *bool
	Comment    *bool
	LogLevel   *bool
	TraceLevel *bool
}

type DropFunctionRequest struct {
	IfExists      *bool
	name          SchemaObjectIdentifier // required
	ArgumentTypes []FunctionArgumentTypeRequest
}

type ShowFunctionRequest struct {
	Like *Like
	In   *In
}

type DescribeFunctionRequest struct {
	name          SchemaObjectIdentifier // required
	ArgumentTypes []FunctionArgumentTypeRequest
}
