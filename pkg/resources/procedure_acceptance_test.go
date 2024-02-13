package resources_test

import (
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func testAccProcedure(t *testing.T, configDirectory string) {
	t.Helper()

	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_procedure.p"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(name),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
			"comment":  config.StringVariable("Terraform acceptance test"),
		}
	}
	variableSet2 := m()
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDynamicTableDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory(configDirectory),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),

					// computed attributes
					resource.TestCheckResourceAttrSet(resourceName, "return_type"),
					resource.TestCheckResourceAttrSet(resourceName, "statement"),
					resource.TestCheckResourceAttrSet(resourceName, "execute_as"),
					resource.TestCheckResourceAttrSet(resourceName, "secure"),
				),
			},

			// test - change comment
			{
				ConfigDirectory: acc.ConfigurationDirectory(configDirectory),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},

			// test - import
			{
				ConfigDirectory:   acc.ConfigurationDirectory(configDirectory),
				ConfigVariables:   variableSet2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"null_input_behavior",
					"return_behavior",
				},
			},
		},
	})
}

func TestAcc_Procedure_SQL(t *testing.T) {
	testAccProcedure(t, "TestAcc_Procedure/sql")
}

/*
Error: 391531 (42601): SQL compilation error: An active warehouse is required for creating Python stored procedures.
func TestAcc_Procedure_Python(t *testing.T) {
	testAccProcedure(t, "TestAcc_Procedure/python")
}
*/

func TestAcc_Procedure_Javascript(t *testing.T) {
	testAccProcedure(t, "TestAcc_Procedure/javascript")
}

func TestAcc_Procedure_Java(t *testing.T) {
	testAccProcedure(t, "TestAcc_Procedure/java")
}

func TestAcc_Procedure_Scala(t *testing.T) {
	testAccProcedure(t, "TestAcc_Procedure/scala")
}

func TestAcc_Procedure_complex(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_procedure.p"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(name),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
			"comment":  config.StringVariable("Terraform acceptance test"),
		}
	}
	variableSet2 := m()
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	statement := "var x = 1\nreturn x\n"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDynamicTableDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Procedure/complex"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "statement", statement),
					resource.TestCheckResourceAttr(resourceName, "execute_as", "CALLER"),
					resource.TestCheckResourceAttr(resourceName, "arguments.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "arguments.0.name", "ARG1"),
					resource.TestCheckResourceAttr(resourceName, "arguments.0.type", "VARCHAR"),
					resource.TestCheckResourceAttr(resourceName, "arguments.1.name", "ARG2"),
					resource.TestCheckResourceAttr(resourceName, "arguments.1.type", "DATE"),
					resource.TestCheckResourceAttr(resourceName, "null_input_behavior", "RETURNS NULL ON NULL INPUT"),

					// computed attributes
					resource.TestCheckResourceAttrSet(resourceName, "return_type"),
					resource.TestCheckResourceAttrSet(resourceName, "statement"),
					resource.TestCheckResourceAttrSet(resourceName, "execute_as"),
					resource.TestCheckResourceAttrSet(resourceName, "secure"),
				),
			},

			// test - change comment
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Procedure/complex"),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},

			// test - import
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_Procedure/complex"),
				ConfigVariables:   variableSet2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"return_behavior",
				},
			},
		},
	})
}
