package resources_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicy_Basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
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
			// create with empty optionals
			{
				Config: networkPolicyConfigBasic(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_ip_list", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_ip_list", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_network_rules", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_network_rules", "0"),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "describe_output.#", "1"),
					resource.TestCheckNoResourceAttr("snowflake_network_policy.test", "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckNoResourceAttr("snowflake_network_policy.test", "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckNoResourceAttr("snowflake_network_policy.test", "describe_output.0.allowed_ip_list"),
					resource.TestCheckNoResourceAttr("snowflake_network_policy.test", "describe_output.0.blocked_ip_list"),
				),
			},
			// import - without optionals
			{
				Config:       networkPolicyConfigBasic(id.Name()),
				ResourceName: "snowflake_network_policy.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrNotInInstanceState(id.Name(), "allowed_ip_list"),
					importchecks.TestCheckResourceAttrNotInInstanceState(id.Name(), "blocked_ip_list"),
					importchecks.TestCheckResourceAttrNotInInstanceState(id.Name(), "allowed_network_rule_list"),
					importchecks.TestCheckResourceAttrNotInInstanceState(id.Name(), "blocked_network_rule_list"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", ""),
				),
			},
			// set optionals
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
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_network_policy.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_network_rules", "2"),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.blocked_ip_list"),
				),
			},
			// import - complete
			{
				Config: networkPolicyConfigComplete(
					id.Name(),
					[]string{allowedNetworkRuleId1.FullyQualifiedName(), allowedNetworkRuleId2.FullyQualifiedName()},
					[]string{blockedNetworkRuleId1.FullyQualifiedName(), blockedNetworkRuleId2.FullyQualifiedName()},
					[]string{"1.1.1.1", "2.2.2.2"},
					[]string{"3.3.3.3", "4.4.4.4"},
					comment,
				),
				ResourceName: "snowflake_network_policy.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_ip_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "blocked_ip_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_network_rule_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "blocked_network_rule_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", comment),
				),
			},
			// change externally
			{
				PreConfig: func() {
					acc.TestClient().NetworkPolicy.Update(t, sdk.NewAlterNetworkPolicyRequest(id).WithUnset(
						*sdk.NewNetworkPolicyUnsetRequest().
							WithAllowedIpList(true).
							WithBlockedIpList(true).
							WithAllowedNetworkRuleList(true).
							WithBlockedNetworkRuleList(true).
							WithComment(true),
					))
				},
				Config: networkPolicyConfigComplete(
					id.Name(),
					[]string{allowedNetworkRuleId1.FullyQualifiedName(), allowedNetworkRuleId2.FullyQualifiedName()},
					[]string{blockedNetworkRuleId1.FullyQualifiedName(), blockedNetworkRuleId2.FullyQualifiedName()},
					[]string{"1.1.1.1", "2.2.2.2"},
					[]string{"3.3.3.3", "4.4.4.4"},
					comment,
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_network_policy.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_network_rules", "2"),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.blocked_ip_list"),
				),
			},
			// unset
			{
				Config: networkPolicyConfigBasic(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_network_policy.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_ip_list", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_ip_list", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_network_rules", "0"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_network_rules", "0"),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "describe_output.#", "1"),
					resource.TestCheckNoResourceAttr("snowflake_network_policy.test", "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckNoResourceAttr("snowflake_network_policy.test", "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckNoResourceAttr("snowflake_network_policy.test", "describe_output.0.allowed_ip_list"),
					resource.TestCheckNoResourceAttr("snowflake_network_policy.test", "describe_output.0.blocked_ip_list"),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_Complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
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
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_ip_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_network_rules", "2"),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.blocked_ip_list"),
				),
			},
			{
				Config: networkPolicyConfigComplete(
					id.Name(),
					[]string{allowedNetworkRuleId1.FullyQualifiedName(), allowedNetworkRuleId2.FullyQualifiedName()},
					[]string{blockedNetworkRuleId1.FullyQualifiedName(), blockedNetworkRuleId2.FullyQualifiedName()},
					[]string{"1.1.1.1", "2.2.2.2"},
					[]string{"3.3.3.3", "4.4.4.4"},
					comment,
				),
				ResourceName: "snowflake_network_policy.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_ip_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "blocked_ip_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_network_rule_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "blocked_network_rule_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", comment),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_Rename(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config: networkPolicyConfigBasic(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.name", id.Name()),
				),
			},
			{
				Config: networkPolicyConfigBasic(newId.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_network_policy.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "id", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.name", newId.Name()),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_InvalidBlockedIpListValue(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config:      networkPolicyConfigInvalidBlockedIpListValue(id.Name()),
				ExpectError: regexp.MustCompile(`invalid value \(0.0.0.0/0\) set for a field \[{{} blocked_ip_list} {{} {{{{}`),
			},
		},
	})
}

func TestAcc_NetworkPolicy_InvalidNetworkRuleIds(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config:      networkPolicyConfigInvalidAllowedNetworkRules(id.Name()),
				ExpectError: regexp.MustCompile(`sdk\.TableColumnIdentifier\. The correct form of the fully qualified name for`),
			},
			{
				Config:      networkPolicyConfigInvalidAllowedNetworkRules(id.Name()),
				ExpectError: regexp.MustCompile(`sdk\.DatabaseObjectIdentifier\. The correct form of the fully qualified name`),
			},
			{
				Config:      networkPolicyConfigInvalidBlockedNetworkRules(id.Name()),
				ExpectError: regexp.MustCompile(`sdk\.TableColumnIdentifier\. The correct form of the fully qualified name for`),
			},
			{
				Config:      networkPolicyConfigInvalidBlockedNetworkRules(id.Name()),
				ExpectError: regexp.MustCompile(`sdk\.DatabaseObjectIdentifier\. The correct form of the fully qualified name`),
			},
		},
	})
}

func networkPolicyConfigBasic(name string) string {
	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%v"
	}`, name)
}

func networkPolicyConfigInvalidBlockedIpListValue(name string) string {
	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%v"
		blocked_ip_list = ["1.1.1.1", "0.0.0.0/0"]
	}`, name)
}

func networkPolicyConfigInvalidAllowedNetworkRules(name string) string {
	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%v"
		allowed_network_rule_list = ["a.b", "a.b.c.d"]
	}`, name)
}

func networkPolicyConfigInvalidBlockedNetworkRules(name string) string {
	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%v"
		blocked_network_rule_list = ["a.b", "a.b.c.d"]
	}`, name)
}

func networkPolicyConfigComplete(
	name string,
	allowedRuleList []string,
	blockedRuleList []string,
	allowedIpList []string,
	blockedIpList []string,
	comment string,
) string {
	allowedRuleListBytes, _ := json.Marshal(allowedRuleList)
	blockedRuleListBytes, _ := json.Marshal(blockedRuleList)
	allowedIpListBytes, _ := json.Marshal(allowedIpList)
	blockedIpListBytes, _ := json.Marshal(blockedIpList)

	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%[1]s"
		allowed_network_rule_list = %[2]s
		blocked_network_rule_list = %[3]s
		allowed_ip_list = %[4]s
		blocked_ip_list = %[5]s
		comment = "%[6]s"
	}`,
		name,
		string(allowedRuleListBytes),
		string(blockedRuleListBytes),
		string(allowedIpListBytes),
		string(blockedIpListBytes),
		comment,
	)
}

func TestAcc_NetworkPolicy_Issue2236(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	allowedNetworkRuleId := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix("ALLOWED")
	allowedNetworkRuleId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix("ALLOWED")
	blockedNetworkRuleId := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix("BLOCKED")
	blockedNetworkRuleId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix("BLOCKED")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.93.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				// Identifier quoting mismatch (no diff suppression)
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					acc.TestClient().NetworkRule.CreateWithIdentifier(t, allowedNetworkRuleId)
					acc.TestClient().NetworkRule.CreateWithIdentifier(t, allowedNetworkRuleId2)
					acc.TestClient().NetworkRule.CreateWithIdentifier(t, blockedNetworkRuleId)
					acc.TestClient().NetworkRule.CreateWithIdentifier(t, blockedNetworkRuleId2)
				},
				Config: networkPolicyConfigWithNetworkRules(
					id.Name(),
					[]string{
						fmt.Sprintf("\"%s\".\"%s\".%s", allowedNetworkRuleId.DatabaseName(), allowedNetworkRuleId.SchemaName(), allowedNetworkRuleId.Name()),
						fmt.Sprintf("\"%s\".\"%s\".%s", allowedNetworkRuleId2.DatabaseName(), allowedNetworkRuleId2.SchemaName(), allowedNetworkRuleId2.Name()),
					},
					[]string{
						fmt.Sprintf("\"%s\".\"%s\".%s", blockedNetworkRuleId.DatabaseName(), blockedNetworkRuleId.SchemaName(), blockedNetworkRuleId.Name()),
						fmt.Sprintf("\"%s\".\"%s\".%s", blockedNetworkRuleId2.DatabaseName(), blockedNetworkRuleId2.SchemaName(), blockedNetworkRuleId2.Name()),
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config: networkPolicyConfigWithNetworkRules(
					id.Name(),
					[]string{
						fmt.Sprintf("\"%s\".\"%s\".%s", allowedNetworkRuleId.DatabaseName(), allowedNetworkRuleId.SchemaName(), allowedNetworkRuleId.Name()),
						fmt.Sprintf("\"%s\".\"%s\".%s", allowedNetworkRuleId2.DatabaseName(), allowedNetworkRuleId2.SchemaName(), allowedNetworkRuleId2.Name()),
					},
					[]string{
						fmt.Sprintf("\"%s\".\"%s\".%s", blockedNetworkRuleId.DatabaseName(), blockedNetworkRuleId.SchemaName(), blockedNetworkRuleId.Name()),
						fmt.Sprintf("\"%s\".\"%s\".%s", blockedNetworkRuleId2.DatabaseName(), blockedNetworkRuleId2.SchemaName(), blockedNetworkRuleId2.Name()),
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "blocked_network_rule_list.#", "2"),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "show_output.0.entries_in_blocked_network_rules", "2"),

					resource.TestCheckResourceAttr("snowflake_network_policy.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet("snowflake_network_policy.test", "describe_output.0.blocked_network_rule_list"),
				),
			},
		},
	})
}

func networkPolicyConfigWithNetworkRules(name string, allowedNetworkRules []string, blockedNetworkRules []string) string {
	allowedRuleListBytes, _ := json.Marshal(allowedNetworkRules)
	blockedRuleListBytes, _ := json.Marshal(blockedNetworkRules)

	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%[1]s"
		allowed_network_rule_list = %[2]s
		blocked_network_rule_list = %[3]s
	}`, name, string(allowedRuleListBytes), string(blockedRuleListBytes))
}
