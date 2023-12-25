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
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: userPasswordPolicyAttachmentConfig(acc.TestDatabaseName, acc.TestSchemaName, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_user_password_policy_attachment.att", "id"),
				),
			},
			// TODO: importer
		},
	})
}

func userPasswordPolicyAttachmentConfig(databaseName, schemaName, prefix string) string {
	s := `
resource "snowflake_user" "user_password_policy_attachment_acceptance_test_user" {
	name = "user_password_policy_attachment_acceptance_test_user"
}
resource "snowflake_password_policy" "pa" {
	database   = "%s"
	schema     = "%s"
	name       = "%v"
}

resource "snowflake_user_password_policy_attachment" "att" {
	password_policy = snowflake_password_policy.pa.qualified_name
	user_name = snowflake_user.user_password_policy_attachment_acceptance_test_user.name
}
`
	return fmt.Sprintf(s, databaseName, schemaName, prefix)
}
