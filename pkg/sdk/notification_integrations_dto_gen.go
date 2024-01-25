package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateNotificationIntegrationOptions]   = new(CreateNotificationIntegrationRequest)
	_ optionsProvider[AlterNotificationIntegrationOptions]    = new(AlterNotificationIntegrationRequest)
	_ optionsProvider[DropNotificationIntegrationOptions]     = new(DropNotificationIntegrationRequest)
	_ optionsProvider[ShowNotificationIntegrationOptions]     = new(ShowNotificationIntegrationRequest)
	_ optionsProvider[DescribeNotificationIntegrationOptions] = new(DescribeNotificationIntegrationRequest)
)

type CreateNotificationIntegrationRequest struct {
	OrReplace                *bool
	IfNotExists              *bool
	name                     AccountObjectIdentifier // required
	Enabled                  bool                    // required
	AutomatedDataLoadsParams *AutomatedDataLoadsParamsRequest
	PushNotificationParams   *PushNotificationParamsRequest
	EmailParams              *EmailParamsRequest
	Comment                  *string
}

func (r *CreateNotificationIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

type AutomatedDataLoadsParamsRequest struct {
	GoogleAutoParams *GoogleAutoParamsRequest
	AzureAutoParams  *AzureAutoParamsRequest
}

type GoogleAutoParamsRequest struct {
	GcpPubsubSubscriptionName string // required
}

type AzureAutoParamsRequest struct {
	AzureStorageQueuePrimaryUri string // required
	AzureTenantId               string // required
}

type PushNotificationParamsRequest struct {
	AmazonPushParams *AmazonPushParamsRequest
	GooglePushParams *GooglePushParamsRequest
	AzurePushParams  *AzurePushParamsRequest
}

type AmazonPushParamsRequest struct {
	AwsSnsTopicArn string // required
	AwsSnsRoleArn  string // required
}

type GooglePushParamsRequest struct {
	GcpPubsubTopicName string // required
}

type AzurePushParamsRequest struct {
	AzureEventGridTopicEndpoint string // required
	AzureTenantId               string // required
}

type EmailParamsRequest struct {
	AllowedRecipients []NotificationIntegrationAllowedRecipient
}

type AlterNotificationIntegrationRequest struct {
	IfExists         *bool
	name             AccountObjectIdentifier // required
	Set              *NotificationIntegrationSetRequest
	UnsetEmailParams *NotificationIntegrationUnsetEmailParamsRequest
	SetTags          []TagAssociation
	UnsetTags        []ObjectIdentifier
}

type NotificationIntegrationSetRequest struct {
	Enabled        *bool
	SetPushParams  *SetPushParamsRequest
	SetEmailParams *SetEmailParamsRequest
	Comment        *string
}

type SetPushParamsRequest struct {
	SetAmazonPush *SetAmazonPushRequest
	SetGooglePush *SetGooglePushRequest
	SetAzurePush  *SetAzurePushRequest
}

type SetAmazonPushRequest struct {
	AwsSnsTopicArn string // required
	AwsSnsRoleArn  string // required
}

type SetGooglePushRequest struct {
	GcpPubsubSubscriptionName string // required
}

type SetAzurePushRequest struct {
	AzureStorageQueuePrimaryUri string // required
	AzureTenantId               string // required
}

type SetEmailParamsRequest struct {
	AllowedRecipients []NotificationIntegrationAllowedRecipient // required
}

type NotificationIntegrationUnsetEmailParamsRequest struct {
	AllowedRecipients *bool
	Comment           *bool
}

type DropNotificationIntegrationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowNotificationIntegrationRequest struct {
	Like *Like
}

type DescribeNotificationIntegrationRequest struct {
	name AccountObjectIdentifier // required
}
