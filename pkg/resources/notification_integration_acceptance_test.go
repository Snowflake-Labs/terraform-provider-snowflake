package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NotificationIntegration_AutoGoogle(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"
	const gcpOtherPubsubSubscriptionName = "projects/project-1234/subscriptions/other"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: googleAutoConfig(accName, gcpPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "GCP_PUBSUB"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "gcp_pubsub_subscription_name", gcpPubsubSubscriptionName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "gcp_pubsub_service_account"),
				),
			},
			// change parameters
			{
				Config: googleAutoConfig(accName, gcpOtherPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "GCP_PUBSUB"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "gcp_pubsub_subscription_name", gcpOtherPubsubSubscriptionName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "gcp_pubsub_service_account"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_NotificationIntegration_AutoAzure(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	const azureStorageQueuePrimaryUri = "azure://great-bucket/great-path/"
	const azureOtherStorageQueuePrimaryUri = "azure://great-bucket/other-great-path/"
	const azureTenantId = "00000000-0000-0000-0000-000000000000"
	const azureOtherTenantId = "11111111-1111-1111-1111-111111111111"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: azureAutoConfig(accName, azureStorageQueuePrimaryUri, azureTenantId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AZURE_STORAGE_QUEUE"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_storage_queue_primary_uri", azureStorageQueuePrimaryUri),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_tenant_id", azureTenantId),
					resource.TestCheckNoResourceAttr("snowflake_notification_integration.test", "direction"),
				),
			},
			// change parameters
			{
				Config: azureAutoConfig(accName, azureOtherStorageQueuePrimaryUri, azureOtherTenantId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AZURE_STORAGE_QUEUE"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_storage_queue_primary_uri", azureOtherStorageQueuePrimaryUri),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_tenant_id", azureOtherTenantId),
					resource.TestCheckNoResourceAttr("snowflake_notification_integration.test", "direction"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
				// it is not returned in DESCRIBE for azure automated data load
				ImportStateVerifyIgnore: []string{"azure_tenant_id"},
			},
		},
	})
}

func TestAcc_NotificationIntegration_PushAmazon(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	const awsSnsTopicArn = "arn:aws:sns:us-east-2:123456789012:MyTopic"
	const awsOtherSnsTopicArn = "arn:aws:sns:us-east-2:123456789012:OtherTopic"
	const awsSnsRoleArn = "arn:aws:iam::000000000001:/role/test"
	const awsOtherSnsRoleArn = "arn:aws:iam::000000000001:/role/other"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: amazonPushConfig(accName, awsSnsTopicArn, awsSnsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AWS_SNS"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_topic_arn", awsSnsTopicArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_role_arn", awsSnsRoleArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "OUTBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_external_id"),
				),
			},
			// change parameters
			{
				Config: amazonPushConfig(accName, awsOtherSnsTopicArn, awsOtherSnsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AWS_SNS"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_topic_arn", awsOtherSnsTopicArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_role_arn", awsOtherSnsRoleArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "OUTBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_external_id"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_NotificationIntegration_changeNotificationProvider(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"
	const awsSnsTopicArn = "arn:aws:sns:us-east-2:123456789012:MyTopic"
	const awsSnsRoleArn = "arn:aws:iam::000000000001:/role/test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				Config: googleAutoConfig(accName, gcpPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "GCP_PUBSUB"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "gcp_pubsub_subscription_name", gcpPubsubSubscriptionName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "gcp_pubsub_service_account"),
				),
			},
			// change provider to AWS
			{
				Config: amazonPushConfig(accName, awsSnsTopicArn, awsSnsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AWS_SNS"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_topic_arn", awsSnsTopicArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "aws_sns_role_arn", awsSnsRoleArn),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "OUTBOUND"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_notification_integration.test", "aws_sns_external_id"),
				),
			},
		},
	})
}

// TODO [SNOW-1017802]: handle after "create and describe notification integration - push google" test passes
func TestAcc_NotificationIntegration_PushGoogle(t *testing.T) {
	t.Skip("Skipping because can't be currently created. Check 'create and describe notification integration - push google' test in the SDK.")
}

// TODO [SNOW-1017802]: handle after "create and describe notification integration - push azure" test passes
// TODO [SNOW-1021713]: handle after it's added to the resource
func TestAcc_NotificationIntegration_PushAzure(t *testing.T) {
	t.Skip("Skipping because can't be currently created. Check 'create and describe notification integration - push azure' test in the SDK.")
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2501
func TestAcc_NotificationIntegration_migrateFromVersion085(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: googleAutoConfig(accName, gcpPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   googleAutoConfigWithoutDirection(accName, gcpPubsubSubscriptionName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", "INBOUND"),
				),
			},
		},
	})
}

func TestAcc_NotificationIntegration_migrateFromVersion085_explicitType(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NotificationIntegration),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: googleAutoConfigWithExplicitType(accName, gcpPubsubSubscriptionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   googleAutoConfig(accName, gcpPubsubSubscriptionName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
				),
			},
		},
	})
}

func googleAutoConfig(name string, gcpPubsubSubscriptionName string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  name                            = "%s"
  notification_provider           = "%s"
  gcp_pubsub_subscription_name    = "%s"
  direction                       = "INBOUND"
}
`
	return fmt.Sprintf(s, name, "GCP_PUBSUB", gcpPubsubSubscriptionName)
}

func googleAutoConfigWithoutDirection(name string, gcpPubsubSubscriptionName string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  name                            = "%s"
  notification_provider           = "%s"
  gcp_pubsub_subscription_name    = "%s"
}
`
	return fmt.Sprintf(s, name, "GCP_PUBSUB", gcpPubsubSubscriptionName)
}

func googleAutoConfigWithExplicitType(name string, gcpPubsubSubscriptionName string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  type                            = "QUEUE"
  name                            = "%s"
  notification_provider           = "%s"
  gcp_pubsub_subscription_name    = "%s"
  direction                       = "INBOUND"
}
`
	return fmt.Sprintf(s, name, "GCP_PUBSUB", gcpPubsubSubscriptionName)
}

func azureAutoConfig(name string, azureStorageQueuePrimaryUri string, azureTenantId string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  name                            = "%s"
  notification_provider			  = "%s"
  azure_storage_queue_primary_uri = "%s"
  azure_tenant_id                 = "%s"
}
`
	return fmt.Sprintf(s, name, "AZURE_STORAGE_QUEUE", azureStorageQueuePrimaryUri, azureTenantId)
}

func amazonPushConfig(name string, awsSnsTopicArn string, awsSnsRoleArn string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  name                            = "%s"
  notification_provider           = "%s"
  aws_sns_topic_arn               = "%s"
  aws_sns_role_arn                = "%s"
  direction                       = "OUTBOUND"
}
`
	return fmt.Sprintf(s, name, "AWS_SNS", awsSnsTopicArn, awsSnsRoleArn)
}
