package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateOnTableStreamOptions]         = new(CreateOnTableStreamRequest)
	_ optionsProvider[CreateOnExternalTableStreamOptions] = new(CreateOnExternalTableStreamRequest)
	_ optionsProvider[CreateOnStageStreamOptions]         = new(CreateOnStageStreamRequest)
	_ optionsProvider[CreateOnViewStreamOptions]          = new(CreateOnViewStreamRequest)
	_ optionsProvider[CopyStreamOptions]                  = new(CopyStreamRequest)
	_ optionsProvider[AlterStreamOptions]                 = new(AlterStreamRequest)
	_ optionsProvider[DropStreamOptions]                  = new(DropStreamRequest)
	_ optionsProvider[ShowStreamOptions]                  = new(ShowStreamRequest)
	_ optionsProvider[DescribeStreamOptions]              = new(DescribeStreamRequest)
)

type CreateOnTableStreamRequest struct {
	OrReplace       *bool
	IfNotExists     *bool
	name            AccountObjectIdentifier // required
	CopyGrants      *bool
	TableId         AccountObjectIdentifier // required
	On              *OnStreamRequest
	AppendOnly      *bool
	ShowInitialRows *bool
	Comment         *string
}

type OnStreamRequest struct {
	At        *bool
	Before    *bool
	Timestamp *string
	Offset    *string
	Statement *string
	Stream    *string
}

type CreateOnExternalTableStreamRequest struct {
	OrReplace       *bool
	IfNotExists     *bool
	name            AccountObjectIdentifier // required
	CopyGrants      *bool
	ExternalTableId AccountObjectIdentifier // required
	On              *OnStreamRequest
	InsertOnly      *bool
	Comment         *string
}

type CreateOnStageStreamRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	CopyGrants  *bool
	StageId     AccountObjectIdentifier // required
	Comment     *string
}

type CreateOnViewStreamRequest struct {
	OrReplace       *bool
	IfNotExists     *bool
	name            AccountObjectIdentifier // required
	CopyGrants      *bool
	ViewId          AccountObjectIdentifier // required
	On              *OnStreamRequest
	AppendOnly      *bool
	ShowInitialRows *bool
	Comment         *string
}

type CopyStreamRequest struct {
	OrReplace    *bool
	name         AccountObjectIdentifier  // required
	sourceStream *AccountObjectIdentifier // required
	CopyGrants   *bool
}

type AlterStreamRequest struct {
	IfExists     *bool
	name         AccountObjectIdentifier // required
	SetComment   *string
	UnsetComment *bool
	SetTags      []TagAssociation
	UnsetTags    []ObjectIdentifier
}

type DropStreamRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowStreamRequest struct {
	Terse      *bool
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}

type ShowByIdStreamRequest struct {
	name AccountObjectIdentifier // required
}

type DescribeStreamRequest struct {
	name AccountObjectIdentifier // required
}
