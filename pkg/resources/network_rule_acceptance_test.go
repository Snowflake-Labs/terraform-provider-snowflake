package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	networkRuleComment = "CREATED BY A TERRAFORM ACCEPTANCE TEST"
)

func TestAcc_NetworkRule(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkRule),
		Steps: []resource.TestStep{
			// basic
			{
				Config: networkRuleIpv4(name, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "comment", networkRuleComment),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "type", "IPV4"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "mode", "INGRESS"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "value_list.#", "2"),
				),
			},
			//// IMPORT - all fields are non-empty
			{
				ResourceName:      "snowflake_network_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			//// CHANGE PROPERTIES - set to empty
			{
				Config: networkRuleIpv4Empty(name, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "comment", networkRuleComment),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "type", "IPV4"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "mode", "INGRESS"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "value_list.#", "0"),
				),
			},
			// IMPORT - incomplete
			{
				ResourceName:      "snowflake_network_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: networkRuleHost(name, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "comment", networkRuleComment),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "type", "HOST_PORT"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "mode", "EGRESS"),
					resource.TestCheckResourceAttr("snowflake_network_rule.test", "value_list.#", "2"),
				),
			},
		},
	})
}

func networkRuleIpv4(name string, database string, schema string) string {
	return fmt.Sprintf(`
resource "snowflake_network_rule" "test" {
	name            = "%v"
	database        = "%v"
	schema          = "%v"
	comment         = "%v"
    type            = "IPV4"
    mode			= "INGRESS"
	value_list      = ["192.168.0.100/24", "29.254.123.20"]
}
`, name, database, schema, networkRuleComment)
}

func networkRuleIpv4Empty(name string, database string, schema string) string {
	return fmt.Sprintf(`
resource "snowflake_network_rule" "test" {
	name            = "%v"
	database        = "%v"
	schema          = "%v"
	comment         = "%v"
    type            = "IPV4"
    mode			= "INGRESS"
	value_list      = []
}
`, name, database, schema, networkRuleComment)
}

func networkRuleHost(name string, database string, schema string) string {
	return fmt.Sprintf(`
resource "snowflake_network_rule" "test" {
	name            = "%v"
	database        = "%v"
	schema          = "%v"
	comment         = "%v"
    type            = "HOST_PORT"
    mode			= "EGRESS"
	value_list      = ["example.com", "company.com:443"]
}
`, name, database, schema, networkRuleComment)
}

func TestAcc_NetworkRule_migrateFromVersion_0_94_1(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_network_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: networkRuleIpv4(id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   networkRuleIpv4(id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceName, "qualified_name"),
				),
			},
		},
	})
}
