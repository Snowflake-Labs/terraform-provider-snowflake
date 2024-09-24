package sdk

import (
	"context"
	"database/sql"
)

type ExternalVolumes interface {
	Create(ctx context.Context, request *CreateExternalVolumeRequest) error
	Alter(ctx context.Context, request *AlterExternalVolumeRequest) error
	Drop(ctx context.Context, request *DropExternalVolumeRequest) error
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]ExternalVolumeProperty, error)
	Show(ctx context.Context, request *ShowExternalVolumeRequest) ([]ExternalVolume, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ExternalVolume, error)
}

// CreateExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-external-volume.
type CreateExternalVolumeOptions struct {
	create           bool                            `ddl:"static" sql:"CREATE"`
	OrReplace        *bool                           `ddl:"keyword" sql:"OR REPLACE"`
	externalVolume   bool                            `ddl:"static" sql:"EXTERNAL VOLUME"`
	IfNotExists      *bool                           `ddl:"keyword" sql:"IF NOT EXISTS"`
	name             AccountObjectIdentifier         `ddl:"identifier"`
	StorageLocations []ExternalVolumeStorageLocation `ddl:"parameter,parentheses" sql:"STORAGE_LOCATIONS"`
	AllowWrites      *bool                           `ddl:"parameter" sql:"ALLOW_WRITES"`
	Comment          *string                         `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type ExternalVolumeStorageLocation struct {
	S3StorageLocationParams    *S3StorageLocationParams    `ddl:"list,parentheses,no_comma"`
	GCSStorageLocationParams   *GCSStorageLocationParams   `ddl:"list,parentheses,no_comma"`
	AzureStorageLocationParams *AzureStorageLocationParams `ddl:"list,parentheses,no_comma"`
}
type S3StorageLocationParams struct {
	Name                 string                      `ddl:"parameter,single_quotes" sql:"NAME"`
	StorageProvider      S3StorageProvider           `ddl:"parameter,single_quotes" sql:"STORAGE_PROVIDER"`
	StorageAwsRoleArn    string                      `ddl:"parameter,single_quotes" sql:"STORAGE_AWS_ROLE_ARN"`
	StorageBaseUrl       string                      `ddl:"parameter,single_quotes" sql:"STORAGE_BASE_URL"`
	StorageAwsExternalId *string                     `ddl:"parameter,single_quotes" sql:"STORAGE_AWS_EXTERNAL_ID"`
	Encryption           *ExternalVolumeS3Encryption `ddl:"list,parentheses,no_comma" sql:"ENCRYPTION ="`
}
type ExternalVolumeS3Encryption struct {
	Type     S3EncryptionType `ddl:"parameter,single_quotes" sql:"TYPE"`
	KmsKeyId *string          `ddl:"parameter,single_quotes" sql:"KMS_KEY_ID"`
}
type GCSStorageLocationParams struct {
	Name               string                       `ddl:"parameter,single_quotes" sql:"NAME"`
	StorageProviderGcs string                       `ddl:"static" sql:"STORAGE_PROVIDER = 'GCS'"`
	StorageBaseUrl     string                       `ddl:"parameter,single_quotes" sql:"STORAGE_BASE_URL"`
	Encryption         *ExternalVolumeGCSEncryption `ddl:"list,parentheses,no_comma" sql:"ENCRYPTION ="`
}
type ExternalVolumeGCSEncryption struct {
	Type     GCSEncryptionType `ddl:"parameter,single_quotes" sql:"TYPE"`
	KmsKeyId *string           `ddl:"parameter,single_quotes" sql:"KMS_KEY_ID"`
}
type AzureStorageLocationParams struct {
	Name                 string `ddl:"parameter,single_quotes" sql:"NAME"`
	StorageProviderAzure string `ddl:"static" sql:"STORAGE_PROVIDER = 'AZURE'"`
	AzureTenantId        string `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
	StorageBaseUrl       string `ddl:"parameter,single_quotes" sql:"STORAGE_BASE_URL"`
}

// AlterExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-external-volume.
type AlterExternalVolumeOptions struct {
	alter                 bool                           `ddl:"static" sql:"ALTER"`
	externalVolume        bool                           `ddl:"static" sql:"EXTERNAL VOLUME"`
	IfExists              *bool                          `ddl:"keyword" sql:"IF EXISTS"`
	name                  AccountObjectIdentifier        `ddl:"identifier"`
	RemoveStorageLocation *string                        `ddl:"parameter,single_quotes,no_equals" sql:"REMOVE STORAGE_LOCATION"`
	Set                   *AlterExternalVolumeSet        `ddl:"keyword" sql:"SET"`
	AddStorageLocation    *ExternalVolumeStorageLocation `ddl:"parameter" sql:"ADD STORAGE_LOCATION"`
}
type AlterExternalVolumeSet struct {
	AllowWrites *bool   `ddl:"parameter" sql:"ALLOW_WRITES"`
	Comment     *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// DropExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-external-volume.
type DropExternalVolumeOptions struct {
	drop           bool                    `ddl:"static" sql:"DROP"`
	externalVolume bool                    `ddl:"static" sql:"EXTERNAL VOLUME"`
	IfExists       *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name           AccountObjectIdentifier `ddl:"identifier"`
}

// DescribeExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-external-volume.
type DescribeExternalVolumeOptions struct {
	describe       bool                    `ddl:"static" sql:"DESCRIBE"`
	externalVolume bool                    `ddl:"static" sql:"EXTERNAL VOLUME"`
	name           AccountObjectIdentifier `ddl:"identifier"`
}
type externalVolumeDescRow struct {
	ParentProperty  string `db:"parent_property"`
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}
type ExternalVolumeProperty struct {
	Parent  string
	Name    string
	Type    string
	Value   string
	Default string
}

// ShowExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-external-volumes.
type ShowExternalVolumeOptions struct {
	show            bool  `ddl:"static" sql:"SHOW"`
	externalVolumes bool  `ddl:"static" sql:"EXTERNAL VOLUMES"`
	Like            *Like `ddl:"keyword" sql:"LIKE"`
}
type externalVolumeShowRow struct {
	Name        string         `db:"name"`
	AllowWrites bool           `db:"allow_writes"`
	Comment     sql.NullString `db:"comment"`
}
type ExternalVolume struct {
	Name        string
	AllowWrites bool
	Comment     string
}
