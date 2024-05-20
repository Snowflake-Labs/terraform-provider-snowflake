package ids

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

var (
	PrecreatedS3StorageIntegration    = sdk.NewAccountObjectIdentifier("S3_STORAGE_INTEGRATION")
	PrecreatedGcpStorageIntegration   = sdk.NewAccountObjectIdentifier("GCP_STORAGE_INTEGRATION")
	PrecreatedAzureStorageIntegration = sdk.NewAccountObjectIdentifier("AZURE_STORAGE_INTEGRATION")
)
