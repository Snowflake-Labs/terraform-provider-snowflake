package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateForJavaProcedureOptions]       = new(CreateForJavaProcedureRequest)
	_ optionsProvider[CreateForJavaScriptProcedureOptions] = new(CreateForJavaScriptProcedureRequest)
	_ optionsProvider[CreateForPythonProcedureOptions]     = new(CreateForPythonProcedureRequest)
	_ optionsProvider[CreateForScalaProcedureOptions]      = new(CreateForScalaProcedureRequest)
	_ optionsProvider[CreateForSQLProcedureOptions]        = new(CreateForSQLProcedureRequest)
	_ optionsProvider[AlterProcedureOptions]               = new(AlterProcedureRequest)
	_ optionsProvider[DropProcedureOptions]                = new(DropProcedureRequest)
	_ optionsProvider[ShowProcedureOptions]                = new(ShowProcedureRequest)
	_ optionsProvider[DescribeProcedureOptions]            = new(DescribeProcedureRequest)
)

type CreateForJavaProcedureRequest struct {
	OrReplace                  *bool
	Secure                     *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []ProcedureArgumentRequest
	CopyGrants                 *bool
	Returns                    ProcedureReturnsRequest   // required
	RuntimeVersion             string                    // required
	Packages                   []ProcedurePackageRequest // required
	Imports                    []ProcedureImportRequest
	Handler                    string // required
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []Secret
	TargetPath                 *string
	NullInputBehavior          *NullInputBehavior
	Comment                    *string
	ExecuteAs                  *ExecuteAs
	ProcedureDefinition        *string
}

type ProcedureArgumentRequest struct {
	ArgName      string   // required
	ArgDataType  DataType // required
	DefaultValue *string
}

type ProcedureReturnsRequest struct {
	ResultDataType *ProcedureReturnsResultDataTypeRequest
	Table          *ProcedureReturnsTableRequest
}

type ProcedureReturnsResultDataTypeRequest struct {
	ResultDataType DataType // required
	Null           *bool
	NotNull        *bool
}

type ProcedureReturnsTableRequest struct {
	Columns []ProcedureColumnRequest
}

type ProcedureColumnRequest struct {
	ColumnName     string   // required
	ColumnDataType DataType // required
}

type ProcedurePackageRequest struct {
	Package string // required
}

type ProcedureImportRequest struct {
	Import string // required
}

type CreateForJavaScriptProcedureRequest struct {
	OrReplace           *bool
	Secure              *bool
	name                SchemaObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	CopyGrants          *bool
	ResultDataType      DataType // required
	NotNull             *bool
	NullInputBehavior   *NullInputBehavior
	Comment             *string
	ExecuteAs           *ExecuteAs
	ProcedureDefinition string // required
}

type CreateForPythonProcedureRequest struct {
	OrReplace                  *bool
	Secure                     *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []ProcedureArgumentRequest
	CopyGrants                 *bool
	Returns                    ProcedureReturnsRequest   // required
	RuntimeVersion             string                    // required
	Packages                   []ProcedurePackageRequest // required
	Imports                    []ProcedureImportRequest
	Handler                    string // required
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []Secret
	NullInputBehavior          *NullInputBehavior
	Comment                    *string
	ExecuteAs                  *ExecuteAs
	ProcedureDefinition        *string
}

type CreateForScalaProcedureRequest struct {
	OrReplace           *bool
	Secure              *bool
	name                SchemaObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	CopyGrants          *bool
	Returns             ProcedureReturnsRequest   // required
	RuntimeVersion      string                    // required
	Packages            []ProcedurePackageRequest // required
	Imports             []ProcedureImportRequest
	Handler             string // required
	TargetPath          *string
	NullInputBehavior   *NullInputBehavior
	Comment             *string
	ExecuteAs           *ExecuteAs
	ProcedureDefinition *string
}

type CreateForSQLProcedureRequest struct {
	OrReplace           *bool
	Secure              *bool
	name                SchemaObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	CopyGrants          *bool
	Returns             ProcedureSQLReturnsRequest // required
	NullInputBehavior   *NullInputBehavior
	Comment             *string
	ExecuteAs           *ExecuteAs
	ProcedureDefinition string // required
}

type ProcedureSQLReturnsRequest struct {
	ResultDataType *ProcedureReturnsResultDataTypeRequest
	Table          *ProcedureReturnsTableRequest
	NotNull        *bool
}

type AlterProcedureRequest struct {
	IfExists          *bool
	name              SchemaObjectIdentifier // required
	ArgumentDataTypes []DataType             // required
	RenameTo          *SchemaObjectIdentifier
	SetComment        *string
	SetLogLevel       *string
	SetTraceLevel     *string
	UnsetComment      *bool
	SetTags           []TagAssociation
	UnsetTags         []ObjectIdentifier
	ExecuteAs         *ExecuteAs
}

type DropProcedureRequest struct {
	IfExists          *bool
	name              SchemaObjectIdentifier // required
	ArgumentDataTypes []DataType             // required
}

type ShowProcedureRequest struct {
	Like *Like
	In   *In
}

type DescribeProcedureRequest struct {
	name              SchemaObjectIdentifier // required
	ArgumentDataTypes []DataType             // required
}
