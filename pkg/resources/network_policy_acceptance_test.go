package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	networkPolicyComment = "CREATED BY A TERRAFORM ACCEPTANCE TEST"
)

func TestAcc_NetworkPolicy(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_NETWORK_POLICY_TESTS"); ok {
		t.Skip("Skipping TestAccNetworkPolicy")
	}

	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: networkPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", networkPolicyComment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "0"),
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
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_network_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func networkPolicyConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_network_policy" "test" {
	name            = "%v"
	comment         = "%v"
	allowed_ip_list = ["192.168.0.100/24", "29.254.123.20"]
}
`, name, networkPolicyComment)
}

func networkPolicyConfig2(name string) string {
	return fmt.Sprintf(`
resource "snowflake_network_policy" "test" {
	name            = "%v"
	comment         = "%v"
	allowed_ip_list = ["192.168.0.100/24"]
	blocked_ip_list = ["192.168.0.101"]
}
`, name, networkPolicyComment)
}
