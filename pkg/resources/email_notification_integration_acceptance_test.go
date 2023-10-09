package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_EmailNotificationIntegration(t *testing.T) {
	emailIntegrationName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	if _, ok := os.LookupEnv("SKIP_EMAIL_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAcc_EmailNotificationIntegration")
	}

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: emailNotificationIntegrationConfig(emailIntegrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", emailIntegrationName),
				),
			},
			{
				ResourceName:      "snowflake_email_notification_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func emailNotificationIntegrationConfig(name string) string {
	s := `
resource "snowflake_email_notification_integration" "test" {
  name               = "%s"
  enabled            = true
  allowed_recipients = ["joe@domain.com"] # Verified Email Addresses is required
  comment            = "test"
  /* 
Error: error creating notification integration: 394209 (22023): 
Email recipients in the given list at indexes [1] are not allowed. 
Either these email addresses are not yet validated or do not belong to any user in the current account.
  */
}
`
	return fmt.Sprintf(s, name)
}
