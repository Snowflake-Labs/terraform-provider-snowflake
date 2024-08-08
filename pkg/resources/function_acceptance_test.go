package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func testAccFunction(t *testing.T, configDirectory string) {
	t.Helper()

	name := acc.TestClient().Ids.Alpha()
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
	name := acc.TestClient().Ids.Alpha()
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
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeVARCHAR)
	name := id.Name()
	comment := random.Comment()
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
				Config: functionConfig(acc.TestDatabaseName, acc.TestSchemaName, name, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|%s|%s|VARCHAR", acc.TestDatabaseName, acc.TestSchemaName, name)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
				),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: functionConfig(acc.TestDatabaseName, acc.TestSchemaName, name, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
				),
			},
		},
	})
}

func TestAcc_Function_Rename(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	newName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()
	newComment := random.Comment()
	resourceName := "snowflake_function.f"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Function),
		Steps: []resource.TestStep{
			{
				Config: functionConfig(acc.TestDatabaseName, acc.TestSchemaName, name, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: functionConfig(acc.TestDatabaseName, acc.TestSchemaName, newName, newComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", newName),
					resource.TestCheckResourceAttr(resourceName, "comment", newComment),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func functionConfig(database string, schema string, name string, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_function" "f" {
  database        = "%[1]s"
  schema          = "%[2]s"
  name            = "%[3]s"
  comment         = "%[4]s"
  return_type     = "VARCHAR"
  return_behavior = "IMMUTABLE"
  statement       = "SELECT PARAM"

  arguments {
    name = "PARAM"
    type = "VARCHAR"
  }
}
`, database, schema, name, comment)
}

// TODO [SNOW-1348103]: do not trim the data type (e.g. NUMBER(10, 2) -> NUMBER loses the information as shown in this test); finish the test
// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2735
func TestAcc_Function_gh2735(t *testing.T) {
	t.Skipf("Will be fixed with functions redesign in SNOW-1348103")
	name := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_function.f"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Function),
		Steps: []resource.TestStep{
			{
				Config: functionConfigGh2735(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func functionConfigGh2735(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_function" "f" {
  database        = "%[1]s"
  schema          = "%[2]s"
  name            = "%[3]s"
  return_type = "TABLE (NUM1 NUMBER, NUM2 NUMBER(10,2))"

  statement = <<EOT
    SELECT 12,13.4
  EOT
}
`, database, schema, name)
}

// TODO: test new state upgrader

func TestAcc_Function_EnsureSmoothResourceIdMigrationToV0950(t *testing.T) {
	name := acc.TestClient().Ids.RandomAccountObjectIdentifier().Name()
	resourceName := "snowflake_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Function),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: functionConfigWithMoreArguments(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%s"."%s"."%s"(VARCHAR, FLOAT, NUMBER)`, acc.TestDatabaseName, acc.TestSchemaName, name)),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   functionConfigWithMoreArguments(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%s"."%s"."%s"(VARCHAR, FLOAT, NUMBER)`, acc.TestDatabaseName, acc.TestSchemaName, name)),
				),
			},
		},
	})
}

func functionConfigWithMoreArguments(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_function" "f" {
  database        = "%[1]s"
  schema          = "%[2]s"
  name            = "%[3]s"
  return_type     = "VARCHAR"
  return_behavior = "IMMUTABLE"
  statement       = "SELECT A"

  arguments {
    name = "A"
    type = "VARCHAR"
  }
  arguments {
    name = "B"
    type = "FLOAT"
  }
  arguments {
    name = "C"
    type = "NUMBER"
  }
}
`, database, schema, name)
}

func TestAcc_Function_EnsureSmoothResourceIdMigrationToV0950_WithoutArguments(t *testing.T) {
	name := acc.TestClient().Ids.RandomAccountObjectIdentifier().Name()
	resourceName := "snowflake_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Function),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: functionConfigWithoutArguments(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%s"."%s"."%s"`, acc.TestDatabaseName, acc.TestSchemaName, name)),
				),
			},
			{
				// TODO: Fails
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   functionConfigWithoutArguments(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%s"."%s"."%s"()`, acc.TestDatabaseName, acc.TestSchemaName, name)),
				),
			},
		},
	})
}

func functionConfigWithoutArguments(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_function" "f" {
  database        = "%[1]s"
  schema          = "%[2]s"
  name            = "%[3]s"
  return_type     = "VARCHAR"
  return_behavior = "IMMUTABLE"
  statement       = "SELECT 'abc'"
}
`, database, schema, name)
}
