// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowNotificationIntegrationSchema represents output of SHOW query for the single NotificationIntegration.
var ShowNotificationIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"notification_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"category": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowNotificationIntegrationSchema

func NotificationIntegrationToSchema(notificationIntegration *sdk.NotificationIntegration) map[string]any {
	notificationIntegrationSchema := make(map[string]any)
	notificationIntegrationSchema["name"] = notificationIntegration.Name
	notificationIntegrationSchema["notification_type"] = notificationIntegration.NotificationType
	notificationIntegrationSchema["category"] = notificationIntegration.Category
	notificationIntegrationSchema["enabled"] = notificationIntegration.Enabled
	notificationIntegrationSchema["comment"] = notificationIntegration.Comment
	notificationIntegrationSchema["created_on"] = notificationIntegration.CreatedOn.String()
	return notificationIntegrationSchema
}

var _ = NotificationIntegrationToSchema
