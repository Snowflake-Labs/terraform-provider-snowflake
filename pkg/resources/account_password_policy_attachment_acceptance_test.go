package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_AccountPasswordPolicyAttachment(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountPasswordPolicyAttachmentConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_account_password_policy_attachment.att", "id"),
				),
			},
			{
				ResourceName:      "snowflake_account_password_policy_attachment.att",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"initially_suspended",
					"wait_for_provisioning",
					"query_acceleration_max_scale_factor",
					"max_concurrency_level",
					"statement_queued_timeout_in_seconds",
					"statement_timeout_in_seconds",
				},
			},
		},
	})
}

func accountPasswordPolicyAttachmentConfig(prefix string) string {
	s := `
resource "snowflake_database" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
	}

resource "snowflake_password_policy" "pa" {
	database   = snowflake_database.test.name
	schema     = snowflake_schema.test.name
	name       = "%v"
}

resource "snowflake_account_password_policy_attachment" "att" {
	password_policy = snowflake_password_policy.pa.qualified_name
}
`
	return fmt.Sprintf(s, prefix, prefix, prefix)
}
