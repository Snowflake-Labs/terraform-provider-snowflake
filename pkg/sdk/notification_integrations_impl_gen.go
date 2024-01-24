package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ NotificationIntegrations = (*notificationIntegrations)(nil)

type notificationIntegrations struct {
	client *Client
}

func (v *notificationIntegrations) Create(ctx context.Context, request *CreateNotificationIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *notificationIntegrations) Alter(ctx context.Context, request *AlterNotificationIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *notificationIntegrations) Drop(ctx context.Context, request *DropNotificationIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *notificationIntegrations) Show(ctx context.Context, request *ShowNotificationIntegrationRequest) ([]NotificationIntegration, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showNotificationIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showNotificationIntegrationsDbRow, NotificationIntegration](dbRows)
	return resultList, nil
}

func (v *notificationIntegrations) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*NotificationIntegration, error) {
	notificationIntegrations, err := v.Show(ctx, NewShowNotificationIntegrationRequest().WithLike(&Like{
		Pattern: String(id.Name()),
	}))
	if err != nil {
		return nil, err
	}
	return collections.FindOne(notificationIntegrations, func(r NotificationIntegration) bool { return r.Name == id.Name() })
}

func (v *notificationIntegrations) Describe(ctx context.Context, id AccountObjectIdentifier) ([]NotificationIntegrationProperty, error) {
	opts := &DescribeNotificationIntegrationOptions{
		name: id,
	}
	rows, err := validateAndQuery[descNotificationIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[descNotificationIntegrationsDbRow, NotificationIntegrationProperty](rows), nil
}

func (r *CreateNotificationIntegrationRequest) toOpts() *CreateNotificationIntegrationOptions {
	opts := &CreateNotificationIntegrationOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Enabled:     r.Enabled,

		Comment: r.Comment,
	}
	if r.AutomatedDataLoadsParams != nil {
		opts.AutomatedDataLoadsParams = &AutomatedDataLoadsParams{}
		if r.AutomatedDataLoadsParams.GoogleAutomatedDataLoad != nil {
			opts.AutomatedDataLoadsParams.GoogleAutomatedDataLoad = &GoogleAutomatedDataLoad{
				GcpPubsubSubscriptionName: r.AutomatedDataLoadsParams.GoogleAutomatedDataLoad.GcpPubsubSubscriptionName,
			}
		}
		if r.AutomatedDataLoadsParams.AzureAutomatedDataLoad != nil {
			opts.AutomatedDataLoadsParams.AzureAutomatedDataLoad = &AzureAutomatedDataLoad{
				AzureStorageQueuePrimaryUri: r.AutomatedDataLoadsParams.AzureAutomatedDataLoad.AzureStorageQueuePrimaryUri,
				AzureTenantId:               r.AutomatedDataLoadsParams.AzureAutomatedDataLoad.AzureTenantId,
			}
		}
	}
	if r.PushNotificationParams != nil {
		opts.PushNotificationParams = &PushNotificationParams{}
		if r.PushNotificationParams.AmazonPush != nil {
			opts.PushNotificationParams.AmazonPush = &AmazonPush{
				AwsSnsTopicArn: r.PushNotificationParams.AmazonPush.AwsSnsTopicArn,
				AwsSnsRoleArn:  r.PushNotificationParams.AmazonPush.AwsSnsRoleArn,
			}
		}
		if r.PushNotificationParams.GooglePush != nil {
			opts.PushNotificationParams.GooglePush = &GooglePush{
				GcpPubsubTopicName: r.PushNotificationParams.GooglePush.GcpPubsubTopicName,
			}
		}
		if r.PushNotificationParams.AzurePush != nil {
			opts.PushNotificationParams.AzurePush = &AzurePush{
				AzureEventGridTopicEndpoint: r.PushNotificationParams.AzurePush.AzureEventGridTopicEndpoint,
				AzureTenantId:               r.PushNotificationParams.AzurePush.AzureTenantId,
			}
		}
	}
	if r.EmailParams != nil {
		opts.EmailParams = &EmailParams{
			AllowedRecipients: r.EmailParams.AllowedRecipients,
		}
	}
	return opts
}

func (r *AlterNotificationIntegrationRequest) toOpts() *AlterNotificationIntegrationOptions {
	opts := &AlterNotificationIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &NotificationIntegrationSet{
			Enabled: r.Set.Enabled,

			Comment: r.Set.Comment,
		}
		if r.Set.SetPushParams != nil {
			opts.Set.SetPushParams = &SetPushParams{}
			if r.Set.SetPushParams.SetAmazonPush != nil {
				opts.Set.SetPushParams.SetAmazonPush = &SetAmazonPush{
					AwsSnsTopicArn: r.Set.SetPushParams.SetAmazonPush.AwsSnsTopicArn,
					AwsSnsRoleArn:  r.Set.SetPushParams.SetAmazonPush.AwsSnsRoleArn,
				}
			}
			if r.Set.SetPushParams.SetGooglePush != nil {
				opts.Set.SetPushParams.SetGooglePush = &SetGooglePush{
					GcpPubsubSubscriptionName: r.Set.SetPushParams.SetGooglePush.GcpPubsubSubscriptionName,
				}
			}
			if r.Set.SetPushParams.SetAzurePush != nil {
				opts.Set.SetPushParams.SetAzurePush = &SetAzurePush{
					AzureStorageQueuePrimaryUri: r.Set.SetPushParams.SetAzurePush.AzureStorageQueuePrimaryUri,
					AzureTenantId:               r.Set.SetPushParams.SetAzurePush.AzureTenantId,
				}
			}
		}
		if r.Set.SetEmailParams != nil {
			opts.Set.SetEmailParams = &SetEmailParams{
				AllowedRecipients: r.Set.SetEmailParams.AllowedRecipients,
			}
		}
	}
	if r.UnsetEmailParams != nil {
		opts.UnsetEmailParams = &NotificationIntegrationUnsetEmailParams{
			AllowedRecipients: r.UnsetEmailParams.AllowedRecipients,
			Comment:           r.UnsetEmailParams.Comment,
		}
	}
	return opts
}

func (r *DropNotificationIntegrationRequest) toOpts() *DropNotificationIntegrationOptions {
	opts := &DropNotificationIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowNotificationIntegrationRequest) toOpts() *ShowNotificationIntegrationOptions {
	opts := &ShowNotificationIntegrationOptions{
		Like: r.Like,
	}
	return opts
}

func (r showNotificationIntegrationsDbRow) convert() *NotificationIntegration {
	s := &NotificationIntegration{
		Name:             r.Name,
		NotificationType: r.Type,
		Category:         r.Category,
		Enabled:          r.Enabled,
		CreatedOn:        r.CreatedOn,
	}
	if r.Comment.Valid {
		s.Comment = r.Comment.String
	}
	return s
}

func (r *DescribeNotificationIntegrationRequest) toOpts() *DescribeNotificationIntegrationOptions {
	opts := &DescribeNotificationIntegrationOptions{
		name: r.name,
	}
	return opts
}

func (r descNotificationIntegrationsDbRow) convert() *NotificationIntegrationProperty {
	return &NotificationIntegrationProperty{
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}
}
