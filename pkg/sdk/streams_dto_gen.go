package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateOnTableStreamOptions]          = new(CreateOnTableStreamRequest)
	_ optionsProvider[CreateOnExternalTableStreamOptions]  = new(CreateOnExternalTableStreamRequest)
	_ optionsProvider[CreateOnDirectoryTableStreamOptions] = new(CreateOnDirectoryTableStreamRequest)
	_ optionsProvider[CreateOnViewStreamOptions]           = new(CreateOnViewStreamRequest)
	_ optionsProvider[CloneStreamOptions]                  = new(CloneStreamRequest)
	_ optionsProvider[AlterStreamOptions]                  = new(AlterStreamRequest)
	_ optionsProvider[DropStreamOptions]                   = new(DropStreamRequest)
	_ optionsProvider[ShowStreamOptions]                   = new(ShowStreamRequest)
	_ optionsProvider[DescribeStreamOptions]               = new(DescribeStreamRequest)
)

type CreateOnTableStreamRequest struct {
	OrReplace       *bool
	IfNotExists     *bool
	name            SchemaObjectIdentifier // required
	Tag             []TagAssociation
	CopyGrants      *bool
	TableId         SchemaObjectIdentifier // required
	On              *OnStreamRequest
	AppendOnly      *bool
	ShowInitialRows *bool
	Comment         *string
}

func (r *CreateOnTableStreamRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type OnStreamRequest struct {
	At        *bool
	Before    *bool
	Statement OnStreamStatementRequest
}

type OnStreamStatementRequest struct {
	Timestamp *string
	Offset    *string
	Statement *string
	Stream    *string
}

type CreateOnExternalTableStreamRequest struct {
	OrReplace       *bool
	IfNotExists     *bool
	name            SchemaObjectIdentifier // required
	Tag             []TagAssociation
	CopyGrants      *bool
	ExternalTableId SchemaObjectIdentifier // required
	On              *OnStreamRequest
	InsertOnly      *bool
	Comment         *string
}

type CreateOnDirectoryTableStreamRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        SchemaObjectIdentifier // required
	Tag         []TagAssociation
	CopyGrants  *bool
	StageId     SchemaObjectIdentifier // required
	Comment     *string
}

type CreateOnViewStreamRequest struct {
	OrReplace       *bool
	IfNotExists     *bool
	name            SchemaObjectIdentifier // required
	Tag             []TagAssociation
	CopyGrants      *bool
	ViewId          SchemaObjectIdentifier // required
	On              *OnStreamRequest
	AppendOnly      *bool
	ShowInitialRows *bool
	Comment         *string
}

type CloneStreamRequest struct {
	OrReplace    *bool
	name         SchemaObjectIdentifier // required
	sourceStream SchemaObjectIdentifier // required
	CopyGrants   *bool
}

type AlterStreamRequest struct {
	IfExists     *bool
	name         SchemaObjectIdentifier // required
	SetComment   *string
	UnsetComment *bool
	SetTags      []TagAssociation
	UnsetTags    []ObjectIdentifier
}

type DropStreamRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowStreamRequest struct {
	Terse      *bool
	Like       *Like
	In         *ExtendedIn
	StartsWith *string
	Limit      *LimitFrom
}

type ShowByIdStreamRequest struct {
	name SchemaObjectIdentifier // required
}

type DescribeStreamRequest struct {
	name SchemaObjectIdentifier // required
}
