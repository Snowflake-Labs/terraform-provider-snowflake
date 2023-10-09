// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateOnTableStreamRequest(
	name AccountObjectIdentifier,
	TableId AccountObjectIdentifier,
) *CreateOnTableStreamRequest {
	s := CreateOnTableStreamRequest{}
	s.name = name
	s.TableId = TableId
	return &s
}

func (s *CreateOnTableStreamRequest) WithOrReplace(OrReplace *bool) *CreateOnTableStreamRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateOnTableStreamRequest) WithIfNotExists(IfNotExists *bool) *CreateOnTableStreamRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateOnTableStreamRequest) WithCopyGrants(CopyGrants *bool) *CreateOnTableStreamRequest {
	s.CopyGrants = CopyGrants
	return s
}

func (s *CreateOnTableStreamRequest) WithOn(On *OnStreamRequest) *CreateOnTableStreamRequest {
	s.On = On
	return s
}

func (s *CreateOnTableStreamRequest) WithAppendOnly(AppendOnly *bool) *CreateOnTableStreamRequest {
	s.AppendOnly = AppendOnly
	return s
}

func (s *CreateOnTableStreamRequest) WithShowInitialRows(ShowInitialRows *bool) *CreateOnTableStreamRequest {
	s.ShowInitialRows = ShowInitialRows
	return s
}

func (s *CreateOnTableStreamRequest) WithComment(Comment *string) *CreateOnTableStreamRequest {
	s.Comment = Comment
	return s
}

func NewOnStreamRequest() *OnStreamRequest {
	return &OnStreamRequest{}
}

func (s *OnStreamRequest) WithAt(At *bool) *OnStreamRequest {
	s.At = At
	return s
}

func (s *OnStreamRequest) WithBefore(Before *bool) *OnStreamRequest {
	s.Before = Before
	return s
}

func (s *OnStreamRequest) WithTimestamp(Timestamp *string) *OnStreamRequest {
	s.Timestamp = Timestamp
	return s
}

func (s *OnStreamRequest) WithOffset(Offset *string) *OnStreamRequest {
	s.Offset = Offset
	return s
}

func (s *OnStreamRequest) WithStatement(Statement *string) *OnStreamRequest {
	s.Statement = Statement
	return s
}

func (s *OnStreamRequest) WithStream(Stream *string) *OnStreamRequest {
	s.Stream = Stream
	return s
}

func NewCreateOnExternalTableStreamRequest(
	name AccountObjectIdentifier,
	ExternalTableId AccountObjectIdentifier,
) *CreateOnExternalTableStreamRequest {
	s := CreateOnExternalTableStreamRequest{}
	s.name = name
	s.ExternalTableId = ExternalTableId
	return &s
}

func (s *CreateOnExternalTableStreamRequest) WithOrReplace(OrReplace *bool) *CreateOnExternalTableStreamRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateOnExternalTableStreamRequest) WithIfNotExists(IfNotExists *bool) *CreateOnExternalTableStreamRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateOnExternalTableStreamRequest) WithCopyGrants(CopyGrants *bool) *CreateOnExternalTableStreamRequest {
	s.CopyGrants = CopyGrants
	return s
}

func (s *CreateOnExternalTableStreamRequest) WithOn(On *OnStreamRequest) *CreateOnExternalTableStreamRequest {
	s.On = On
	return s
}

func (s *CreateOnExternalTableStreamRequest) WithInsertOnly(InsertOnly *bool) *CreateOnExternalTableStreamRequest {
	s.InsertOnly = InsertOnly
	return s
}

func (s *CreateOnExternalTableStreamRequest) WithComment(Comment *string) *CreateOnExternalTableStreamRequest {
	s.Comment = Comment
	return s
}

func NewCreateOnStageStreamRequest(
	name AccountObjectIdentifier,
	StageId AccountObjectIdentifier,
) *CreateOnStageStreamRequest {
	s := CreateOnStageStreamRequest{}
	s.name = name
	s.StageId = StageId
	return &s
}

func (s *CreateOnStageStreamRequest) WithOrReplace(OrReplace *bool) *CreateOnStageStreamRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateOnStageStreamRequest) WithIfNotExists(IfNotExists *bool) *CreateOnStageStreamRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateOnStageStreamRequest) WithCopyGrants(CopyGrants *bool) *CreateOnStageStreamRequest {
	s.CopyGrants = CopyGrants
	return s
}

func (s *CreateOnStageStreamRequest) WithComment(Comment *string) *CreateOnStageStreamRequest {
	s.Comment = Comment
	return s
}

func NewCreateOnViewStreamRequest(
	name AccountObjectIdentifier,
	ViewId AccountObjectIdentifier,
) *CreateOnViewStreamRequest {
	s := CreateOnViewStreamRequest{}
	s.name = name
	s.ViewId = ViewId
	return &s
}

func (s *CreateOnViewStreamRequest) WithOrReplace(OrReplace *bool) *CreateOnViewStreamRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateOnViewStreamRequest) WithIfNotExists(IfNotExists *bool) *CreateOnViewStreamRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateOnViewStreamRequest) WithCopyGrants(CopyGrants *bool) *CreateOnViewStreamRequest {
	s.CopyGrants = CopyGrants
	return s
}

func (s *CreateOnViewStreamRequest) WithOn(On *OnStreamRequest) *CreateOnViewStreamRequest {
	s.On = On
	return s
}

func (s *CreateOnViewStreamRequest) WithAppendOnly(AppendOnly *bool) *CreateOnViewStreamRequest {
	s.AppendOnly = AppendOnly
	return s
}

func (s *CreateOnViewStreamRequest) WithShowInitialRows(ShowInitialRows *bool) *CreateOnViewStreamRequest {
	s.ShowInitialRows = ShowInitialRows
	return s
}

func (s *CreateOnViewStreamRequest) WithComment(Comment *string) *CreateOnViewStreamRequest {
	s.Comment = Comment
	return s
}

func NewCopyStreamRequest(
	name AccountObjectIdentifier,
	sourceStream *AccountObjectIdentifier,
) *CopyStreamRequest {
	s := CopyStreamRequest{}
	s.name = name
	s.sourceStream = sourceStream
	return &s
}

func (s *CopyStreamRequest) WithOrReplace(OrReplace *bool) *CopyStreamRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CopyStreamRequest) WithCopyGrants(CopyGrants *bool) *CopyStreamRequest {
	s.CopyGrants = CopyGrants
	return s
}

func NewAlterStreamRequest(
	name AccountObjectIdentifier,
) *AlterStreamRequest {
	s := AlterStreamRequest{}
	s.name = name
	return &s
}

func (s *AlterStreamRequest) WithIfExists(IfExists *bool) *AlterStreamRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterStreamRequest) WithSetComment(SetComment *string) *AlterStreamRequest {
	s.SetComment = SetComment
	return s
}

func (s *AlterStreamRequest) WithUnsetComment(UnsetComment *bool) *AlterStreamRequest {
	s.UnsetComment = UnsetComment
	return s
}

func (s *AlterStreamRequest) WithSetTags(SetTags []TagAssociation) *AlterStreamRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterStreamRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterStreamRequest {
	s.UnsetTags = UnsetTags
	return s
}

func NewDropStreamRequest(
	name AccountObjectIdentifier,
) *DropStreamRequest {
	s := DropStreamRequest{}
	s.name = name
	return &s
}

func (s *DropStreamRequest) WithIfExists(IfExists *bool) *DropStreamRequest {
	s.IfExists = IfExists
	return s
}

func NewShowStreamRequest() *ShowStreamRequest {
	return &ShowStreamRequest{}
}

func (s *ShowStreamRequest) WithTerse(Terse *bool) *ShowStreamRequest {
	s.Terse = Terse
	return s
}

func (s *ShowStreamRequest) WithLike(Like *Like) *ShowStreamRequest {
	s.Like = Like
	return s
}

func (s *ShowStreamRequest) WithIn(In *In) *ShowStreamRequest {
	s.In = In
	return s
}

func (s *ShowStreamRequest) WithStartsWith(StartsWith *string) *ShowStreamRequest {
	s.StartsWith = StartsWith
	return s
}

func (s *ShowStreamRequest) WithLimit(Limit *LimitFrom) *ShowStreamRequest {
	s.Limit = Limit
	return s
}

func NewDescribeStreamRequest(
	name AccountObjectIdentifier,
) *DescribeStreamRequest {
	s := DescribeStreamRequest{}
	s.name = name
	return &s
}
