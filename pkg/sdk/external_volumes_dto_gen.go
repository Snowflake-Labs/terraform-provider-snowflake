package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateExternalVolumeOptions]   = new(CreateExternalVolumeRequest)
	_ optionsProvider[AlterExternalVolumeOptions]    = new(AlterExternalVolumeRequest)
	_ optionsProvider[DropExternalVolumeOptions]     = new(DropExternalVolumeRequest)
	_ optionsProvider[DescribeExternalVolumeOptions] = new(DescribeExternalVolumeRequest)
	_ optionsProvider[ShowExternalVolumeOptions]     = new(ShowExternalVolumeRequest)
)

type CreateExternalVolumeRequest struct {
	OrReplace        *bool
	IfNotExists      *bool
	name             AccountObjectIdentifier         // required
	StorageLocations []ExternalVolumeStorageLocation // required
	AllowWrites      *bool
	Comment          *string
}

type AlterExternalVolumeRequest struct {
	IfExists              *bool
	name                  AccountObjectIdentifier // required
	RemoveStorageLocation *string
	Set                   *AlterExternalVolumeSetRequest
	AddStorageLocation    *ExternalVolumeStorageLocationRequest
}

type AlterExternalVolumeSetRequest struct {
	AllowWrites *bool
	Comment     *string
}

type ExternalVolumeStorageLocationRequest struct {
	S3StorageLocationParams    *S3StorageLocationParamsRequest
	GCSStorageLocationParams   *GCSStorageLocationParamsRequest
	AzureStorageLocationParams *AzureStorageLocationParamsRequest
}

type S3StorageLocationParamsRequest struct {
	Name                 string            // required
	StorageProvider      S3StorageProvider // required
	StorageAwsRoleArn    string            // required
	StorageBaseUrl       string            // required
	StorageAwsExternalId *string
	Encryption           *ExternalVolumeS3EncryptionRequest
}

type ExternalVolumeS3EncryptionRequest struct {
	Type     S3EncryptionType // required
	KmsKeyId *string
}

type GCSStorageLocationParamsRequest struct {
	Name           string // required
	StorageBaseUrl string // required
	Encryption     *ExternalVolumeGCSEncryptionRequest
}

type ExternalVolumeGCSEncryptionRequest struct {
	Type     GCSEncryptionType // required
	KmsKeyId *string
}

type AzureStorageLocationParamsRequest struct {
	Name           string // required
	AzureTenantId  string // required
	StorageBaseUrl string // required
}

type DropExternalVolumeRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type DescribeExternalVolumeRequest struct {
	name AccountObjectIdentifier // required
}

type ShowExternalVolumeRequest struct {
	Like *Like
}
