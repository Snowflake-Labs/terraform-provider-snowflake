package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExternalFunctions(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	apiName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	externalFunctionName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: externalFunctions(databaseName, schemaName, apiName, externalFunctionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_external_functions.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_external_functions.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_external_functions.t", "external_functions.#"),
					resource.TestCheckResourceAttr("data.snowflake_external_functions.t", "external_functions.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_external_functions.t", "external_functions.0.name", externalFunctionName),
				),
			},
		},
	})
}

func externalFunctions(databaseName string, schemaName string, apiName string, externalFunctionName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "test_database" {
		name = "%v"
	}

	resource snowflake_schema "test_schema"{
		name 	 = "%v"
		database = snowflake_database.test_database.name
	}

	resource "snowflake_api_integration" "test_api_int" {
		name = "%v"
		api_provider = "aws_api_gateway"
		api_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
		api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
		enabled = true
	}

	resource "snowflake_external_function" "test_func" {
		name     = "%v"
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
		url_of_proxy_and_resource = "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"
	}

	data snowflake_external_functions "t" {
		database = snowflake_external_function.test_func.database
		schema = snowflake_external_function.test_func.schema
		depends_on = [snowflake_external_function.test_func]
	}
	`, databaseName, schemaName, apiName, externalFunctionName)
}
