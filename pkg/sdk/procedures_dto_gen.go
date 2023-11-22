package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateProcedureForJavaProcedureOptions]       = new(CreateProcedureForJavaProcedureRequest)
	_ optionsProvider[CreateProcedureForJavaScriptProcedureOptions] = new(CreateProcedureForJavaScriptProcedureRequest)
	_ optionsProvider[CreateProcedureForPythonProcedureOptions]     = new(CreateProcedureForPythonProcedureRequest)
	_ optionsProvider[CreateProcedureForScalaProcedureOptions]      = new(CreateProcedureForScalaProcedureRequest)
	_ optionsProvider[CreateProcedureForSQLProcedureOptions]        = new(CreateProcedureForSQLProcedureRequest)
	_ optionsProvider[AlterProcedureOptions]                        = new(AlterProcedureRequest)
	_ optionsProvider[DropProcedureOptions]                         = new(DropProcedureRequest)
	_ optionsProvider[ShowProcedureOptions]                         = new(ShowProcedureRequest)
	_ optionsProvider[DescribeProcedureOptions]                     = new(DescribeProcedureRequest)
)

type CreateProcedureForJavaProcedureRequest struct {
	OrReplace                  *bool
	Secure                     *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []ProcedureArgumentRequest
	CopyGrants                 *bool
	Returns                    *ProcedureReturnsRequest
	RuntimeVersion             *string
	Packages                   []ProcedurePackageRequest
	Imports                    []ProcedureImportRequest
	Handler                    string
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []ProcedureSecretRequest
	TargetPath                 *string
	NullInputBehavior          *ProcedureNullInputBehavior
	Comment                    *string
	ExecuteAs                  *ProcedureExecuteAs
	ProcedureDefinition        *string
}

type ProcedureArgumentRequest struct {
	ArgName     string
	ArgDataType DataType
}

type ProcedureReturnsRequest struct {
	ResultDataType *ProcedureReturnsResultDataTypeRequest
	Table          *ProcedureReturnsTableRequest
}

type ProcedureReturnsResultDataTypeRequest struct {
	ResultDataType DataType
	Null           *bool
	NotNull        *bool
}

type ProcedureReturnsTableRequest struct {
	Columns []ProcedureColumnRequest
}

type ProcedureColumnRequest struct {
	ColumnName     string
	ColumnDataType DataType
}

type ProcedurePackageRequest struct {
	Package string
}

type ProcedureImportRequest struct {
	Import string
}

type ProcedureSecretRequest struct {
	SecretVariableName string
	SecretName         string
}

type CreateProcedureForJavaScriptProcedureRequest struct {
	OrReplace           *bool
	Secure              *bool
	name                SchemaObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	CopyGrants          *bool
	Returns             *ProcedureReturns2Request
	NullInputBehavior   *ProcedureNullInputBehavior
	Comment             *string
	ExecuteAs           *ProcedureExecuteAs
	ProcedureDefinition *string
}

type ProcedureReturns2Request struct {
	ResultDataType DataType
	NotNull        *bool
}

type CreateProcedureForPythonProcedureRequest struct {
	OrReplace                  *bool
	Secure                     *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []ProcedureArgumentRequest
	CopyGrants                 *bool
	Returns                    *ProcedureReturnsRequest
	RuntimeVersion             *string
	Packages                   []ProcedurePackageRequest
	Imports                    []ProcedureImportRequest
	Handler                    string
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []ProcedureSecretRequest
	NullInputBehavior          *ProcedureNullInputBehavior
	Comment                    *string
	ExecuteAs                  *ProcedureExecuteAs
	ProcedureDefinition        *string
}

type CreateProcedureForScalaProcedureRequest struct {
	OrReplace           *bool
	Secure              *bool
	name                SchemaObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	CopyGrants          *bool
	Returns             *ProcedureReturnsRequest
	RuntimeVersion      *string
	Packages            []ProcedurePackageRequest
	Imports             []ProcedureImportRequest
	Handler             string
	TargetPath          *string
	NullInputBehavior   *ProcedureNullInputBehavior
	Comment             *string
	ExecuteAs           *ProcedureExecuteAs
	ProcedureDefinition *string
}

type CreateProcedureForSQLProcedureRequest struct {
	OrReplace           *bool
	Secure              *bool
	name                SchemaObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	CopyGrants          *bool
	Returns             *ProcedureReturns3Request
	NullInputBehavior   *ProcedureNullInputBehavior
	Comment             *string
	ExecuteAs           *ProcedureExecuteAs
	ProcedureDefinition *string
}

type ProcedureReturns3Request struct {
	ResultDataType *ProcedureReturnsResultDataTypeRequest
	Table          *ProcedureReturnsTableRequest
	NotNull        *bool
}

type AlterProcedureRequest struct {
	IfExists      *bool
	name          SchemaObjectIdentifier // required
	ArgumentTypes []ProcedureArgumentTypeRequest
	Set           *ProcedureSetRequest
	Unset         *ProcedureUnsetRequest
	ExecuteAs     *ProcedureExecuteAs
	RenameTo      *SchemaObjectIdentifier
	SetTags       []TagAssociation
	UnsetTags     []ObjectIdentifier
}

type ProcedureArgumentTypeRequest struct {
	ArgDataType DataType
}

type ProcedureSetRequest struct {
	LogLevel   *string
	TraceLevel *string
	Comment    *string
}

type ProcedureUnsetRequest struct {
	Comment *bool
}

type DropProcedureRequest struct {
	IfExists      *bool
	name          SchemaObjectIdentifier // required
	ArgumentTypes []ProcedureArgumentTypeRequest
}

type ShowProcedureRequest struct {
	Like *Like
	In   *In
}

type DescribeProcedureRequest struct {
	name          SchemaObjectIdentifier // required
	ArgumentTypes []ProcedureArgumentTypeRequest
}
