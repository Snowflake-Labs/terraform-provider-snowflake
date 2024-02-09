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

func testAccFunction(t *testing.T, configDirectory string) {
	t.Helper()

	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_function.f"
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
					resource.TestCheckResourceAttrSet(resourceName, "is_secure"),
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
					"language",
					"null_input_behavior",
					"return_behavior",
				},
			},
		},
	})
}

func TestAcc_Function_Javascript(t *testing.T) {
	testAccFunction(t, "TestAcc_Function/javascript")
}

func TestAcc_Function_SQL(t *testing.T) {
	testAccFunction(t, "TestAcc_Function/sql")
}

func TestAcc_Function_Java(t *testing.T) {
	testAccFunction(t, "TestAcc_Function/java")
}

func TestAcc_Function_Scala(t *testing.T) {
	testAccFunction(t, "TestAcc_Function/scala")
}

/*
 Error: 391528 (42601): SQL compilation error: An active warehouse is required for creating Python UDFs.
func TestAcc_Function_Python(t *testing.T) {
	testAccFunction(t, "TestAcc_Function/python")
}
*/

func TestAcc_Function_complex(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_function.f"
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

	statement := "\t\tif (D <= 0) {\n\t\t\treturn 1;\n\t\t} else {\n\t\t\tvar result = 1;\n\t\t\tfor (var i = 2; i <= D; i++) {\n\t\t\t\tresult = result * i;\n\t\t\t}\n\t\t\treturn result;\n\t\t}\n"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDynamicTableDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Function/complex"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "statement", statement),
					resource.TestCheckResourceAttr(resourceName, "arguments.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "arguments.0.name", "D"),
					resource.TestCheckResourceAttr(resourceName, "arguments.0.type", "FLOAT"),
					resource.TestCheckResourceAttr(resourceName, "return_behavior", "VOLATILE"),
					resource.TestCheckResourceAttr(resourceName, "return_type", "FLOAT"),
					resource.TestCheckResourceAttr(resourceName, "language", "javascript"),
					resource.TestCheckResourceAttr(resourceName, "null_input_behavior", "CALLED ON NULL INPUT"),

					// computed attributes
					resource.TestCheckResourceAttrSet(resourceName, "return_type"),
					resource.TestCheckResourceAttrSet(resourceName, "statement"),
					resource.TestCheckResourceAttrSet(resourceName, "is_secure"),
				),
			},

			// test - change comment
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Function/complex"),
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
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_Function/complex"),
				ConfigVariables:   variableSet2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"language",
				},
			},
		},
	})
}
