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
				Config: notificationIntegrationConfig(accName, storageUri, tenant),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_storage_queue_primary_uri", storageUri),
					resource.TestCheckResourceAttr("snowflake_notification_integration.test", "azure_tenant_id", tenant),
				),
			},
		},
	})
}

func notificationIntegrationConfig(name string, azureStorageQueuePrimaryUri string, azureTenantId string) string {
	s := `
resource "snowflake_notification_integration" "test" {
  name                            = "%s"
  azure_storage_queue_primary_uri = "%s"
  azure_tenant_id                 = "%s"
}
`
	return fmt.Sprintf(s, name, azureStorageQueuePrimaryUri, azureTenantId)
}
