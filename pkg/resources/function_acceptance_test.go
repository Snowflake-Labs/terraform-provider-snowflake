package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	ignoreDuringImport := []string{"null_input_behavior"}
	if strings.Contains(configDirectory, "/sql") {
		ignoreDuringImport = append(ignoreDuringImport, "return_behavior")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Function),
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
				ConfigDirectory:         acc.ConfigurationDirectory(configDirectory),
				ConfigVariables:         variableSet2,
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreDuringImport,
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
		CheckDestroy: acc.CheckDestroy(t, resources.Function),
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
					resource.TestCheckResourceAttr(resourceName, "language", "JAVASCRIPT"),
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

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2490
func TestAcc_Function_migrateFromVersion085(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Function),

		// Using the string config because of the validation in teststep_validate.go:
		// teststep.Config.HasConfigurationFiles() returns true both for ConfigFile and ConfigDirectory.
		// It returns false for Config. I don't understand why they have such a validation, but we will work around it later.
		// Added as subtask SNOW-1057066 to SNOW-926148.
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: functionConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|%s|%s|VARCHAR", acc.TestDatabaseName, acc.TestSchemaName, name)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   functionConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", sdk.NewSchemaObjectIdentifierWithArguments(acc.TestDatabaseName, acc.TestSchemaName, name, []sdk.DataType{sdk.DataTypeVARCHAR}).FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
				),
			},
		},
	})
}

func functionConfig(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_function" "f" {
  database        = "%[1]s"
  schema          = "%[2]s"
  name            = "%[3]s"
  return_type     = "VARCHAR"
  return_behavior = "IMMUTABLE"
  statement       = "SELECT PARAM"

  arguments {
    name = "PARAM"
    type = "VARCHAR"
  }
}
`, database, schema, name)
}
