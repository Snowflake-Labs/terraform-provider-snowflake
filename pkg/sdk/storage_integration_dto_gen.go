package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateStorageIntegrationOptions] = new(CreateStorageIntegrationRequest)
	_ optionsProvider[AlterStorageIntegrationOptions]  = new(AlterStorageIntegrationRequest)
)

type CreateStorageIntegrationRequest struct {
	OrReplace                  *bool
	IfNotExists                *bool
	name                       AccountObjectIdentifier // required
	S3StorageProviderParams    *S3StorageParamsRequest
	GCSStorageProviderParams   *GCSStorageParamsRequest
	AzureStorageProviderParams *AzureStorageParamsRequest
	Enabled                    bool     // required
	StorageAllowedLocations    []string // required
	StorageBlockedLocations    []string
	Comment                    *string
}

type S3StorageParamsRequest struct {
	StorageAwsRoleArn   string // required
	StorageAwsObjectAcl *string
}

type GCSStorageParamsRequest struct {
}

type AzureStorageParamsRequest struct {
	AzureTenantId *string // required
}

type AlterStorageIntegrationRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	Set       *StorageIntegrationSetRequest
	Unset     *StorageIntegrationUnsetRequest
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
}

type StorageIntegrationSetRequest struct {
	SetS3Params             *SetS3StorageParamsRequest
	SetAzureParams          *SetAzureStorageParamsRequest
	Enabled                 bool
	StorageAllowedLocations []string
	StorageBlockedLocations []string
	Comment                 *string
}

type SetS3StorageParamsRequest struct {
	StorageAwsRoleArn   string // required
	StorageAwsObjectAcl *string
}

type SetAzureStorageParamsRequest struct {
	AzureTenantId *string // required
}

type StorageIntegrationUnsetRequest struct {
	Enabled                 *bool
	StorageBlockedLocations *bool
	Comment                 *bool
}
