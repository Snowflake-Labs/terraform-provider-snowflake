package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateInternalStageOptions] = new(CreateInternalStageRequest)
	_ optionsProvider[DropStageOptions]           = new(DropStageRequest)
	_ optionsProvider[DescribeStageOptions]       = new(DescribeStageRequest)
	_ optionsProvider[ShowStageOptions]           = new(ShowStageRequest)
)

type CreateInternalStageRequest struct {
	OrReplace             *bool
	Temporary             *bool
	IfNotExists           *bool
	name                  SchemaObjectIdentifier // required
	Encryption            *InternalStageEncryptionRequest
	DirectoryTableOptions *InternalDirectoryTableOptionsRequest
	FileFormat            *StageFileFormatRequest
	CopyOptions           *StageCopyOptionsRequest
	Comment               *string
	Tag                   []TagAssociation
}

type InternalStageEncryptionRequest struct {
	Type *InternalStageEncryptionOption
}

type InternalDirectoryTableOptionsRequest struct {
	Enable          *bool
	RefreshOnCreate *bool
}

type StageFileFormatRequest struct {
	FormatName *string
	Type       *FileFormatType
	TYPE       []FileFormatType
}

type StageCopyOptionsRequest struct {
	OnError           *StageCopyOnErrorOptionsRequest
	SizeLimit         *int
	Purge             *bool
	ReturnFailedOnly  *bool
	MatchByColumnName *StageCopyColumnMapOption
	EnforceLength     *bool
	Truncatecolumns   *bool
	Force             *bool
	Type              *FileFormatType
	TYPE              []FileFormatType
}

type StageCopyOnErrorOptionsRequest struct {
	Continue *bool
	SkipFile *bool
	SkipFile *bool
}

type DropStageRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type DescribeStageRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowStageRequest struct {
	Like *Like
	In   *In
}
