package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO [SNOW-1021713]: split this resource into smaller ones as part of SNOW-867235
// TODO [SNOW-1021713]: remove SQS entirely
// TODO [SNOW-1021713]: support Azure push notifications (AZURE_EVENT_GRID)
var notificationIntegrationSchema = map[string]*schema.Schema{
	// The first part of the schema is shared between all integration vendors
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "QUEUE",
		ValidateFunc: validation.StringInSlice([]string{"QUEUE"}, true),
		Description:  "A type of integration",
		ForceNew:     true,
		Deprecated:   "Will be removed - it is added automatically on the SDK level.",
	},
	"direction": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"INBOUND", "OUTBOUND"}, true),
		Description:  "Direction of the cloud messaging with respect to Snowflake (required only for error notifications)",
		ForceNew:     true,
		Deprecated:   "Will be removed - it is added automatically on the SDK level.",
	},
	// This part of the schema is the cloudProviderParams in the Snowflake documentation and differs between vendors
	"notification_provider": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"AZURE_STORAGE_QUEUE", "AWS_SQS", "AWS_SNS", "GCP_PUBSUB"}, true),
		Description:  "The third-party cloud message queuing service (supported values: AZURE_STORAGE_QUEUE, AWS_SNS, GCP_PUBSUB; AWS_SQS is deprecated and will be removed in the future provider versions)",
		ForceNew:     true,
	},
	"azure_storage_queue_primary_uri": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The queue ID for the Azure Queue Storage queue created for Event Grid notifications. Required for AZURE_STORAGE_QUEUE provider",
	},
	"azure_tenant_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The ID of the Azure Active Directory tenant used for identity management. Required for AZURE_STORAGE_QUEUE provider",
	},
	"aws_sqs_external_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role",
		Deprecated:  "No longer supported notification method",
	},
	"aws_sqs_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
		Deprecated:  "No longer supported notification method",
	},
	"aws_sqs_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "AWS SQS queue ARN for notification integration to connect to",
		Deprecated:  "No longer supported notification method",
	},
	"aws_sqs_role_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "AWS IAM role ARN for notification integration to assume",
		Deprecated:  "No longer supported notification method",
	},
	"aws_sns_external_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role",
	},
	"aws_sns_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
	},
	"aws_sns_topic_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "AWS SNS Topic ARN for notification integration to connect to. Required for AWS_SNS provider.",
	},
	"aws_sns_role_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "AWS IAM role ARN for notification integration to assume. Required for AWS_SNS provider",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A comment for the integration",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the notification integration was created.",
	},
	"gcp_pubsub_subscription_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The subscription id that Snowflake will listen to when using the GCP_PUBSUB provider.",
	},
	"gcp_pubsub_topic_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The topic id that Snowflake will use to push notifications.",
		// There is no alter SQL for gcp_pubsub_topic_name, therefore it has to be recreated.
		ForceNew: true,
	},
	"gcp_pubsub_service_account": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The GCP service account identifier that Snowflake will use when assuming the GCP role",
	},
}

// NotificationIntegration returns a pointer to the resource representing a notification integration.
func NotificationIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateNotificationIntegration,
		Read:   ReadNotificationIntegration,
		Update: UpdateNotificationIntegration,
		Delete: DeleteNotificationIntegration,

		Schema: notificationIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateNotificationIntegration implements schema.CreateFunc.
func CreateNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	enabled := d.Get("enabled").(bool)

	createRequest := sdk.NewCreateNotificationIntegrationRequest(id, enabled)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	notificationProvider := strings.ToUpper(d.Get("notification_provider").(string))
	switch notificationProvider {
	case "AWS_SNS":
		topic, ok := d.GetOk("aws_sns_topic_arn")
		if !ok {
			return fmt.Errorf("if you use AWS_SNS provider you must specify an aws_sns_topic_arn")
		}
		role, ok := d.GetOk("aws_sns_role_arn")
		if !ok {
			return fmt.Errorf("if you use AWS_SNS provider you must specify an aws_sns_role_arn")
		}
		createRequest.WithPushNotificationParams(
			sdk.NewPushNotificationParamsRequest().WithAmazonPushParams(sdk.NewAmazonPushParamsRequest(topic.(string), role.(string))),
		)
	case "GCP_PUBSUB":
		if v, ok := d.GetOk("gcp_pubsub_subscription_name"); ok {
			createRequest.WithAutomatedDataLoadsParams(
				sdk.NewAutomatedDataLoadsParamsRequest().WithGoogleAutoParams(sdk.NewGoogleAutoParamsRequest(v.(string))),
			)
		}
		if v, ok := d.GetOk("gcp_pubsub_topic_name"); ok {
			createRequest.WithPushNotificationParams(
				sdk.NewPushNotificationParamsRequest().WithGooglePushParams(sdk.NewGooglePushParamsRequest(v.(string))),
			)
		}
	case "AZURE_STORAGE_QUEUE":
		uri, ok := d.GetOk("azure_storage_queue_primary_uri")
		if !ok {
			return fmt.Errorf("if you use AZURE_STORAGE_QUEUE provider you must specify an azure_storage_queue_primary_uri")
		}
		tenantId, ok := d.GetOk("azure_tenant_id")
		if !ok {
			return fmt.Errorf("if you use AZURE_STORAGE_QUEUE provider you must specify an azure_tenant_id")
		}
		createRequest.WithAutomatedDataLoadsParams(
			sdk.NewAutomatedDataLoadsParamsRequest().WithAzureAutoParams(sdk.NewAzureAutoParamsRequest(uri.(string), tenantId.(string))),
		)
	default:
		return fmt.Errorf("unexpected provider %v", notificationProvider)
	}

	err := client.NotificationIntegrations.Create(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("error creating notification integration: %w", err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadNotificationIntegration(d, meta)
}

// ReadNotificationIntegration implements schema.ReadFunc.
func ReadNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.NotificationIntegrations.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] notification integration (%s) not found", d.Id())
		d.SetId("")
		return err
	}

	// Note: category must be NOTIFICATION or something is broken
	if c := integration.Category; c != "NOTIFICATION" {
		return fmt.Errorf("expected %v to be a NOTIFICATION integration, got %v", id, c)
	}

	if err := d.Set("name", integration.Name); err != nil {
		return err
	}

	if err := d.Set("comment", integration.Comment); err != nil {
		return err
	}

	if err := d.Set("created_on", integration.CreatedOn.String()); err != nil {
		return err
	}

	if err := d.Set("enabled", integration.Enabled); err != nil {
		return err
	}

	// Snowflake returns "QUEUE - AZURE_STORAGE_QUEUE" instead of simple "QUEUE" as a type
	// It needs to be parsed in order to not show a diff in Terraform.
	typeParts := strings.Split(integration.NotificationType, "-")
	parsedType := strings.TrimSpace(typeParts[0])
	if err := d.Set("type", parsedType); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	integrationProperties, err := client.NotificationIntegrations.Describe(ctx, id)
	if err != nil {
		return fmt.Errorf("could not describe notification integration: %w", err)
	}
	for _, property := range integrationProperties {
		name := property.Name
		value := property.Value
		switch name {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "DIRECTION":
			if err := d.Set("direction", value); err != nil {
				return err
			}
		case "NOTIFICATION_PROVIDER":
			if err := d.Set("notification_provider", value); err != nil {
				return err
			}
		case "AZURE_STORAGE_QUEUE_PRIMARY_URI":
			if err := d.Set("azure_storage_queue_primary_uri", value); err != nil {
				return err
			}
		case "AZURE_TENANT_ID":
			if err := d.Set("azure_tenant_id", value); err != nil {
				return err
			}
		case "AWS_SNS_TOPIC_ARN":
			if err := d.Set("aws_sns_topic_arn", value); err != nil {
				return err
			}
		case "AWS_SNS_ROLE_ARN":
			if err := d.Set("aws_sns_role_arn", value); err != nil {
				return err
			}
		case "SF_AWS_EXTERNAL_ID":
			if err := d.Set("aws_sns_external_id", value); err != nil {
				return err
			}
		case "SF_AWS_IAM_USER_ARN":
			if err := d.Set("aws_sns_iam_user_arn", value); err != nil {
				return err
			}
		case "GCP_PUBSUB_SUBSCRIPTION_NAME":
			if err := d.Set("gcp_pubsub_subscription_name", value); err != nil {
				return err
			}
		case "GCP_PUBSUB_TOPIC_NAME":
			if err := d.Set("gcp_pubsub_topic_name", value); err != nil {
				return err
			}
		case "GCP_PUBSUB_SERVICE_ACCOUNT":
			if err := d.Set("gcp_pubsub_service_account", value); err != nil {
				return err
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", name)
		}
	}

	return err
}

// UpdateNotificationIntegration implements schema.UpdateFunc.
func UpdateNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	var runSetStatement bool
	setRequest := sdk.NewNotificationIntegrationSetRequest()
	if d.HasChange("comment") {
		runSetStatement = true
		setRequest.WithComment(sdk.String(d.Get("comment").(string)))
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		setRequest.WithEnabled(sdk.Bool(d.Get("enabled").(bool)))
	}

	notificationProvider := strings.ToUpper(d.Get("notification_provider").(string))
	switch notificationProvider {
	case "AWS_SNS":
		if d.HasChange("aws_sns_topic_arn") || d.HasChange("aws_sns_role_arn") {
			runSetStatement = true
			setAmazonPush := sdk.NewSetAmazonPushRequest(d.Get("aws_sns_topic_arn").(string), d.Get("aws_sns_role_arn").(string))
			setRequest.WithSetPushParams(sdk.NewSetPushParamsRequest().WithSetAmazonPush(setAmazonPush))
		}
	case "GCP_PUBSUB":
		if d.HasChange("gcp_pubsub_subscription_name") {
			runSetStatement = true
			setGooglePush := sdk.NewSetGooglePushRequest(d.Get("gcp_pubsub_subscription_name").(string))
			setRequest.WithSetPushParams(sdk.NewSetPushParamsRequest().WithSetGooglePush(setGooglePush))
		}
	case "AZURE_STORAGE_QUEUE":
		if d.HasChange("azure_storage_queue_primary_uri") || d.HasChange("azure_tenant_id") {
			runSetStatement = true
			setAzurePush := sdk.NewSetAzurePushRequest(d.Get("azure_storage_queue_primary_uri").(string), d.Get("azure_tenant_id").(string))
			setRequest.WithSetPushParams(sdk.NewSetPushParamsRequest().WithSetAzurePush(setAzurePush))
		}
	default:
		return fmt.Errorf("unexpected provider %v", notificationProvider)
	}

	if runSetStatement {
		err := client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithSet(setRequest))
		if err != nil {
			return fmt.Errorf("error updating notification integration: %w", err)
		}
	}

	return ReadNotificationIntegration(d, meta)
}

// DeleteNotificationIntegration implements schema.DeleteFunc.
func DeleteNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.NotificationIntegrations.Drop(ctx, sdk.NewDropNotificationIntegrationRequest(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
