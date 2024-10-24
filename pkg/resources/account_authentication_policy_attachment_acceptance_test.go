package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountAuthenticationPolicyAttachment(t *testing.T) {
	policyName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountAuthenticationPolicyAttachmentConfig(acc.TestDatabaseName, acc.TestSchemaName, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_account_authentication_policy_attachment.att", "id"),
				),
			},
			{
				ResourceName:      "snowflake_account_authentication_policy_attachment.att",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func accountAuthenticationPolicyAttachmentConfig(databaseName, schemaName, policyName string) string {
	s := `
resource "snowflake_authentication_policy" "pa" {
	database   = "%s"
	schema     = "%s"
	name       = "%v"
}

resource "snowflake_account_authentication_policy_attachment" "att" {
	authentication_policy = snowflake_authentication_policy.pa.fully_qualified_name
}
`
	return fmt.Sprintf(s, databaseName, schemaName, policyName)
}
