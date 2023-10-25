// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ExternalFunction(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_EXTERNAL_FUNCTION_TESTS"); ok {
		t.Skip("Skipping TestAccExternalFunction")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalFunctionConfig(accName, []string{"https://123456.execute-api.us-west-2.amazonaws.com/prod/"}, "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_function.test_func", "name", accName),
					resource.TestCheckResourceAttr("snowflake_external_function.test_func", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrSet("snowflake_external_function.test_func", "created_on"),
					resource.TestCheckResourceAttr("snowflake_external_function.test_func_2", "request_translator", fmt.Sprintf("%s.%s.TEST_FUNC_REQ_TRANSLATOR", accName, accName)),
					resource.TestCheckResourceAttr("snowflake_external_function.test_func_2", "response_translator", fmt.Sprintf("%s.%s.TEST_FUNC_RES_TRANSLATOR", accName, accName)),
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

	resource "snowflake_function" "test_func_req_translator" {
	  	name     =  upper("test_func_req_translator")
	  	database = snowflake_database.test_database.name
	  	schema   = snowflake_schema.test_schema.name
	  	arguments {
			name = "EVENT"
			type = "OBJECT"
	  	}
	  	comment             = "Terraform acceptance test"
	  	return_type         = "OBJECT"
	  	language            = "javascript"
	  	statement           = <<EOH
		  	let exeprimentName = EVENT.body.data[0][1]
		  	return { "body": { "name": test }}
	  	EOH
	}


		resource "snowflake_function" "test_func_res_translator" {
		  name     =  upper("test_func_res_translator")
		  database = snowflake_database.test_database.name
          schema   = snowflake_schema.test_schema.name
		  arguments {
			name = "EVENT"
			type = "OBJECT"
		  }
		  comment             = "Terraform acceptance test"
		  return_type         = "OBJECT"
		  language            = "javascript"
		  statement           = <<EOH
			  return { "body": { "data" :  [[0, EVENT]] } };
		  EOH
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
		return_type = "variant"
		return_behavior = "IMMUTABLE"
		api_integration = snowflake_api_integration.test_api_int.name
		url_of_proxy_and_resource = "%s"
	}

	resource "snowflake_external_function" "test_func_2" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		comment = "Terraform acceptance test"
		return_type = "variant"
		return_behavior = "IMMUTABLE"
		api_integration = snowflake_api_integration.test_api_int.name
		header {
			name = "x-custom-header"
			value = "snowflake"
		}
		max_batch_rows = 500
		request_translator = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_function.test_func_req_translator.name}"
		response_translator = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_function.test_func_res_translator.name}"
		url_of_proxy_and_resource = "%s"
	}
	`, name, name, name, prefixes, name, url, name, url+"_2")
}
