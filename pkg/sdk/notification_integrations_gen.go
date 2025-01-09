package sdk

import (
	"context"
	"database/sql"
	"time"
)

type NotificationIntegrations interface {
	Create(ctx context.Context, request *CreateNotificationIntegrationRequest) error
	Alter(ctx context.Context, request *AlterNotificationIntegrationRequest) error
	Drop(ctx context.Context, request *DropNotificationIntegrationRequest) error
	Show(ctx context.Context, request *ShowNotificationIntegrationRequest) ([]NotificationIntegration, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*NotificationIntegration, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]NotificationIntegrationProperty, error)
}

// CreateNotificationIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration.
type CreateNotificationIntegrationOptions struct {
	create                   bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	notificationIntegration  bool                      `ddl:"static" sql:"NOTIFICATION INTEGRATION"`
	IfNotExists              *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                     AccountObjectIdentifier   `ddl:"identifier"`
	Enabled                  bool                      `ddl:"parameter" sql:"ENABLED"`
	AutomatedDataLoadsParams *AutomatedDataLoadsParams `ddl:"keyword"`
	PushNotificationParams   *PushNotificationParams   `ddl:"keyword"`
	EmailParams              *EmailParams              `ddl:"keyword"`
	Comment                  *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type NotificationIntegrationAllowedRecipient struct {
	Email string `ddl:"keyword,single_quotes"`
}

type AutomatedDataLoadsParams struct {
	notificationType string            `ddl:"static" sql:"TYPE = QUEUE"`
	GoogleAutoParams *GoogleAutoParams `ddl:"keyword"`
	AzureAutoParams  *AzureAutoParams  `ddl:"keyword"`
}

type GoogleAutoParams struct {
	notificationProvider      string `ddl:"static" sql:"NOTIFICATION_PROVIDER = GCP_PUBSUB"`
	GcpPubsubSubscriptionName string `ddl:"parameter,single_quotes" sql:"GCP_PUBSUB_SUBSCRIPTION_NAME"`
}

type AzureAutoParams struct {
	notificationProvider        string `ddl:"static" sql:"NOTIFICATION_PROVIDER = AZURE_STORAGE_QUEUE"`
	AzureStorageQueuePrimaryUri string `ddl:"parameter,single_quotes" sql:"AZURE_STORAGE_QUEUE_PRIMARY_URI"`
	AzureTenantId               string `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
}

type PushNotificationParams struct {
	direction        string            `ddl:"static" sql:"DIRECTION = OUTBOUND"`
	notificationType string            `ddl:"static" sql:"TYPE = QUEUE"`
	AmazonPushParams *AmazonPushParams `ddl:"keyword"`
	GooglePushParams *GooglePushParams `ddl:"keyword"`
	AzurePushParams  *AzurePushParams  `ddl:"keyword"`
}

type AmazonPushParams struct {
	notificationProvider string `ddl:"static" sql:"NOTIFICATION_PROVIDER = AWS_SNS"`
	AwsSnsTopicArn       string `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC_ARN"`
	AwsSnsRoleArn        string `ddl:"parameter,single_quotes" sql:"AWS_SNS_ROLE_ARN"`
}

type GooglePushParams struct {
	notificationProvider string `ddl:"static" sql:"NOTIFICATION_PROVIDER = GCP_PUBSUB"`
	GcpPubsubTopicName   string `ddl:"parameter,single_quotes" sql:"GCP_PUBSUB_TOPIC_NAME"`
}

type AzurePushParams struct {
	notificationProvider        string `ddl:"static" sql:"NOTIFICATION_PROVIDER = AZURE_EVENT_GRID"`
	AzureEventGridTopicEndpoint string `ddl:"parameter,single_quotes" sql:"AZURE_EVENT_GRID_TOPIC_ENDPOINT"`
	AzureTenantId               string `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
}

type EmailParams struct {
	notificationType  string                                    `ddl:"static" sql:"TYPE = EMAIL"`
	AllowedRecipients []NotificationIntegrationAllowedRecipient `ddl:"parameter,parentheses" sql:"ALLOWED_RECIPIENTS"`
}

// AlterNotificationIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration.
type AlterNotificationIntegrationOptions struct {
	alter                   bool                                     `ddl:"static" sql:"ALTER"`
	notificationIntegration bool                                     `ddl:"static" sql:"NOTIFICATION INTEGRATION"`
	IfExists                *bool                                    `ddl:"keyword" sql:"IF EXISTS"`
	name                    AccountObjectIdentifier                  `ddl:"identifier"`
	Set                     *NotificationIntegrationSet              `ddl:"keyword" sql:"SET"`
	UnsetEmailParams        *NotificationIntegrationUnsetEmailParams `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTags                 []TagAssociation                         `ddl:"keyword" sql:"SET TAG"`
	UnsetTags               []ObjectIdentifier                       `ddl:"keyword" sql:"UNSET TAG"`
}

type NotificationIntegrationSet struct {
	Enabled        *bool           `ddl:"parameter" sql:"ENABLED"`
	SetPushParams  *SetPushParams  `ddl:"keyword"`
	SetEmailParams *SetEmailParams `ddl:"keyword"`
	Comment        *string         `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type SetPushParams struct {
	SetAmazonPush *SetAmazonPush `ddl:"keyword"`
	SetGooglePush *SetGooglePush `ddl:"keyword"`
	SetAzurePush  *SetAzurePush  `ddl:"keyword"`
}

type SetAmazonPush struct {
	AwsSnsTopicArn string `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC_ARN"`
	AwsSnsRoleArn  string `ddl:"parameter,single_quotes" sql:"AWS_SNS_ROLE_ARN"`
}

type SetGooglePush struct {
	GcpPubsubSubscriptionName string `ddl:"parameter,single_quotes" sql:"GCP_PUBSUB_SUBSCRIPTION_NAME"`
}

type SetAzurePush struct {
	AzureStorageQueuePrimaryUri string `ddl:"parameter,single_quotes" sql:"AZURE_STORAGE_QUEUE_PRIMARY_URI"`
	AzureTenantId               string `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
}

type SetEmailParams struct {
	AllowedRecipients []NotificationIntegrationAllowedRecipient `ddl:"parameter,parentheses" sql:"ALLOWED_RECIPIENTS"`
}

type NotificationIntegrationUnsetEmailParams struct {
	AllowedRecipients *bool `ddl:"keyword" sql:"ALLOWED_RECIPIENTS"`
	Comment           *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropNotificationIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-integration.
type DropNotificationIntegrationOptions struct {
	drop                    bool                    `ddl:"static" sql:"DROP"`
	notificationIntegration bool                    `ddl:"static" sql:"NOTIFICATION INTEGRATION"`
	IfExists                *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                    AccountObjectIdentifier `ddl:"identifier"`
}

// ShowNotificationIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-integrations.
type ShowNotificationIntegrationOptions struct {
	show                     bool  `ddl:"static" sql:"SHOW"`
	notificationIntegrations bool  `ddl:"static" sql:"NOTIFICATION INTEGRATIONS"`
	Like                     *Like `ddl:"keyword" sql:"LIKE"`
}

type showNotificationIntegrationsDbRow struct {
	Name      string         `db:"name"`
	Type      string         `db:"type"`
	Category  string         `db:"category"`
	Enabled   bool           `db:"enabled"`
	Comment   sql.NullString `db:"comment"`
	CreatedOn time.Time      `db:"created_on"`
}

type NotificationIntegration struct {
	Name             string
	NotificationType string
	Category         string
	Enabled          bool
	Comment          string
	CreatedOn        time.Time
}

func (v *NotificationIntegration) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

// DescribeNotificationIntegrationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-integration.
type DescribeNotificationIntegrationOptions struct {
	describe                bool                    `ddl:"static" sql:"DESCRIBE"`
	notificationIntegration bool                    `ddl:"static" sql:"NOTIFICATION INTEGRATION"`
	name                    AccountObjectIdentifier `ddl:"identifier"`
}

type descNotificationIntegrationsDbRow struct {
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

type NotificationIntegrationProperty struct {
	Name    string
	Type    string
	Value   string
	Default string
}
