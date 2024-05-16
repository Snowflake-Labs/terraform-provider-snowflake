package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	networkPolicyComment = "CREATED BY A TERRAFORM ACCEPTANCE TEST"
)

func TestAcc_NetworkPolicy(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	nameRule1 := acc.TestClient().Ids.Alpha()
	nameRule2 := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config: networkPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", networkPolicyComment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "0"),
				),
			},
			// CHANGE PROPERTIES - add to and remove from ip lists
			{
				Config: networkPolicyConfig2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", networkPolicyComment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "0"),
				),
			},
			{
				Config: networkPolicyConfigNetworkRules(name, nameRule1, nameRule2, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", networkPolicyComment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "1"),
				),
			},
			{
				Config: networkPolicyConfigIPsAndRules(name, nameRule1, nameRule2, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", networkPolicyComment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "1"),
				),
			},
			// IMPORT - all fields are non-empty
			{
				ResourceName:      "snowflake_network_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: networkPolicyConfigAllEmpty(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", networkPolicyComment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "0"),
				),
			},
			// IMPORT - incomplete
			{
				ResourceName:      "snowflake_network_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_NetworkPolicyBadNetworkRuleNames(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			// Checks the case when a network rule name, which is not a schema object identifier, is passed to a network policy
			{
				Config:      networkPolicyConfigBadNetworkRule(name),
				ExpectError: regexp.MustCompile(`.*Invalid identifier type.*`),
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

func networkPolicyConfigNetworkRules(name string, nameRule1 string, nameRule2 string, database string, schema string) string {
	return fmt.Sprintf(`
resource "snowflake_network_rule" "test1" {
	name            = "%[2]v"
	database        = "%[4]v"
	schema          = "%[5]v"
	comment         = "%[6]v"
    type            = "IPV4"
    mode			= "INGRESS"
	value_list      = ["192.168.0.100/24", "29.254.123.20"]
}

resource "snowflake_network_rule" "test2" {
	name            = "%[3]v"
	database        = "%[4]v"
	schema          = "%[5]v"
	comment         = "%[6]v"
    type            = "HOST_PORT"
    mode			= "EGRESS"
	value_list      = ["example.com", "company.com:443"]
}

resource "snowflake_network_policy" "test" {
	name            = "%[1]v"
	comment         = "%[6]v"
	allowed_network_rule_list = [snowflake_network_rule.test1.qualified_name]
	blocked_network_rule_list = [snowflake_network_rule.test2.qualified_name]
}
`, name, nameRule1, nameRule2, database, schema, networkPolicyComment)
}

func networkPolicyConfigIPsAndRules(name string, nameRule1 string, nameRule2 string, database string, schema string) string {
	return fmt.Sprintf(`
resource "snowflake_network_rule" "test1" {
	name            = "%[2]v"
	database        = "%[4]v"
	schema          = "%[5]v"
	comment         = "%[6]v"
    type            = "IPV4"
    mode			= "INGRESS"
	value_list      = ["192.168.0.100/24", "29.254.123.20"]
}

resource "snowflake_network_rule" "test2" {
	name            = "%[3]v"
	database        = "%[4]v"
	schema          = "%[5]v"
	comment         = "%[6]v"
    type            = "HOST_PORT"
    mode			= "EGRESS"
	value_list      = ["example.com", "company.com:443"]
}

resource "snowflake_network_policy" "test" {
	name            = "%[1]v"
	comment         = "%[6]v"
	allowed_ip_list = ["192.168.0.100/24"]
	blocked_ip_list = ["192.168.0.101"]
	allowed_network_rule_list = [snowflake_network_rule.test1.qualified_name]
	blocked_network_rule_list = [snowflake_network_rule.test2.qualified_name]
}
`, name, nameRule1, nameRule2, database, schema, networkPolicyComment)
}

func networkPolicyConfigAllEmpty(name string) string {
	return fmt.Sprintf(`
resource "snowflake_network_policy" "test" {
	name            = "%v"
	comment         = "%v"
	allowed_ip_list = []
	blocked_ip_list = []
	allowed_network_rule_list = []
	blocked_network_rule_list = []
}
`, name, networkPolicyComment)
}

func networkPolicyConfigBadNetworkRule(name string) string {
	return fmt.Sprintf(`
resource "snowflake_network_policy" "test" {
	name            = "%v"
	comment         = "%v"
	allowed_network_rule_list = ["FOO"]
}
`, name, networkPolicyComment)
}
