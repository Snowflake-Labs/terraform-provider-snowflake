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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalFunction_basic(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"schema":                    config.StringVariable(acc.TestSchemaName),
			"name":                      config.StringVariable(accName),
			"api_allowed_prefixes":      config.ListVariable(config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/")),
			"url_of_proxy_and_resource": config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
			"comment":                   config.StringVariable("Terraform acceptance test"),
		}
	}

	resourceName := "snowflake_external_function.external_function"
	configVariables := m()
	configVariables2 := m()
	configVariables2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalFunction/basic"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", accName),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "arg.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "arg.0.name", "ARG1"),
					resource.TestCheckResourceAttr(resourceName, "arg.0.type", "VARCHAR"),
					resource.TestCheckResourceAttr(resourceName, "arg.1.name", "ARG2"),
					resource.TestCheckResourceAttr(resourceName, "arg.1.type", "VARCHAR"),
					resource.TestCheckResourceAttr(resourceName, "null_input_behavior", "CALLED ON NULL INPUT"),
					resource.TestCheckResourceAttr(resourceName, "return_type", "VARIANT"),
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName, "return_behavior", "IMMUTABLE"),
					resource.TestCheckResourceAttrSet(resourceName, "api_integration"),
					resource.TestCheckResourceAttr(resourceName, "compression", "AUTO"),
					resource.TestCheckResourceAttr(resourceName, "url_of_proxy_and_resource", "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalFunction/basic"),
				ConfigVariables: configVariables2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},
			// IMPORT
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ExternalFunction/basic"),
				ConfigVariables:   configVariables2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// these two are not found in either the show or describe command
				ImportStateVerifyIgnore: []string{"return_null_allowed", "api_integration"},
			},
		},
	})
}

func TestAcc_ExternalFunction_no_arguments(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"schema":                    config.StringVariable(acc.TestSchemaName),
			"name":                      config.StringVariable(accName),
			"api_allowed_prefixes":      config.ListVariable(config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/")),
			"url_of_proxy_and_resource": config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
			"comment":                   config.StringVariable("Terraform acceptance test"),
		}
	}

	resourceName := "snowflake_external_function.external_function"
	configVariables := m()
	configVariables2 := m()
	configVariables2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalFunction/no_arguments"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", accName),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "arg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "null_input_behavior", "CALLED ON NULL INPUT"),
					resource.TestCheckResourceAttr(resourceName, "return_type", "VARIANT"),
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName, "return_behavior", "IMMUTABLE"),
					resource.TestCheckResourceAttrSet(resourceName, "api_integration"),
					resource.TestCheckResourceAttr(resourceName, "compression", "AUTO"),
					resource.TestCheckResourceAttr(resourceName, "url_of_proxy_and_resource", "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalFunction/no_arguments"),
				ConfigVariables: configVariables2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},
			// IMPORT
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ExternalFunction/no_arguments"),
				ConfigVariables:   configVariables2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// these two are not found in either the show or describe command
				ImportStateVerifyIgnore: []string{"return_null_allowed", "api_integration"},
			},
		},
	})
}

func TestAcc_ExternalFunction_complete(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"schema":                    config.StringVariable(acc.TestSchemaName),
			"name":                      config.StringVariable(accName),
			"api_allowed_prefixes":      config.ListVariable(config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/")),
			"url_of_proxy_and_resource": config.StringVariable("https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
			"comment":                   config.StringVariable("Terraform acceptance test"),
		}
	}

	resourceName := "snowflake_external_function.external_function"
	configVariables := m()
	configVariables2 := m()
	configVariables2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalFunction/complete"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", accName),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "arg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "null_input_behavior", "CALLED ON NULL INPUT"),
					resource.TestCheckResourceAttr(resourceName, "return_type", "VARIANT"),
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName, "return_behavior", "IMMUTABLE"),
					resource.TestCheckResourceAttrSet(resourceName, "api_integration"),
					resource.TestCheckResourceAttr(resourceName, "compression", "AUTO"),
					resource.TestCheckResourceAttr(resourceName, "url_of_proxy_and_resource", "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
					resource.TestCheckResourceAttr(resourceName, "header.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "header.0.name", "x-custom-header"),
					resource.TestCheckResourceAttr(resourceName, "header.0.value", "snowflake"),
					resource.TestCheckResourceAttr(resourceName, "max_batch_rows", "500"),
					resource.TestCheckResourceAttr(resourceName, "request_translator", fmt.Sprintf("%s.%s.%s%s", acc.TestDatabaseName, acc.TestSchemaName, accName, "_request_translator")),
					resource.TestCheckResourceAttr(resourceName, "response_translator", fmt.Sprintf("%s.%s.%s%s", acc.TestDatabaseName, acc.TestSchemaName, accName, "_response_translator")),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalFunction/complete"),
				ConfigVariables: configVariables2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},
			// IMPORT
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ExternalFunction/complete"),
				ConfigVariables:   configVariables2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// these four are not found in either the show or describe command
				ImportStateVerifyIgnore: []string{"return_null_allowed", "api_integration", "request_translator", "response_translator"},
			},
		},
	})
}

func TestAcc_ExternalFunction_migrateFromVersion085(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalFunction),

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: externalFunctionConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|%s|%s|VARCHAR-VARCHAR", acc.TestDatabaseName, acc.TestSchemaName, name)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", "\""+acc.TestDatabaseName+"\""),
					resource.TestCheckResourceAttr(resourceName, "schema", "\""+acc.TestSchemaName+"\""),
					resource.TestCheckNoResourceAttr(resourceName, "return_null_allowed"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   externalFunctionConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", sdk.NewSchemaObjectIdentifierWithArguments(acc.TestDatabaseName, acc.TestSchemaName, name, []sdk.DataType{sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR}).FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
			},
		},
	})
}

func TestAcc_ExternalFunction_migrateFromVersion085_issue2694_previousValuePresent(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalFunction),

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: externalFunctionConfigWithReturnNullAllowed(acc.TestDatabaseName, acc.TestSchemaName, name, sdk.Bool(true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   externalFunctionConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
			},
		},
	})
}

func TestAcc_ExternalFunction_migrateFromVersion085_issue2694_previousValueRemoved(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalFunction),

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: externalFunctionConfigWithReturnNullAllowed(acc.TestDatabaseName, acc.TestSchemaName, name, sdk.Bool(true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: externalFunctionConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "return_null_allowed"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   externalFunctionConfig(acc.TestDatabaseName, acc.TestSchemaName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "return_null_allowed", "true"),
				),
			},
		},
	})
}

// Proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2528.
// The problem originated from ShowById without IN clause. There was no IN clause in the docs at the time.
// It was raised with the appropriate team in Snowflake.
func TestAcc_ExternalFunction_issue2528(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	secondSchema := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resourceName := "snowflake_external_function.f"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalFunction),
		Steps: []resource.TestStep{
			{
				Config: externalFunctionConfigIssue2528(acc.TestDatabaseName, acc.TestSchemaName, accName, secondSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", accName),
				),
			},
		},
	})
}

func externalFunctionConfig(database string, schema string, name string) string {
	return externalFunctionConfigWithReturnNullAllowed(database, schema, name, nil)
}

func externalFunctionConfigWithReturnNullAllowed(database string, schema string, name string, returnNullAllowed *bool) string {
	returnNullAllowedText := ""
	if returnNullAllowed != nil {
		returnNullAllowedText = fmt.Sprintf("return_null_allowed = \"%t\"", *returnNullAllowed)
	}

	return fmt.Sprintf(`
resource "snowflake_api_integration" "test_api_int" {
  name                 = "%[3]s"
  api_provider         = "aws_api_gateway"
  api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
  api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
  enabled              = true
}

resource "snowflake_external_function" "f" {
  name     = "%[3]s"
  database = "%[1]s"
  schema   = "%[2]s"
  arg {
    name = "ARG1"
    type = "VARCHAR"
  }
  arg {
    name = "ARG2"
    type = "VARCHAR"
  }
  return_type               = "VARIANT"
  return_behavior           = "IMMUTABLE"
  api_integration           = snowflake_api_integration.test_api_int.name
  url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
  %[4]s
}

`, database, schema, name, returnNullAllowedText)
}

func externalFunctionConfigIssue2528(database string, schema string, name string, schema2 string) string {
	return fmt.Sprintf(`
resource "snowflake_api_integration" "test_api_int" {
  name                 = "%[3]s"
  api_provider         = "aws_api_gateway"
  api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
  api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
  enabled              = true
}

resource "snowflake_schema" "s2" {
  database            = "%[1]s"
  name                = "%[4]s"
}

resource "snowflake_external_function" "f" {
  name     = "%[3]s"
  database = "%[1]s"
  schema   = "%[2]s"
  arg {
    name = "SNS_NOTIF"
    type = "OBJECT"
  }
  return_type = "VARIANT"
  return_behavior = "VOLATILE"
  api_integration = snowflake_api_integration.test_api_int.name
  url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}

resource "snowflake_external_function" "f2" {
  depends_on = [snowflake_schema.s2]

  name     = "%[3]s"
  database = "%[1]s"
  schema   = "%[4]s"
  arg {
    name = "SNS_NOTIF"
    type = "OBJECT"
  }
  return_type = "VARIANT"
  return_behavior = "VOLATILE"
  api_integration = snowflake_api_integration.test_api_int.name
  url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
}
`, database, schema, name, schema2)
}
