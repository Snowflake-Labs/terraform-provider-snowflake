package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_NotificationAzureIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_NOTIFICATION_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAcc_NotificationAzureIntegration")
	}
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	storageURI := "azure://great-bucket/great-path/"
	tenant := "some-guid"

	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: azureNotificationIntegrationConfig(accName, storageURI, tenant),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AZURE_STORAGE_QUEUE"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_storage_queue_primary_uri", storageURI),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_tenant_id", tenant),
				),
			},
		},
	})
}

func TestAcc_NotificationGCPIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_NOTIFICATION_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAcc_NotificationGCPIntegration")
	}
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	gcpNotificationDirection := "INBOUND"

	pubsubName := "projects/project-1234/subscriptions/sub2"
	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: gcpNotificationIntegrationConfig(accName, pubsubName, gcpNotificationDirection),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "GCP_PUBSUB"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "gcp_pubsub_subscription_name", pubsubName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "direction", gcpNotificationDirection),
				),
			},
		},
	})
}

func azureNotificationIntegrationConfig(name string, azureStorageQueuePrimaryURI string, azureTenantID string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  name                            = "%s"
  notification_provider			  = "%s"
  azure_storage_queue_primary_uri = "%s"
  azure_tenant_id                 = "%s"
}
`
	return fmt.Sprintf(s, name, "AZURE_STORAGE_QUEUE", azureStorageQueuePrimaryURI, azureTenantID)
}

func gcpNotificationIntegrationConfig(name string, gcpPubsubSubscriptionName string, gcpNotificationDirection string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  name                            = "%s"
  notification_provider           = "%s"
  gcp_pubsub_subscription_name    = "%s"
  direction                       = "%s"
}
`
	return fmt.Sprintf(s, name, "GCP_PUBSUB", gcpPubsubSubscriptionName, gcpNotificationDirection)
}
