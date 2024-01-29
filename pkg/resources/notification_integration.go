package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO [SNOW-TODO]: split this resource into smaller ones as part of SNOW-867235
// TODO [SNOW-TODO]: remove SQS entirely
// TODO [SNOW-TODO]: support Azure push notifications
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
	},
	"direction": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"INBOUND", "OUTBOUND"}, true),
		Description:  "Direction of the cloud messaging with respect to Snowflake (required only for error notifications)",
		ForceNew:     true,
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
	id := d.Id()

	stmt := snowflake.NewNotificationIntegrationBuilder(d.Id()).Show()
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call

	s, err := snowflake.ScanNotificationIntegration(row)
	if err != nil {
		return fmt.Errorf("could not show notification integration: %w", err)
	}

	// Note: category must be NOTIFICATION or something is broken
	if c := s.Category.String; c != "NOTIFICATION" {
		return fmt.Errorf("expected %v to be a NOTIFICATION integration, got %v", id, c)
	}

	if err := d.Set("name", s.Name.String); err != nil {
		return err
	}

	// Snowflake returns "QUEUE - AZURE_STORAGE_QUEUE" instead of simple "QUEUE" as a type
	// so it needs to be parsed in order to not show a diff in Terraform
	typeParts := strings.Split(s.Type.String, "-")
	parsedType := strings.TrimSpace(typeParts[0])
	if err := d.Set("type", parsedType); err != nil {
		return err
	}

	if err := d.Set("created_on", s.CreatedOn.String); err != nil {
		return err
	}

	if err := d.Set("enabled", s.Enabled.Bool); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	// We need to grab them in a loop
	var k, pType string
	var v, n interface{}
	stmt = snowflake.NewNotificationIntegrationBuilder(d.Id()).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("could not describe notification integration: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &n); err != nil {
			return err
		}
		switch k {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "DIRECTION":
			if err := d.Set("direction", v.(string)); err != nil {
				return err
			}
		case "NOTIFICATION_PROVIDER":
			if err := d.Set("notification_provider", v.(string)); err != nil {
				return err
			}
		case "AZURE_STORAGE_QUEUE_PRIMARY_URI":
			if err := d.Set("azure_storage_queue_primary_uri", v.(string)); err != nil {
				return err
			}
		case "AZURE_TENANT_ID":
			if err := d.Set("azure_tenant_id", v.(string)); err != nil {
				return err
			}
		case "AWS_SQS_ARN":
			if err := d.Set("aws_sqs_arn", v.(string)); err != nil {
				return err
			}
		case "AWS_SQS_ROLE_ARN":
			if err := d.Set("aws_sqs_role_arn", v.(string)); err != nil {
				return err
			}
		case "AWS_SQS_EXTERNAL_ID":
			if err := d.Set("aws_sqs_external_id", v.(string)); err != nil {
				return err
			}
		case "AWS_SQS_IAM_USER_ARN":
			if err := d.Set("aws_sqs_iam_user_arn", v.(string)); err != nil {
				return err
			}
		case "AWS_SNS_TOPIC_ARN":
			if err := d.Set("aws_sns_topic_arn", v.(string)); err != nil {
				return err
			}
		case "AWS_SNS_ROLE_ARN":
			if err := d.Set("aws_sns_role_arn", v.(string)); err != nil {
				return err
			}
		case "SF_AWS_EXTERNAL_ID":
			if err := d.Set("aws_sns_external_id", v.(string)); err != nil {
				return err
			}
		case "SF_AWS_IAM_USER_ARN":
			if err := d.Set("aws_sns_iam_user_arn", v.(string)); err != nil {
				return err
			}
		case "GCP_PUBSUB_SUBSCRIPTION_NAME":
			if err := d.Set("gcp_pubsub_subscription_name", v.(string)); err != nil {
				return err
			}
		case "GCP_PUBSUB_TOPIC_NAME":
			if err := d.Set("gcp_pubsub_topic_name", v.(string)); err != nil {
				return err
			}
		case "GCP_PUBSUB_SERVICE_ACCOUNT":
			if err := d.Set("gcp_pubsub_service_account", v.(string)); err != nil {
				return err
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateNotificationIntegration implements schema.UpdateFunc.
func UpdateNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.NewNotificationIntegrationBuilder(id).Alter()

	// This is required in case the only change is to UNSET STORAGE_ALLOWED_LOCATIONS.
	// Not sure if there is a more elegant way of determining this
	var runSetStatement bool

	if d.HasChange("comment") {
		runSetStatement = true
		stmt.SetString("COMMENT", d.Get("comment").(string))
	}

	if d.HasChange("type") {
		runSetStatement = true
		stmt.SetString("TYPE", d.Get("type").(string))
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}

	if d.HasChange("direction") {
		runSetStatement = true
		stmt.SetString("DIRECTION", d.Get("direction").(string))
	}

	if d.HasChange("notification_provider") {
		runSetStatement = true
		stmt.SetString("NOTIFICATION_PROVIDER", d.Get("notification_provider").(string))
	}

	if d.HasChange("azure_storage_queue_primary_uri") {
		runSetStatement = true
		stmt.SetString("AZURE_STORAGE_QUEUE_PRIMARY_URI", d.Get("azure_storage_queue_primary_uri").(string))
	}

	if d.HasChange("azure_tenant_id") {
		runSetStatement = true
		stmt.SetString("AZURE_TENANT_ID", d.Get("azure_tenant_id").(string))
	}

	if d.HasChange("aws_sqs_arn") {
		runSetStatement = true
		stmt.SetString("AWS_SQS_ARN", d.Get("aws_sqs_arn").(string))
	}

	if d.HasChange("aws_sqs_role_arn") {
		runSetStatement = true
		stmt.SetString("AWS_SQS_ROLE_ARN", d.Get("aws_sqs_role_arn").(string))
	}

	if d.HasChange("aws_sns_topic_arn") {
		runSetStatement = true
		stmt.SetString("AWS_SNS_TOPIC_ARN", d.Get("aws_sns_topic_arn").(string))
	}

	if d.HasChange("aws_sns_role_arn") {
		runSetStatement = true
		stmt.SetString("AWS_SNS_ROLE_ARN", d.Get("aws_sns_role_arn").(string))
	}

	if d.HasChange("gcp_pubsub_subscription_name") {
		runSetStatement = true
		stmt.SetString("GCP_PUBSUB_SUBSCRIPTION_NAME", d.Get("gcp_pubsub_subscription_name").(string))
	}

	if d.HasChange("gcp_pubsub_topic_name") {
		runSetStatement = true
		stmt.SetString("GCP_PUBSUB_TOPIC_NAME", d.Get("gcp_pubsub_topic_name").(string))
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
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
