package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RowAccessPolicy(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "body", "case when current_role() in ('ANALYST') then true else false end"),
					resource.TestCheckResourceAttr(resourceName, "argument.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.name", "N"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.type", "VARCHAR"),
					resource.TestCheckResourceAttr(resourceName, "argument.1.name", "V"),
					resource.TestCheckResourceAttr(resourceName, "argument.1.type", "VARCHAR"),
				),
			},
			// change comment and expression
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - changed comment"),
					resource.TestCheckResourceAttr(resourceName, "body", "case when current_role() in ('ANALYST') then false else true end"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.name", "N"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.type", "VARCHAR"),
					resource.TestCheckResourceAttr(resourceName, "argument.1.name", "V"),
					resource.TestCheckResourceAttr(resourceName, "argument.1.type", "VARCHAR"),
				),
			},
			// change signature
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - changed comment"),
					resource.TestCheckResourceAttr(resourceName, "body", "case when current_role() in ('ANALYST') then false else true end"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.name", "V"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.type", "BOOLEAN"),
					resource.TestCheckResourceAttr(resourceName, "argument.1.name", "X"),
					resource.TestCheckResourceAttr(resourceName, "argument.1.type", "TIMESTAMP_NTZ"),
				),
			},
			// external change on body
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				PreConfig: func() {
					acc.TestClient().RowAccessPolicy.Alter(t, *sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.Pointer("case when current_role() in ('EXTERNAL') then false else true end")))
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - changed comment"),
					resource.TestCheckResourceAttr(resourceName, "body", "case when current_role() in ('ANALYST') then false else true end"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.name", "V"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.type", "BOOLEAN"),
					resource.TestCheckResourceAttr(resourceName, "argument.1.name", "X"),
					resource.TestCheckResourceAttr(resourceName, "argument.1.type", "TIMESTAMP_NTZ"),
				),
			},
			// IMPORT
			{
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2053 is fixed
func TestAcc_RowAccessPolicy_Issue2053(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
			"arguments": config.SetVariable(
				config.MapVariable(map[string]config.Variable{
					"name": config.StringVariable("A"),
					"type": config.StringVariable("VARCHAR"),
				}),
			),
		}
	}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.95.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: rowAccessPolicy_v0_95_0_WithHeredoc(id, `    case
      when current_role() in ('ANALYST') then true
      else false
    end
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables:          m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "body", `case
  when current_role() in ('ANALYST') then true
  else false
end`),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_Rename(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	m := func(identifier sdk.SchemaObjectIdentifier) config.Variables {
		return config.Variables{
			"name":     config.StringVariable(identifier.Name()),
			"database": config.StringVariable(identifier.DatabaseName()),
			"schema":   config.StringVariable(identifier.SchemaName()),
			"arguments": config.SetVariable(
				config.MapVariable(map[string]config.Variable{
					"name": config.StringVariable("A"),
					"type": config.StringVariable("VARCHAR"),
				}),
			),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: m(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
			// rename with one param changed
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: m(newId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", newId.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", newId.FullyQualifiedName()),
				),
			},
		},
	})
}

func rowAccessPolicy_v0_95_0(id sdk.SchemaObjectIdentifier, expr string) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
  name     = "%s"
  database = "%s"
  schema   = "%s"
  signature = {
    A = "VARCHAR",
  }
  row_access_expression = "%s"
}`, id.Name(), id.DatabaseName(), id.SchemaName(), expr)
}

func rowAccessPolicy_v0_95_0_WithHeredoc(id sdk.SchemaObjectIdentifier, expr string) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
  name     = "%s"
  database = "%s"
  schema   = "%s"
  signature = {
    A = "VARCHAR",
  }
  row_access_expression = <<-EOT
%s
EOT
}`, id.Name(), id.DatabaseName(), id.SchemaName(), expr)
}

func rowAccessPolicy_v0_96_0(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
  name     = "%s"
  database = "%s"
  schema   = "%s"
  argument {
    name = "A"
    type = "VARCHAR"
  }
  row_access_expression = <<-EOT
    case
      when current_role() in ('ANALYST') then true
      else false
    end
EOT
}`, id.Name(), id.DatabaseName(), id.SchemaName())
}

func TestAcc_RowAccessPolicy_InvalidDataType(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
			"arguments": config.SetVariable(
				config.MapVariable(map[string]config.Variable{
					"name": config.StringVariable("A"),
					"type": config.StringVariable("invalid-int"),
				}),
			),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: m(),
				ExpectError:     regexp.MustCompile(`invalid data type: invalid-int`),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_DataTypeAliases(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
			"arguments": config.SetVariable(
				config.MapVariable(map[string]config.Variable{
					"name": config.StringVariable("A"),
					"type": config.StringVariable("TEXT"),
				}),
			),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "argument.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.type", "VARCHAR")),
			},
		},
	})
}

func TestAcc_view_migrateFromVersion_0_95_0(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	m := config.Variables{
		"name":     config.StringVariable(id.Name()),
		"database": config.StringVariable(id.DatabaseName()),
		"schema":   config.StringVariable(id.SchemaName()),
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.95.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: rowAccessPolicy_v0_95_0(id, "case when current_role() in ('ANALYST') then true else false end"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "row_access_expression", "case when current_role() in ('ANALYST') then true else false end"),
					resource.TestCheckResourceAttr(resourceName, "signature.A", "VARCHAR"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables:          m,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "body", "case when current_role() in ('ANALYST') then true else false end"),
					resource.TestCheckNoResourceAttr(resourceName, "row_access_expression"),
					resource.TestCheckResourceAttr(resourceName, "argument.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.name", "A"),
					resource.TestCheckResourceAttr(resourceName, "argument.0.type", "VARCHAR"),
					resource.TestCheckNoResourceAttr(resourceName, "signature.A"),
				),
			},
		},
	})
}
