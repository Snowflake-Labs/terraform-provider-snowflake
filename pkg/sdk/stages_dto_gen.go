package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateInternalStageOptions]          = new(CreateInternalStageRequest)
	_ optionsProvider[CreateOnS3StageOptions]              = new(CreateOnS3StageRequest)
	_ optionsProvider[CreateOnGCSStageOptions]             = new(CreateOnGCSStageRequest)
	_ optionsProvider[CreateOnAzureStageOptions]           = new(CreateOnAzureStageRequest)
	_ optionsProvider[CreateOnS3CompatibleStageOptions]    = new(CreateOnS3CompatibleStageRequest)
	_ optionsProvider[AlterStageOptions]                   = new(AlterStageRequest)
	_ optionsProvider[AlterInternalStageStageOptions]      = new(AlterInternalStageStageRequest)
	_ optionsProvider[AlterExternalS3StageStageOptions]    = new(AlterExternalS3StageStageRequest)
	_ optionsProvider[AlterExternalGCSStageStageOptions]   = new(AlterExternalGCSStageStageRequest)
	_ optionsProvider[AlterExternalAzureStageStageOptions] = new(AlterExternalAzureStageStageRequest)
	_ optionsProvider[AlterDirectoryTableStageOptions]     = new(AlterDirectoryTableStageRequest)
	_ optionsProvider[DropStageOptions]                    = new(DropStageRequest)
	_ optionsProvider[DescribeStageOptions]                = new(DescribeStageRequest)
	_ optionsProvider[ShowStageOptions]                    = new(ShowStageRequest)
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
	Type *InternalStageEncryptionOption // required
}

type InternalDirectoryTableOptionsRequest struct {
	Enable          *bool
	RefreshOnCreate *bool
}

type StageFileFormatRequest struct {
	FormatName *string
	Type       *FileFormatType
	Options    *FileFormatTypeOptionsRequest
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
}

type StageCopyOnErrorOptionsRequest struct {
	Continue       *bool
	SkipFile       *string
	AbortStatement *bool
}

type CreateOnS3StageRequest struct {
	OrReplace             *bool
	Temporary             *bool
	IfNotExists           *bool
	name                  SchemaObjectIdentifier // required
	ExternalStageParams   *ExternalS3StageParamsRequest
	DirectoryTableOptions *ExternalS3DirectoryTableOptionsRequest
	FileFormat            *StageFileFormatRequest
	CopyOptions           *StageCopyOptionsRequest
	Comment               *string
	Tag                   []TagAssociation
}

type ExternalS3StageParamsRequest struct {
	Url                string // required
	StorageIntegration *AccountObjectIdentifier
	Credentials        *ExternalStageS3CredentialsRequest
	Encryption         *ExternalStageS3EncryptionRequest
}

type ExternalStageS3CredentialsRequest struct {
	AWSKeyId     *string
	AWSSecretKey *string
	AWSToken     *string
	AWSRole      *string
}

type ExternalStageS3EncryptionRequest struct {
	Type      *ExternalStageS3EncryptionOption // required
	MasterKey *string
	KmsKeyId  *string
}

type ExternalS3DirectoryTableOptionsRequest struct {
	Enable          *bool
	RefreshOnCreate *bool
	AutoRefresh     *bool
}

type CreateOnGCSStageRequest struct {
	OrReplace             *bool
	Temporary             *bool
	IfNotExists           *bool
	name                  SchemaObjectIdentifier // required
	ExternalStageParams   *ExternalGCSStageParamsRequest
	DirectoryTableOptions *ExternalGCSDirectoryTableOptionsRequest
	FileFormat            *StageFileFormatRequest
	CopyOptions           *StageCopyOptionsRequest
	Comment               *string
	Tag                   []TagAssociation
}

type ExternalGCSStageParamsRequest struct {
	Url                string // required
	StorageIntegration *AccountObjectIdentifier
	Encryption         *ExternalStageGCSEncryptionRequest
}

type ExternalStageGCSEncryptionRequest struct {
	Type     *ExternalStageGCSEncryptionOption // required
	KmsKeyId *string
}

type ExternalGCSDirectoryTableOptionsRequest struct {
	Enable                  *bool
	RefreshOnCreate         *bool
	AutoRefresh             *bool
	NotificationIntegration *string
}

type CreateOnAzureStageRequest struct {
	OrReplace             *bool
	Temporary             *bool
	IfNotExists           *bool
	name                  SchemaObjectIdentifier // required
	ExternalStageParams   *ExternalAzureStageParamsRequest
	DirectoryTableOptions *ExternalAzureDirectoryTableOptionsRequest
	FileFormat            *StageFileFormatRequest
	CopyOptions           *StageCopyOptionsRequest
	Comment               *string
	Tag                   []TagAssociation
}

type ExternalAzureStageParamsRequest struct {
	Url                string // required
	StorageIntegration *AccountObjectIdentifier
	Credentials        *ExternalStageAzureCredentialsRequest
	Encryption         *ExternalStageAzureEncryptionRequest
}

type ExternalStageAzureCredentialsRequest struct {
	AzureSasToken string // required
}

type ExternalStageAzureEncryptionRequest struct {
	Type      *ExternalStageAzureEncryptionOption // required
	MasterKey *string
}

type ExternalAzureDirectoryTableOptionsRequest struct {
	Enable                  *bool
	RefreshOnCreate         *bool
	AutoRefresh             *bool
	NotificationIntegration *string
}

type CreateOnS3CompatibleStageRequest struct {
	OrReplace             *bool
	Temporary             *bool
	IfNotExists           *bool
	name                  SchemaObjectIdentifier // required
	Url                   string                 // required
	Endpoint              string                 // required
	Credentials           *ExternalStageS3CompatibleCredentialsRequest
	DirectoryTableOptions *ExternalS3DirectoryTableOptionsRequest
	FileFormat            *StageFileFormatRequest
	CopyOptions           *StageCopyOptionsRequest
	Comment               *string
	Tag                   []TagAssociation
}

type ExternalStageS3CompatibleCredentialsRequest struct {
	AWSKeyId     *string // required
	AWSSecretKey *string // required
}

type AlterStageRequest struct {
	IfExists  *bool
	name      SchemaObjectIdentifier // required
	RenameTo  *SchemaObjectIdentifier
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
}

type AlterInternalStageStageRequest struct {
	IfExists    *bool
	name        SchemaObjectIdentifier // required
	FileFormat  *StageFileFormatRequest
	CopyOptions *StageCopyOptionsRequest
	Comment     *string
}

type AlterExternalS3StageStageRequest struct {
	IfExists            *bool
	name                SchemaObjectIdentifier // required
	ExternalStageParams *ExternalS3StageParamsRequest
	FileFormat          *StageFileFormatRequest
	CopyOptions         *StageCopyOptionsRequest
	Comment             *string
}

type AlterExternalGCSStageStageRequest struct {
	IfExists            *bool
	name                SchemaObjectIdentifier // required
	ExternalStageParams *ExternalGCSStageParamsRequest
	FileFormat          *StageFileFormatRequest
	CopyOptions         *StageCopyOptionsRequest
	Comment             *string
}

type AlterExternalAzureStageStageRequest struct {
	IfExists            *bool
	name                SchemaObjectIdentifier // required
	ExternalStageParams *ExternalAzureStageParamsRequest
	FileFormat          *StageFileFormatRequest
	CopyOptions         *StageCopyOptionsRequest
	Comment             *string
}

type AlterDirectoryTableStageRequest struct {
	IfExists     *bool
	name         SchemaObjectIdentifier // required
	SetDirectory *DirectoryTableSetRequest
	Refresh      *DirectoryTableRefreshRequest
}

type DirectoryTableSetRequest struct {
	Enable bool // required
}

type DirectoryTableRefreshRequest struct {
	Subpath *string
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
