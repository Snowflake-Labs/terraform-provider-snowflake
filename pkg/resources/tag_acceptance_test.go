package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tag(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_tag.t"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":           config.StringVariable(name),
			"database":       config.StringVariable(acc.TestDatabaseName),
			"schema":         config.StringVariable(acc.TestSchemaName),
			"comment":        config.StringVariable("Terraform acceptance test"),
			"allowed_values": config.ListVariable(config.StringVariable("")),
		}
	}

	variableSet2 := m()
	variableSet2["allowed_values"] = config.ListVariable(config.StringVariable("alv1"), config.StringVariable("alv2"))

	variableSet3 := m()
	variableSet3["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Tag/basic"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "allowed_values.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allowed_values.0", ""),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
				),
			},

			// test - change allowed values
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Tag/basic"),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "allowed_values.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "allowed_values.0", "alv1"),
					resource.TestCheckResourceAttr(resourceName, "allowed_values.1", "alv2"),
				),
			},

			// test - change comment
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Tag/basic"),
				ConfigVariables: variableSet3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},

			// test - import
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_Tag/basic"),
				ConfigVariables:   variableSet3,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
