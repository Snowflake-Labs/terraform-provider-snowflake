package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicyAttachmentUser(t *testing.T) {
	user1 := acc.TestClient().Ids.Alpha()
	user2 := acc.TestClient().Ids.Alpha()
	policyNameUser := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: networkPolicyAttachmentConfigSingle(user1, policyNameUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyNameUser),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "false"),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "users.#", "1"),
				),
			},
			{
				Config: networkPolicyAttachmentConfig(user1, user2, policyNameUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyNameUser),
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

func TestAcc_NetworkPolicyAttachmentAccount(t *testing.T) {
	policyNameAccount := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckNetworkPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: networkPolicyAttachmentConfigAccount(policyNameAccount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "network_policy_name", policyNameAccount),
					resource.TestCheckResourceAttr("snowflake_network_policy_attachment.test", "set_for_account", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkPolicyAttachmentDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_network_policy_attachment" {
			continue
		}
		ctx := context.Background()
		parameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterNetworkPolicy)
		if err != nil {
			fmt.Printf("[WARN] network policy (%s) not found on account", rs.Primary.Attributes["Id"])
			return nil
		}
		if err == nil && parameter.Level == "ACCOUNT" && parameter.Key == "NETWORK_POLICY" && parameter.Value == rs.Primary.Attributes["network_policy_name"] {
			return fmt.Errorf("network policy attachment %v still exists", rs.Primary.Attributes["Id"])
		}
	}
	return nil
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

func networkPolicyAttachmentConfigAccount(policyName string) string {
	return fmt.Sprintf(`
resource "snowflake_network_policy" "test" {
	name            = "%v"
	allowed_ip_list = ["0.0.0.0/0"]
}

resource "snowflake_network_policy_attachment" "test" {
	network_policy_name = snowflake_network_policy.test.name
	set_for_account     = true
}
`, policyName)
}
