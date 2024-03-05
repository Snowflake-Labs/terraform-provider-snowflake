package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicyAttachment(t *testing.T) {
	user1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	policyName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: networkPolicyAttachmentConfigSingle(user1, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyName),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "false"),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "users.#", "1"),
				),
			},
			{
				Config: networkPolicyAttachmentConfig(user1, user2, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyName),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "false"),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "users.#", "2"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_network_policy_attachment.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func networkPolicyAttachmentConfigSingle(user1, policyName string) string {
	return fmt.Sprintf(`
resource "snowflake_user" "test-user1" {
	name = "%s"
}

resource "snowflake_network_policy" "test" {
	name            = "%v"
	allowed_ip_list = ["192.168.0.100/24", "29.254.123.20"]
}

resource "snowflake_network_policy_attachment" "test" {
	network_policy_name = snowflake_network_policy.test.name
	set_for_account     = false
	users               = [snowflake_user.test-user1.name]
}
`, user1, policyName)
}

func networkPolicyAttachmentConfig(user1, user2, policyName string) string {
	return fmt.Sprintf(`
resource "snowflake_user" "test-user1" {
	name = "%s"
}

resource "snowflake_user" "test-user2" {
	name = "%s"
}

resource "snowflake_network_policy" "test" {
	name            = "%v"
	allowed_ip_list = ["192.168.0.100/24", "29.254.123.20"]
}

resource "snowflake_network_policy_attachment" "test" {
	network_policy_name = snowflake_network_policy.test.name
	set_for_account     = false
	users               = [snowflake_user.test-user1.name, snowflake_user.test-user2.name]
}
`, user1, user2, policyName)
}
