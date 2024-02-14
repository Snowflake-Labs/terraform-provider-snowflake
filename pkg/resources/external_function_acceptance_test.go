package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalFunction_basic(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_EXTERNAL_FUNCTION_TESTS"); ok {
		t.Skip("Skipping TestAcc_ExternalFunction")
	}
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
		CheckDestroy: nil,
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
	if _, ok := os.LookupEnv("SKIP_EXTERNAL_FUNCTION_TESTS"); ok {
		t.Skip("Skipping TestAcc_ExternalFunction")
	}
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
		CheckDestroy: nil,
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
	if _, ok := os.LookupEnv("SKIP_EXTERNAL_FUNCTION_TESTS"); ok {
		t.Skip("Skipping TestAcc_ExternalFunction")
	}
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
		CheckDestroy: nil,
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
