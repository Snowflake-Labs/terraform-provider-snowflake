package sdk

import "context"

type StorageIntegrations interface {
	Create(ctx context.Context, request *CreateStorageIntegrationRequest) error
	Alter(ctx context.Context, request *AlterStorageIntegrationRequest) error
}

// CreateStorageIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration.
type CreateStorageIntegrationOptions struct {
	create                     bool                    `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	storageIntegration         bool                    `ddl:"static" sql:"STORAGE INTEGRATION"`
	IfNotExists                *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       AccountObjectIdentifier `ddl:"identifier"`
	externalStageType          string                  `ddl:"static" sql:"TYPE = EXTERNAL_STAGE"`
	S3StorageProviderParams    *S3StorageParams        `ddl:"keyword"`
	GCSStorageProviderParams   *GCSStorageParams       `ddl:"keyword"`
	AzureStorageProviderParams *AzureStorageParams     `ddl:"keyword"`
	Enabled                    bool                    `ddl:"parameter" sql:"ENABLED"`
	StorageAllowedLocations    []StorageLocation       `ddl:"parameter,parentheses" sql:"STORAGE_ALLOWED_LOCATIONS"`
	StorageBlockedLocations    []StorageLocation       `ddl:"parameter,parentheses" sql:"STORAGE_BLOCKED_LOCATIONS"`
	Comment                    *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type StorageLocation struct {
	Path string `ddl:"keyword,single_quotes"`
}

type S3StorageParams struct {
	storageProvider     string  `ddl:"static" sql:"STORAGE_PROVIDER = 'S3'"`
	StorageAwsRoleArn   string  `ddl:"parameter,single_quotes" sql:"STORAGE_AWS_ROLE_ARN"`
	StorageAwsObjectAcl *string `ddl:"parameter,single_quotes" sql:"STORAGE_AWS_OBJECT_ACL"`
}

type GCSStorageParams struct {
	storageProvider string `ddl:"static" sql:"STORAGE_PROVIDER = 'GCS'"`
}

type AzureStorageParams struct {
	storageProvider string  `ddl:"static" sql:"STORAGE_PROVIDER = 'AZURE'"`
	AzureTenantId   *string `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
}

// AlterStorageIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-storage-integration.
type AlterStorageIntegrationOptions struct {
	alter              bool                     `ddl:"static" sql:"ALTER"`
	storageIntegration bool                     `ddl:"static" sql:"STORAGE INTEGRATION"`
	IfExists           *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name               AccountObjectIdentifier  `ddl:"identifier"`
	Set                *StorageIntegrationSet   `ddl:"keyword" sql:"SET"`
	Unset              *StorageIntegrationUnset `ddl:"list" sql:"UNSET"`
	SetTags            []TagAssociation         `ddl:"keyword" sql:"SET TAG"`
	UnsetTags          []ObjectIdentifier       `ddl:"keyword" sql:"UNSET TAG"`
}

type StorageIntegrationSet struct {
	SetS3Params             *SetS3StorageParams    `ddl:"keyword"`
	SetAzureParams          *SetAzureStorageParams `ddl:"keyword"`
	Enabled                 bool                   `ddl:"parameter" sql:"ENABLED"`
	StorageAllowedLocations []StorageLocation      `ddl:"parameter,parentheses" sql:"STORAGE_ALLOWED_LOCATIONS"`
	StorageBlockedLocations []StorageLocation      `ddl:"parameter,parentheses" sql:"STORAGE_BLOCKED_LOCATIONS"`
	Comment                 *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type SetS3StorageParams struct {
	StorageAwsRoleArn   string  `ddl:"parameter,single_quotes" sql:"STORAGE_AWS_ROLE_ARN"`
	StorageAwsObjectAcl *string `ddl:"parameter,single_quotes" sql:"STORAGE_AWS_OBJECT_ACL"`
}

type SetAzureStorageParams struct {
	AzureTenantId *string `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
}

type StorageIntegrationUnset struct {
	Enabled                 *bool `ddl:"keyword" sql:"ENABLED"`
	StorageBlockedLocations *bool `ddl:"keyword" sql:"STORAGE_BLOCKED_LOCATIONS"`
	Comment                 *bool `ddl:"keyword" sql:"COMMENT"`
}
