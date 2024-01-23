package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateApiIntegrationOptions]   = new(CreateApiIntegrationRequest)
	_ optionsProvider[AlterApiIntegrationOptions]    = new(AlterApiIntegrationRequest)
	_ optionsProvider[DropApiIntegrationOptions]     = new(DropApiIntegrationRequest)
	_ optionsProvider[ShowApiIntegrationOptions]     = new(ShowApiIntegrationRequest)
	_ optionsProvider[DescribeApiIntegrationOptions] = new(DescribeApiIntegrationRequest)
)

type CreateApiIntegrationRequest struct {
	OrReplace              *bool
	IfNotExists            *bool
	name                   AccountObjectIdentifier // required
	S3ApiProviderParams    *S3ApiParamsRequest
	AzureApiProviderParams *AzureApiParamsRequest
	GCSApiProviderParams   *GCSApiParamsRequest
	ApiAllowedPrefixes     []ApiIntegrationEndpointPrefix // required
	ApiBlockedPrefixes     []ApiIntegrationEndpointPrefix
	Enabled                bool // required
	Comment                *string
}

type S3ApiParamsRequest struct {
	ApiProvider   ApiIntegrationAwsApiProviderType // required
	ApiAwsRoleArn string                           // required
	ApiKey        *string
}

type AzureApiParamsRequest struct {
	AzureTenantId        string // required
	AzureAdApplicationId string // required
	ApiKey               *string
}

type GCSApiParamsRequest struct {
	GoogleAudience string // required
}

type AlterApiIntegrationRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	Set       *ApiIntegrationSetRequest
	Unset     *ApiIntegrationUnsetRequest
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
}

type ApiIntegrationSetRequest struct {
	S3Params           *SetS3ApiParamsRequest
	AzureParams        *SetAzureApiParamsRequest
	Enabled            *bool
	ApiAllowedPrefixes []ApiIntegrationEndpointPrefix
	ApiBlockedPrefixes []ApiIntegrationEndpointPrefix
	Comment            *string
}

type SetS3ApiParamsRequest struct {
	ApiAwsRoleArn *string
	ApiKey        *string
}

type SetAzureApiParamsRequest struct {
	AzureAdApplicationId string // required
	ApiKey               *string
}

type ApiIntegrationUnsetRequest struct {
	ApiKey             *bool
	Enabled            *bool
	ApiBlockedPrefixes *bool
	Comment            *bool
}

type DropApiIntegrationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowApiIntegrationRequest struct {
	Like *Like
}

type DescribeApiIntegrationRequest struct {
	name AccountObjectIdentifier // required
}
