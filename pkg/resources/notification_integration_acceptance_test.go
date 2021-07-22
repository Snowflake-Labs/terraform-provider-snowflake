package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_NotificationIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_NOTIFICATION_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAccNotificationIntegration")
	}
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	storageUri := "azure://great-bucket/great-path/"
	tenant := "some-guid"

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: azureNotificationIntegrationConfig(accName, storageUri, tenant),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "AZURE_STORAGE_QUEUE"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_storage_queue_primary_uri", storageUri),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_tenant_id", tenant),
				),
			},
		},
	})

	pubsubName := "projects/project-1234/subscriptions/sub2"
	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: gcpNotificationIntegrationConfig(accName, pubsubName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "notification_provider", "GCP_PUBSUB"),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "gcp_pubsub_subscription_name", pubsubName),
				),
			},
		},
	})
}

func azureNotificationIntegrationConfig(name string, azureStorageQueuePrimaryUri string, azureTenantId string) string {
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

func gcpNotificationIntegrationConfig(name string, gcpPubsubSubscriptionName string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  name                            = "%s"
  notification_provider           = "%s"
  gcp_pubsub_subscription_name    = "%s"
}
`
	return fmt.Sprintf(s, name, "GCP_PUBSUB", gcpPubsubSubscriptionName)
}
