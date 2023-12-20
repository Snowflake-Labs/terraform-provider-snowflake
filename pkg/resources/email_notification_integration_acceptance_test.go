package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_EmailNotificationIntegration(t *testing.T) {
	env := os.Getenv("SKIP_EMAIL_INTEGRATION_TESTS")
	if env != "" {
		t.Skip("Skipping TestAcc_EmailNotificationIntegration")
	}

	emailIntegrationName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	verifiedEmail := "artur.sawicki@snowflake.com"

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: emailNotificationIntegrationConfig(emailIntegrationName, verifiedEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", emailIntegrationName),
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients.0", verifiedEmail),
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

func emailNotificationIntegrationConfig(name string, email string) string {
	s := `
resource "snowflake_email_notification_integration" "test" {
  name               = "%s"
  enabled            = true
  allowed_recipients = ["%s"] # Verified Email Addresses is required
  comment            = "test"
  /*
Error: error creating notification integration: 394209 (22023):
Email recipients in the given list at indexes [1] are not allowed.
Either these email addresses are not yet validated or do not belong to any user in the current account.
  */
}
`
	return fmt.Sprintf(s, name, email)
}

// TestAcc_EmailNotificationIntegration_issue2223 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2223 issue.
// Snowflake allowed empty allowed recipients in https://docs.snowflake.com/en/release-notes/2023/7_40#email-notification-integrations-allowed-recipients-no-longer-required.
func TestAcc_EmailNotificationIntegration_issue2223(t *testing.T) {
	emailIntegrationName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: emailNotificationIntegrationWithoutRecipientsConfig(emailIntegrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_email_notification_integration.test", "name", emailIntegrationName),
					resource.TestCheckNoResourceAttr("snowflake_email_notification_integration.test", "allowed_recipients"),
				),
			},
		},
	})
}

func emailNotificationIntegrationWithoutRecipientsConfig(name string) string {
	s := `
resource "snowflake_email_notification_integration" "test" {
  name               = "%s"
  enabled            = true
}`
	return fmt.Sprintf(s, name)
}
