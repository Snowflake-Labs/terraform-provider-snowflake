package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserPasswordPolicyAttachment(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: acc.TestAccProviders(),
		// TODO: what is this for?
		// PreCheck:     func() { acc.TestAccPreCheck(t) },
		// TODO: this probably needs to be setup
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: userPasswordPolicyAttachmentConfig(acc.TestDatabaseName, acc.TestSchemaName, prefix),
				// TODO: this only checks if the id is set, but not which value
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_user_password_policy_attachment.att", "id"),
				),
			},
			// TODO: I need to implement the importer
			// {
			// 	ResourceName:      "snowflake_user_password_policy_attachment.att",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateVerifyIgnore: []string{
			// 		"initially_suspended",
			// 		"wait_for_provisioning",
			// 		"query_acceleration_max_scale_factor",
			// 		"max_concurrency_level",
			// 		"statement_queued_timeout_in_seconds",
			// 		"statement_timeout_in_seconds",
			// 	},
			// },
		},
	})
}

// TODO: change the username
func userPasswordPolicyAttachmentConfig(databaseName, schemaName, prefix string) string {
	s := `
resource "snowflake_password_policy" "pa" {
	database   = "%s"
	schema     = "%s"
	name       = "%v"
}

resource "snowflake_user_password_policy_attachment" "att" {
	password_policy = snowflake_password_policy.pa.qualified_name
	user_name = "RBONET"
}
`
	return fmt.Sprintf(s, databaseName, schemaName, prefix)
}
