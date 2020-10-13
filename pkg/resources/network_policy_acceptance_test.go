package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	networkPolicyComment = "Created by a Terraform acceptance test"
)

func TestAccNetworkPolicy(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_NETWORK_POLICY_TESTS"); ok {
		t.Skip("Skipping TestAccNetworkPolicy")
	}

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: networkPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", networkPolicyComment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "users.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "set_for_account", "false"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: networkPolicyConfig2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", networkPolicyComment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "users.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "set_for_account", "false"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_network_policy.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"set_for_account", "users"},
			},
		},
	})
}

func networkPolicyConfig(name string) string {
	return fmt.Sprintf(`

resource "snowflake_user" "test_user" {
  name = "TEST_USER"
}

resource "snowflake_network_policy" "test" {
	name            = "%v"
	comment         = "%v"
	allowed_ip_list = ["192.168.0.100/24", "29.254.123.20"]
	blocked_ip_list = ["192.168.0.101"]
	users           = [snowflake_user.test_user.name]
}
`, name, networkPolicyComment)
}

func networkPolicyConfig2(name string) string {
	return fmt.Sprintf(`

resource "snowflake_user" "test_user" {
  name = "TEST_USER"
}

resource "snowflake_network_policy" "test" {
	name            = "%v"
	comment         = "%v"
	allowed_ip_list = ["192.168.0.100/24"]
	blocked_ip_list = ["192.168.0.101"]
	users           = [snowflake_user.test_user.name]
}
`, name, networkPolicyComment)
}
