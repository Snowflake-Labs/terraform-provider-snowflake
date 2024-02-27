package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateExternalFunctionOptions]   = new(CreateExternalFunctionRequest)
	_ optionsProvider[AlterExternalFunctionOptions]    = new(AlterExternalFunctionRequest)
	_ optionsProvider[ShowExternalFunctionOptions]     = new(ShowExternalFunctionRequest)
	_ optionsProvider[DescribeExternalFunctionOptions] = new(DescribeExternalFunctionRequest)
)

type CreateExternalFunctionRequest struct {
	OrReplace             *bool
	Secure                *bool
	name                  SchemaObjectIdentifier // required
	Arguments             []ExternalFunctionArgumentRequest
	ResultDataType        DataType // required
	ReturnNullValues      *ReturnNullValues
	NullInputBehavior     *NullInputBehavior
	ReturnResultsBehavior *ReturnResultsBehavior
	Comment               *string
	ApiIntegration        *AccountObjectIdentifier // required
	Headers               []ExternalFunctionHeaderRequest
	ContextHeaders        []ExternalFunctionContextHeaderRequest
	MaxBatchRows          *int
	Compression           *string
	RequestTranslator     *SchemaObjectIdentifier
	ResponseTranslator    *SchemaObjectIdentifier
	As                    string // required
}

type ExternalFunctionArgumentRequest struct {
	ArgName     string   // required
	ArgDataType DataType // required
}

type ExternalFunctionHeaderRequest struct {
	Name  string // required
	Value string // required
}

type ExternalFunctionContextHeaderRequest struct {
	ContextFunction string // required
}

type AlterExternalFunctionRequest struct {
	IfExists          *bool
	name              SchemaObjectIdentifier // required
	ArgumentDataTypes []DataType             // required
	Set               *ExternalFunctionSetRequest
	Unset             *ExternalFunctionUnsetRequest
}

type ExternalFunctionSetRequest struct {
	ApiIntegration     *AccountObjectIdentifier
	Headers            []ExternalFunctionHeaderRequest
	ContextHeaders     []ExternalFunctionContextHeaderRequest
	MaxBatchRows       *int
	Compression        *string
	RequestTranslator  *SchemaObjectIdentifier
	ResponseTranslator *SchemaObjectIdentifier
}

type ExternalFunctionUnsetRequest struct {
	Comment            *bool
	Headers            *bool
	ContextHeaders     *bool
	MaxBatchRows       *bool
	Compression        *bool
	Secure             *bool
	RequestTranslator  *bool
	ResponseTranslator *bool
}

type ShowExternalFunctionRequest struct {
	Like *Like
	In   *In
}

type DescribeExternalFunctionRequest struct {
	name              SchemaObjectIdentifier // required
	ArgumentDataTypes []DataType             // required
}
