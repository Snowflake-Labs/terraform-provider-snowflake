package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Stages interface {
	CreateInternal(ctx context.Context, request *CreateInternalStageRequest) error
	Drop(ctx context.Context, request *DropStageRequest) error
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]StageProperty, error)
	Show(ctx context.Context, request *ShowStageRequest) ([]Stage, error)
}

// CreateInternalStageOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stage.
type CreateInternalStageOptions struct {
	create                bool                           `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                          `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                          `ddl:"keyword" sql:"TEMPORARY"`
	stage                 bool                           `ddl:"static" sql:"STAGE"`
	IfNotExists           *bool                          `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                  SchemaObjectIdentifier         `ddl:"identifier"`
	Encryption            *InternalStageEncryption       `ddl:"parameter,parentheses" sql:"ENCRYPTION"`
	DirectoryTableOptions *InternalDirectoryTableOptions `ddl:"parameter,parentheses"`
	FileFormat            *StageFileFormat               `ddl:"parameter,parentheses" sql:"FILE_FORMAT"`
	CopyOptions           *StageCopyOptions              `ddl:"parameter,parentheses" sql:"FILE_FORMAT"`
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
	FormatName *string          `ddl:"parameter,single_quotes" sql:"FORMAT_NAME"`
	Type       *FileFormatType  `ddl:"parameter" sql:"TYPE"`
	TYPE       []FileFormatType `ddl:"list"`
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
	Type              *FileFormatType           `ddl:"parameter" sql:"TYPE"`
	TYPE              []FileFormatType          `ddl:"list"`
}

type StageCopyOnErrorOptions struct {
	Continue *bool `ddl:"keyword" sql:"CONTINUE"`
	SkipFile *bool `ddl:"keyword" sql:"SKIP_FILE"`
	SkipFile *bool `ddl:"keyword" sql:"SKIP_FILE"`
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
	ParentProperty  string         `db:"parent_property"`
	Property        string         `db:"property"`
	PropertyType    string         `db:"property_type"`
	PropertyValue   sql.NullString `db:"property_value"`
	PropertyDefault sql.NullString `db:"property_default"`
}

type StageProperty struct {
	Parent  string
	Name    string
	Type    string
	Value   *string
	Default *string
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
	Region             string         `db:"region"`
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
	Region             string
	Type               string
	Cloud              *string
	StorageIntegration *string
	Endpoint           *string
	OwnerRoleType      *string
	DirectoryEnabled   bool
}
