package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateStreamlitOptions]   = new(CreateStreamlitRequest)
	_ optionsProvider[AlterStreamlitOptions]    = new(AlterStreamlitRequest)
	_ optionsProvider[DropStreamlitOptions]     = new(DropStreamlitRequest)
	_ optionsProvider[ShowStreamlitOptions]     = new(ShowStreamlitRequest)
	_ optionsProvider[DescribeStreamlitOptions] = new(DescribeStreamlitRequest)
)

type CreateStreamlitRequest struct {
	OrReplace    *bool
	IfNotExists  *bool
	name         SchemaObjectIdentifier // required
	RootLocation string                 // required
	MainFile     string                 // required
	Warehouse    *AccountObjectIdentifier
	Comment      *string
}

type AlterStreamlitRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Set      *StreamlitSetRequest
	RenameTo *SchemaObjectIdentifier
}

type StreamlitSetRequest struct {
	RootLocation *string // required
	MainFile     *string // required
	Warehouse    *AccountObjectIdentifier
	Comment      *string
}

type DropStreamlitRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowStreamlitRequest struct {
	Terse *bool
	Like  *Like
	In    *In
	Limit *LimitFrom
}

type DescribeStreamlitRequest struct {
	name SchemaObjectIdentifier // required
}
