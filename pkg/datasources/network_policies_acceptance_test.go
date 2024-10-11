package datasources_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicies_Complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	id2 := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	allowedNetworkRuleId1 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	allowedNetworkRuleId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	blockedNetworkRuleId1 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	blockedNetworkRuleId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acc.TestClient().NetworkRule.CreateWithIdentifier(t, allowedNetworkRuleId1)
					acc.TestClient().NetworkRule.CreateWithIdentifier(t, allowedNetworkRuleId2)
					acc.TestClient().NetworkRule.CreateWithIdentifier(t, blockedNetworkRuleId1)
					acc.TestClient().NetworkRule.CreateWithIdentifier(t, blockedNetworkRuleId2)
				},
				Config: networkPolicyConfigComplete(
					id.Name(),
					[]string{allowedNetworkRuleId1.FullyQualifiedName(), allowedNetworkRuleId2.FullyQualifiedName()},
					[]string{blockedNetworkRuleId1.FullyQualifiedName(), blockedNetworkRuleId2.FullyQualifiedName()},
					[]string{"1.1.1.1", "2.2.2.2"},
					[]string{"3.3.3.3", "4.4.4.4"},
					comment,
					id2.Name(),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_network_policies.test", "network_policies.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_blocked_network_rules", "2"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.comment", comment),

					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.blocked_ip_list"),
					resource.TestCheckResourceAttrSet("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.blocked_network_rule_list"),
				),
			},
			{
				Config: networkPolicyConfigBasic(id.Name(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_network_policies.test", "network_policies.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_allowed_ip_list", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_blocked_ip_list", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_allowed_network_rules", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_blocked_network_rules", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.comment", ""),

					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.#", "1"),
					resource.TestCheckNoResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.allowed_ip_list"),
					resource.TestCheckNoResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.blocked_ip_list"),
					resource.TestCheckNoResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.allowed_network_rule_list"),
					resource.TestCheckNoResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.blocked_network_rule_list"),
				),
			},
			{
				Config: networkPolicyConfigBasic(id2.Name(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_network_policies.test", "network_policies.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.name", id2.Name()),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_allowed_ip_list", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_blocked_ip_list", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_allowed_network_rules", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_blocked_network_rules", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.comment", ""),

					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.#", "1"),
					resource.TestCheckNoResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.allowed_ip_list"),
					resource.TestCheckNoResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.blocked_ip_list"),
					resource.TestCheckNoResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.allowed_network_rule_list"),
					resource.TestCheckNoResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.0.blocked_network_rule_list"),
				),
			},
			{
				Config: networkPolicyConfigBasic(id.Name(), false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_network_policies.test", "network_policies.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_allowed_ip_list", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_blocked_ip_list", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_allowed_network_rules", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.entries_in_blocked_network_rules", "0"),
					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.show_output.0.comment", ""),

					resource.TestCheckResourceAttr("data.snowflake_network_policies.test", "network_policies.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func networkPolicyConfigBasic(name string, withDescribe bool) string {
	return fmt.Sprintf(`
	resource "snowflake_network_policy" "test" {
		name = "%v"
	}

	data "snowflake_network_policies" "test" {
		with_describe = %t
		like = snowflake_network_policy.test.name
	}
`, name, withDescribe)
}

func networkPolicyConfigComplete(
	name string,
	allowedRuleList []string,
	blockedRuleList []string,
	allowedIpList []string,
	blockedIpList []string,
	comment string,
	name2 string,
) string {
	allowedRuleListBytes, _ := json.Marshal(allowedRuleList)
	blockedRuleListBytes, _ := json.Marshal(blockedRuleList)
	allowedIpListBytes, _ := json.Marshal(allowedIpList)
	blockedIpListBytes, _ := json.Marshal(blockedIpList)

	return fmt.Sprintf(`
	resource "snowflake_network_policy" "test" {
		name = "%[1]s"
		allowed_network_rule_list = %[2]s
		blocked_network_rule_list = %[3]s
		allowed_ip_list = %[4]s
		blocked_ip_list = %[5]s
		comment = "%[6]s"
	}

	resource "snowflake_network_policy" "test2" {
		name = "%[7]s"
	}

	data "snowflake_network_policies" "test" {
		like = snowflake_network_policy.test.name
	}
`,
		name,
		string(allowedRuleListBytes),
		string(blockedRuleListBytes),
		string(allowedIpListBytes),
		string(blockedIpListBytes),
		comment,
		name2,
	)
}

func TestAcc_NetworkPolicies_NetworkPolicyNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      networkPolicyConfigBasicWithPostConditions(),
				ExpectError: regexp.MustCompile("there should be at least one network policy"),
			},
		},
	})
}

func networkPolicyConfigBasicWithPostConditions() string {
	return `
	data "snowflake_network_policies" "test" {
		like = "non_existing_network_policy"
	  	lifecycle {
			postcondition {
		  		condition     = length(self.network_policies) > 0
		  		error_message = "there should be at least one network policy"
			}
	  	}
	}
	`
}
