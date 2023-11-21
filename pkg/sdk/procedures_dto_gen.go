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
	StrictOrNot                *ProcedureStrictOrNotRequest
	VolatileOrNot              *ProcedureVolatileOrNotRequest
	Comment                    *string
	ExecuteAs                  *ProcedureExecuteAsRequest
	As                         *string
}

type ProcedureArgumentRequest struct {
	ArgName     string
	ArgDataType string
}

type ProcedureReturnsRequest struct {
	ResultDataType *ProcedureReturnsResultDataTypeRequest
	Table          *ProcedureReturnsTableRequest
}

type ProcedureReturnsResultDataTypeRequest struct {
	ResultDataType string
	Null           *bool
	NotNull        *bool
}

type ProcedureReturnsTableRequest struct {
	Columns []ProcedureColumnRequest
}

type ProcedureColumnRequest struct {
	ColumnName     string
	ColumnDataType string
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

type ProcedureStrictOrNotRequest struct {
	Strict            *bool
	CalledOnNullInput *bool
}

type ProcedureVolatileOrNotRequest struct {
	Volatile  *bool
	Immutable *bool
}

type ProcedureExecuteAsRequest struct {
	Caller *bool
	Owner  *bool
}

type CreateProcedureForJavaScriptProcedureRequest struct {
	OrReplace     *bool
	Secure        *bool
	name          SchemaObjectIdentifier // required
	Arguments     []ProcedureArgumentRequest
	CopyGrants    *bool
	Returns       *ProcedureReturns2Request
	StrictOrNot   *ProcedureStrictOrNotRequest
	VolatileOrNot *ProcedureVolatileOrNotRequest
	Comment       *string
	ExecuteAs     *ProcedureExecuteAsRequest
	As            *string
}

type ProcedureReturns2Request struct {
	ResultDataType string
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
	StrictOrNot                *ProcedureStrictOrNotRequest
	VolatileOrNot              *ProcedureVolatileOrNotRequest
	Comment                    *string
	ExecuteAs                  *ProcedureExecuteAsRequest
	As                         *string
}

type CreateProcedureForScalaProcedureRequest struct {
	OrReplace      *bool
	Secure         *bool
	name           SchemaObjectIdentifier // required
	Arguments      []ProcedureArgumentRequest
	CopyGrants     *bool
	Returns        *ProcedureReturnsRequest
	RuntimeVersion *string
	Packages       []ProcedurePackageRequest
	Imports        []ProcedureImportRequest
	Handler        string
	TargetPath     *string
	StrictOrNot    *ProcedureStrictOrNotRequest
	VolatileOrNot  *ProcedureVolatileOrNotRequest
	Comment        *string
	ExecuteAs      *ProcedureExecuteAsRequest
	As             *string
}

type CreateProcedureForSQLProcedureRequest struct {
	OrReplace     *bool
	Secure        *bool
	name          SchemaObjectIdentifier // required
	Arguments     []ProcedureArgumentRequest
	CopyGrants    *bool
	Returns       *ProcedureReturns3Request
	StrictOrNot   *ProcedureStrictOrNotRequest
	VolatileOrNot *ProcedureVolatileOrNotRequest
	Comment       *string
	ExecuteAs     *ProcedureExecuteAsRequest
	As            string
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
	ExecuteAs     *ProcedureExecuteAsRequest
	RenameTo      *SchemaObjectIdentifier
	SetTags       []TagAssociation
	UnsetTags     []ObjectIdentifier
}

type ProcedureArgumentTypeRequest struct {
	ArgDataType string
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
