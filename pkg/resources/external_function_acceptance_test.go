package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ExternalFunction(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_EXTERNAL_FUNCTION_TESTS"); ok {
		t.Skip("Skipping TestAccExternalFunction")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: externalFunctionConfig(accName, []string{"https://123456.execute-api.us-west-2.amazonaws.com/prod/"}, "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_function.test_func", "name", accName),
					resource.TestCheckResourceAttr("snowflake_external_function.test_func", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrSet("snowflake_external_function.test_func", "created_on"),
				),
			},
		},
	})
}

func externalFunctionConfig(name string, prefixes []string, url string) string {
	return fmt.Sprintf(`
	resource "snowflake_database" "test_database" {
		name    = "%s"
		comment = "Terraform acceptance test"
	}
	
	resource "snowflake_schema" "test_schema" {
		name     = "%s"
		database = snowflake_database.test_database.name
		comment  = "Terraform acceptance test"
	}

	resource "snowflake_api_integration" "test_api_int" {
		name = "%s"
		api_provider = "aws_api_gateway"
		api_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
		api_allowed_prefixes = %q
		enabled = true
	}

	resource "snowflake_external_function" "test_func" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		arg {
			name = "arg1"
			type = "varchar"
		}
		arg {
			name = "arg2"
			type = "varchar"
		}
		comment = "Terraform acceptance test"
		return_type = "varchar"
		return_behavior = "IMMUTABLE"
		api_integration = snowflake_api_integration.test_api_int.name
		url_of_proxy_and_resource = "%s"
	}

	resource "snowflake_external_function" "test_func_2" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		comment = "Terraform acceptance test"
		return_type = "varchar"
		return_behavior = "IMMUTABLE"
		api_integration = snowflake_api_integration.test_api_int.name
		header {
			name = "x-custom-header"
			value = "snowflake"
		}
		max_batch_rows = 500
		url_of_proxy_and_resource = "%s"
	}
	`, name, name, name, prefixes, name, url, name, url+"_2")
}
