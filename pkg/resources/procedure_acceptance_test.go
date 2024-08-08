package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func testAccProcedure(t *testing.T, configDirectory string) {
	t.Helper()

	name := acc.TestClient().Ids.Alpha()
	newName := acc.TestClient().Ids.Alpha()

	resourceName := "snowflake_procedure.p"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(name),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"comment":    config.StringVariable("Terraform acceptance test"),
			"execute_as": config.StringVariable("CALLER"),
		}
	}
	variableSet2 := m()
	variableSet2["name"] = config.StringVariable(newName)
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")
	variableSet2["execute_as"] = config.StringVariable("OWNER")

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
		CheckDestroy: acc.CheckDestroy(t, resources.Procedure),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory(configDirectory),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "return_behavior", "VOLATILE"),
					resource.TestCheckResourceAttr(resourceName, "execute_as", "CALLER"),

					// computed attributes
					resource.TestCheckResourceAttrSet(resourceName, "return_type"),
					resource.TestCheckResourceAttrSet(resourceName, "statement"),
					resource.TestCheckResourceAttrSet(resourceName, "secure"),
				),
			},

			// test - rename + change comment and caller (proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2642)
			{
				ConfigDirectory: acc.ConfigurationDirectory(configDirectory),
				ConfigVariables: variableSet2,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", newName),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
					resource.TestCheckResourceAttr(resourceName, "execute_as", "OWNER"),
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
	name := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_procedure.p"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(name),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"comment":    config.StringVariable("Terraform acceptance test"),
			"execute_as": config.StringVariable("CALLER"),
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
		CheckDestroy: acc.CheckDestroy(t, resources.Procedure),
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

func TestAcc_Procedure_migrateFromVersion085(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	resourceName := "snowflake_procedure.p"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Procedure),

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: procedureConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|%s|%s|", acc.TestDatabaseName, acc.TestSchemaName, name)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   procedureConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
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

func procedureConfig(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_procedure" "p" {
 database    = "%[1]s"
 schema      = "%[2]s"
 name        = "%[3]s"
 language    = "JAVASCRIPT"
 return_type = "VARCHAR"
 statement   = <<EOT
   return "Hi"
 EOT
}
`, database, schema, name)
}

func TestAcc_Procedure_proveArgsPermanentDiff(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeVARCHAR, sdk.DataTypeNumber)
	name := id.Name()
	resourceName := "snowflake_procedure.p"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Procedure),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.89.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: sqlProcedureConfigArgsPermanentDiff(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id.FullyQualifiedName()),
				),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
				},
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   sqlProcedureConfigArgsPermanentDiff(acc.TestDatabaseName, acc.TestSchemaName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

// TODO [SNOW-1348106]: diff suppression for the return type (the same with functions); finish this test
func TestAcc_Procedure_returnTypePermanentDiff(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeVARCHAR)
	name := id.Name()
	resourceName := "snowflake_procedure.p"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Procedure),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.89.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: sqlProcedureConfigReturnTypePermanentDiff(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id.FullyQualifiedName()),
				),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate)},
				},
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   sqlProcedureConfigReturnTypePermanentDiff(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id.FullyQualifiedName()),
				),
				// should be empty after SNOW-1348106
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate)},
				},
			},
		},
	})
}

func sqlProcedureConfigArgsPermanentDiff(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_procedure" "p" {
 database    = "%[1]s"
 schema      = "%[2]s"
 name        = "%[3]s"
 language    = "SQL"
 return_type = "NUMBER(38,0)"
 arguments {
   name = "arg1"
   type = "VARCHAR"
 }
 arguments {
   name = "MY_INT"
   type = "int"
 }
 statement   = <<EOT
BEGIN
 RETURN 13.4;
END;
 EOT
}
`, database, schema, name)
}

func sqlProcedureConfigReturnTypePermanentDiff(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_procedure" "p" {
 database    = "%[1]s"
 schema      = "%[2]s"
 name        = "%[3]s"
 language    = "SQL"
 return_type = "TABLE (NUM1 NUMBER(10,2))"
 arguments {
   name = "ARG1"
   type = "VARCHAR"
 }
 statement   = <<EOT
BEGIN
 RETURN 13.4;
END;
 EOT
}
`, database, schema, name)
}
