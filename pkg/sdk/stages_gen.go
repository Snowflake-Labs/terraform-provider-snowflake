package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Stages interface {
	CreateInternal(ctx context.Context, request *CreateInternalStageRequest) error
	CreateOnS3(ctx context.Context, request *CreateOnS3StageRequest) error
	CreateOnGCS(ctx context.Context, request *CreateOnGCSStageRequest) error
	CreateOnAzure(ctx context.Context, request *CreateOnAzureStageRequest) error
	CreateOnS3Compatible(ctx context.Context, request *CreateOnS3CompatibleStageRequest) error
	Alter(ctx context.Context, request *AlterStageRequest) error
	AlterInternalStage(ctx context.Context, request *AlterInternalStageStageRequest) error
	AlterExternalS3Stage(ctx context.Context, request *AlterExternalS3StageStageRequest) error
	AlterExternalGCSStage(ctx context.Context, request *AlterExternalGCSStageStageRequest) error
	AlterExternalAzureStage(ctx context.Context, request *AlterExternalAzureStageStageRequest) error
	AlterDirectoryTable(ctx context.Context, request *AlterDirectoryTableStageRequest) error
	Drop(ctx context.Context, request *DropStageRequest) error
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]StageProperty, error)
	Show(ctx context.Context, request *ShowStageRequest) ([]Stage, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Stage, error)
}

// CreateInternalStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stage.
type CreateInternalStageOptions struct {
	create                bool                           `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                          `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                          `ddl:"keyword" sql:"TEMPORARY"`
	stage                 bool                           `ddl:"static" sql:"STAGE"`
	IfNotExists           *bool                          `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier         `ddl:"identifier"`
	Encryption            *InternalStageEncryption       `ddl:"list,parentheses,no_comma" sql:"ENCRYPTION ="`
	DirectoryTableOptions *InternalDirectoryTableOptions `ddl:"list,parentheses,no_comma" sql:"DIRECTORY ="`
	FileFormat            *StageFileFormat               `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions           *StageCopyOptions              `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment               *string                        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag                   []TagAssociation               `ddl:"keyword,parentheses" sql:"TAG"`
}

type InternalStageEncryption struct {
	Type *InternalStageEncryptionOption `ddl:"parameter,single_quotes" sql:"TYPE"`
}

type InternalDirectoryTableOptions struct {
	Enable          *bool `ddl:"parameter" sql:"ENABLE"`
	RefreshOnCreate *bool `ddl:"parameter" sql:"REFRESH_ON_CREATE"`
}

type StageFileFormat struct {
	FormatName *string                `ddl:"parameter,single_quotes" sql:"FORMAT_NAME"`
	Type       *FileFormatType        `ddl:"parameter" sql:"TYPE"`
	Options    *FileFormatTypeOptions `ddl:"list,no_comma"`
}

type StageCopyOptions struct {
	OnError           *StageCopyOnErrorOptions  `ddl:"parameter" sql:"ON_ERROR"`
	SizeLimit         *int                      `ddl:"parameter" sql:"SIZE_LIMIT"`
	Purge             *bool                     `ddl:"parameter" sql:"PURGE"`
	ReturnFailedOnly  *bool                     `ddl:"parameter" sql:"RETURN_FAILED_ONLY"`
	MatchByColumnName *StageCopyColumnMapOption `ddl:"parameter" sql:"MATCH_BY_COLUMN_NAME"`
	EnforceLength     *bool                     `ddl:"parameter" sql:"ENFORCE_LENGTH"`
	Truncatecolumns   *bool                     `ddl:"parameter" sql:"TRUNCATECOLUMNS"`
	Force             *bool                     `ddl:"parameter" sql:"FORCE"`
}

type StageCopyOnErrorOptions struct {
	Continue       *bool   `ddl:"keyword" sql:"CONTINUE"`
	SkipFile       *string `ddl:"keyword"`
	AbortStatement *bool   `ddl:"keyword" sql:"ABORT_STATEMENT"`
}

// CreateOnS3StageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stage.
type CreateOnS3StageOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	stage                 bool                   `ddl:"static" sql:"STAGE"`
	IfNotExists           *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	ExternalStageParams   *ExternalS3StageParams
	DirectoryTableOptions *ExternalS3DirectoryTableOptions `ddl:"list,parentheses,no_comma" sql:"DIRECTORY ="`
	FileFormat            *StageFileFormat                 `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions           *StageCopyOptions                `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment               *string                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag                   []TagAssociation                 `ddl:"keyword,parentheses" sql:"TAG"`
}

type ExternalS3StageParams struct {
	Url                string                      `ddl:"parameter,single_quotes" sql:"URL"`
	StorageIntegration *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"STORAGE_INTEGRATION"`
	Credentials        *ExternalStageS3Credentials `ddl:"list,parentheses,no_comma" sql:"CREDENTIALS ="`
	Encryption         *ExternalStageS3Encryption  `ddl:"list,parentheses,no_comma" sql:"ENCRYPTION ="`
}

type ExternalStageS3Credentials struct {
	AWSKeyId     *string `ddl:"parameter,single_quotes" sql:"AWS_KEY_ID"`
	AWSSecretKey *string `ddl:"parameter,single_quotes" sql:"AWS_SECRET_KEY"`
	AWSToken     *string `ddl:"parameter,single_quotes" sql:"AWS_TOKEN"`
	AWSRole      *string `ddl:"parameter,single_quotes" sql:"AWS_ROLE"`
}

type ExternalStageS3Encryption struct {
	Type      *ExternalStageS3EncryptionOption `ddl:"parameter,single_quotes" sql:"TYPE"`
	MasterKey *string                          `ddl:"parameter,single_quotes" sql:"MASTER_KEY"`
	KmsKeyId  *string                          `ddl:"parameter,single_quotes" sql:"KMS_KEY_ID"`
}

type ExternalS3DirectoryTableOptions struct {
	Enable          *bool `ddl:"parameter" sql:"ENABLE"`
	RefreshOnCreate *bool `ddl:"parameter" sql:"REFRESH_ON_CREATE"`
	AutoRefresh     *bool `ddl:"parameter" sql:"AUTO_REFRESH"`
}

// CreateOnGCSStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stage.
type CreateOnGCSStageOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	stage                 bool                   `ddl:"static" sql:"STAGE"`
	IfNotExists           *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	ExternalStageParams   *ExternalGCSStageParams
	DirectoryTableOptions *ExternalGCSDirectoryTableOptions `ddl:"list,parentheses,no_comma" sql:"DIRECTORY ="`
	FileFormat            *StageFileFormat                  `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions           *StageCopyOptions                 `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment               *string                           `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag                   []TagAssociation                  `ddl:"keyword,parentheses" sql:"TAG"`
}

type ExternalGCSStageParams struct {
	Url                string                      `ddl:"parameter,single_quotes" sql:"URL"`
	StorageIntegration *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"STORAGE_INTEGRATION"`
	Encryption         *ExternalStageGCSEncryption `ddl:"list,parentheses,no_comma" sql:"ENCRYPTION ="`
}

type ExternalStageGCSEncryption struct {
	Type     *ExternalStageGCSEncryptionOption `ddl:"parameter,single_quotes" sql:"TYPE"`
	KmsKeyId *string                           `ddl:"parameter,single_quotes" sql:"KMS_KEY_ID"`
}

type ExternalGCSDirectoryTableOptions struct {
	Enable                  *bool   `ddl:"parameter" sql:"ENABLE"`
	RefreshOnCreate         *bool   `ddl:"parameter" sql:"REFRESH_ON_CREATE"`
	AutoRefresh             *bool   `ddl:"parameter" sql:"AUTO_REFRESH"`
	NotificationIntegration *string `ddl:"parameter,single_quotes" sql:"NOTIFICATION_INTEGRATION"`
}

// CreateOnAzureStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stage.
type CreateOnAzureStageOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	stage                 bool                   `ddl:"static" sql:"STAGE"`
	IfNotExists           *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	ExternalStageParams   *ExternalAzureStageParams
	DirectoryTableOptions *ExternalAzureDirectoryTableOptions `ddl:"list,parentheses,no_comma" sql:"DIRECTORY ="`
	FileFormat            *StageFileFormat                    `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions           *StageCopyOptions                   `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment               *string                             `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag                   []TagAssociation                    `ddl:"keyword,parentheses" sql:"TAG"`
}

type ExternalAzureStageParams struct {
	Url                string                         `ddl:"parameter,single_quotes" sql:"URL"`
	StorageIntegration *AccountObjectIdentifier       `ddl:"identifier,equals" sql:"STORAGE_INTEGRATION"`
	Credentials        *ExternalStageAzureCredentials `ddl:"list,parentheses,no_comma" sql:"CREDENTIALS ="`
	Encryption         *ExternalStageAzureEncryption  `ddl:"list,parentheses,no_comma" sql:"ENCRYPTION ="`
}

type ExternalStageAzureCredentials struct {
	AzureSasToken string `ddl:"parameter,single_quotes" sql:"AZURE_SAS_TOKEN"`
}

type ExternalStageAzureEncryption struct {
	Type      *ExternalStageAzureEncryptionOption `ddl:"parameter,single_quotes" sql:"TYPE"`
	MasterKey *string                             `ddl:"parameter,single_quotes" sql:"MASTER_KEY"`
}

type ExternalAzureDirectoryTableOptions struct {
	Enable                  *bool   `ddl:"parameter" sql:"ENABLE"`
	RefreshOnCreate         *bool   `ddl:"parameter" sql:"REFRESH_ON_CREATE"`
	AutoRefresh             *bool   `ddl:"parameter" sql:"AUTO_REFRESH"`
	NotificationIntegration *string `ddl:"parameter,single_quotes" sql:"NOTIFICATION_INTEGRATION"`
}

// CreateOnS3CompatibleStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stage.
type CreateOnS3CompatibleStageOptions struct {
	create                bool                                  `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                                 `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                                 `ddl:"keyword" sql:"TEMPORARY"`
	stage                 bool                                  `ddl:"static" sql:"STAGE"`
	IfNotExists           *bool                                 `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier                `ddl:"identifier"`
	Url                   string                                `ddl:"parameter,single_quotes" sql:"URL"`
	Endpoint              string                                `ddl:"parameter,single_quotes" sql:"ENDPOINT"`
	Credentials           *ExternalStageS3CompatibleCredentials `ddl:"list,parentheses,no_comma" sql:"CREDENTIALS ="`
	DirectoryTableOptions *ExternalS3DirectoryTableOptions      `ddl:"list,parentheses,no_comma" sql:"DIRECTORY ="`
	FileFormat            *StageFileFormat                      `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions           *StageCopyOptions                     `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment               *string                               `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag                   []TagAssociation                      `ddl:"keyword,parentheses" sql:"TAG"`
}

type ExternalStageS3CompatibleCredentials struct {
	AWSKeyId     *string `ddl:"parameter,single_quotes" sql:"AWS_KEY_ID"`
	AWSSecretKey *string `ddl:"parameter,single_quotes" sql:"AWS_SECRET_KEY"`
}

// AlterStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-stage.
type AlterStageOptions struct {
	alter     bool                    `ddl:"static" sql:"ALTER"`
	stage     bool                    `ddl:"static" sql:"STAGE"`
	IfExists  *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name      SchemaObjectIdentifier  `ddl:"identifier"`
	RenameTo  *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SetTags   []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

// AlterInternalStageStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-stage.
type AlterInternalStageStageOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"`
	stage    bool                   `ddl:"static" sql:"STAGE"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
	set      bool                   `ddl:"static" sql:"SET"`
	// TODO (SNOW-1019005): Move parameters below to the AlterStageOptions as they're common across stage types (+ remove from other alter option structs)
	FileFormat  *StageFileFormat  `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions *StageCopyOptions `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment     *string           `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterExternalS3StageStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-stage.
type AlterExternalS3StageStageOptions struct {
	alter               bool                   `ddl:"static" sql:"ALTER"`
	stage               bool                   `ddl:"static" sql:"STAGE"`
	IfExists            *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	set                 bool                   `ddl:"static" sql:"SET"`
	ExternalStageParams *ExternalS3StageParams
	FileFormat          *StageFileFormat  `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions         *StageCopyOptions `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment             *string           `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterExternalGCSStageStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-stage.
type AlterExternalGCSStageStageOptions struct {
	alter               bool                   `ddl:"static" sql:"ALTER"`
	stage               bool                   `ddl:"static" sql:"STAGE"`
	IfExists            *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	set                 bool                   `ddl:"static" sql:"SET"`
	ExternalStageParams *ExternalGCSStageParams
	FileFormat          *StageFileFormat  `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions         *StageCopyOptions `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment             *string           `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterExternalAzureStageStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-stage.
type AlterExternalAzureStageStageOptions struct {
	alter               bool                   `ddl:"static" sql:"ALTER"`
	stage               bool                   `ddl:"static" sql:"STAGE"`
	IfExists            *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	set                 bool                   `ddl:"static" sql:"SET"`
	ExternalStageParams *ExternalAzureStageParams
	FileFormat          *StageFileFormat  `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyOptions         *StageCopyOptions `ddl:"list,parentheses,no_comma" sql:"COPY_OPTIONS ="`
	Comment             *string           `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterDirectoryTableStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-stage.
type AlterDirectoryTableStageOptions struct {
	alter        bool                   `ddl:"static" sql:"ALTER"`
	stage        bool                   `ddl:"static" sql:"STAGE"`
	IfExists     *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
	SetDirectory *DirectoryTableSet     `ddl:"list,parentheses,no_comma" sql:"SET DIRECTORY ="`
	Refresh      *DirectoryTableRefresh `ddl:"keyword" sql:"REFRESH"`
}

type DirectoryTableSet struct {
	Enable bool `ddl:"parameter" sql:"ENABLE"`
}

type DirectoryTableRefresh struct {
	Subpath *string `ddl:"parameter,single_quotes" sql:"SUBPATH"`
}

// DropStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-stage.
type DropStageOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	stage    bool                   `ddl:"static" sql:"STAGE"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// DescribeStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-stage.
type DescribeStageOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	stage    bool                   `ddl:"static" sql:"STAGE"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

type stageDescRow struct {
	ParentProperty  string `db:"parent_property"`
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

type StageProperty struct {
	Parent  string
	Name    string
	Type    string
	Value   string
	Default string
}

// ShowStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-stages.
type ShowStageOptions struct {
	show   bool  `ddl:"static" sql:"SHOW"`
	stages bool  `ddl:"static" sql:"STAGES"`
	Like   *Like `ddl:"keyword" sql:"LIKE"`
	In     *In   `ddl:"keyword" sql:"IN"`
}

type stageShowRow struct {
	CreatedOn          time.Time      `db:"created_on"`
	Name               string         `db:"name"`
	DatabaseName       string         `db:"database_name"`
	SchemaName         string         `db:"schema_name"`
	Url                string         `db:"url"`
	HasCredentials     string         `db:"has_credentials"`
	HasEncryptionKey   string         `db:"has_encryption_key"`
	Owner              string         `db:"owner"`
	Comment            string         `db:"comment"`
	Region             sql.NullString `db:"region"`
	Type               string         `db:"type"`
	Cloud              sql.NullString `db:"cloud"`
	StorageIntegration sql.NullString `db:"storage_integration"`
	Endpoint           sql.NullString `db:"endpoint"`
	OwnerRoleType      sql.NullString `db:"owner_role_type"`
	DirectoryEnabled   string         `db:"directory_enabled"`
}

type Stage struct {
	CreatedOn          time.Time
	Name               string
	DatabaseName       string
	SchemaName         string
	Url                string
	HasCredentials     bool
	HasEncryptionKey   bool
	Owner              string
	Comment            string
	Region             *string
	Type               string
	Cloud              *string
	StorageIntegration *string
	Endpoint           *string
	OwnerRoleType      *string
	DirectoryEnabled   bool
}

func (s *Stage) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(s.DatabaseName, s.SchemaName, s.Name)
}

func (s *Stage) ObjectType() ObjectType {
	return ObjectTypeStage
}

func (s *Stage) Location() string {
	return "@" + s.ID().FullyQualifiedName()
}
