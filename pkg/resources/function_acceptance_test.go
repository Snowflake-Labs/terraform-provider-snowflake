package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Function(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_FUNCTION_TESTS"); ok {
		t.Skip("Skipping TestAcc_Function")
	}

	dbName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	functName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	expBody1 := "3.141592654::FLOAT"
	expBody2 := "var X=3\nreturn X"
	expBody3 := "select 1, 2\nunion all\nselect 3, 4\n"
	expBody4 := `class CoolFunc {public static String test(int n) {return "hello!";}}`

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: functionConfig(dbName, schemaName, functName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_function.test_funct", "name", functName),
					resource.TestCheckResourceAttr("snowflake_function.test_funct", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct", "statement", expBody2),
					resource.TestCheckResourceAttr("snowflake_function.test_funct", "arguments.#", "1"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct", "arguments.0.name", "ARG1"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct", "arguments.0.type", "VARCHAR"),

					resource.TestCheckResourceAttr("snowflake_function.test_funct_simple", "name", functName),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_simple", "comment", "user-defined function"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_simple", "statement", expBody1),

					resource.TestCheckResourceAttr("snowflake_function.test_funct_complex", "name", functName),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_complex", "comment", "Table func with 2 args"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_complex", "statement", expBody3),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_complex", "arguments.#", "2"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_complex", "arguments.1.name", "ARG2"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_complex", "arguments.1.type", "DATE"),

					resource.TestCheckResourceAttr("snowflake_function.test_funct_java", "name", functName),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_java", "comment", "Terraform acceptance test for java"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_java", "statement", expBody4),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_java", "arguments.#", "1"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_java", "arguments.0.name", "ARG1"),
					resource.TestCheckResourceAttr("snowflake_function.test_funct_java", "arguments.0.type", "NUMBER"),
					checkBool("snowflake_function.test_funct_java", "is_secure", false), // this is from user_acceptance_test.go

					// TODO: temporarily remove unit tests to allow for urgent release
					// resource.TestCheckResourceAttr("snowflake_function.test_funct_python", "name", functName),
					// resource.TestCheckResourceAttr("snowflake_function.test_funct_python", "comment", "Terraform acceptance test for python"),
					// resource.TestCheckResourceAttr("snowflake_function.test_funct_python", "statement", expBody5),
					// resource.TestCheckResourceAttr("snowflake_function.test_funct_python", "arguments.#", "2"),
					// resource.TestCheckResourceAttr("snowflake_function.test_funct_python", "arguments.0.name", "ARG1"),
					// resource.TestCheckResourceAttr("snowflake_function.test_funct_python", "arguments.0.type", "NUMBER"),
				),
			},
		},
	})
}

func functionConfig(db, schema, name string) string {
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
	
	resource "snowflake_function" "test_funct_simple" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		return_type = "float"
		statement = "3.141592654::FLOAT"
	}

	resource "snowflake_function" "test_funct" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		arguments {
			name = "arg1"
			type = "varchar"
		}
		comment = "Terraform acceptance test"
		return_type = "varchar"
		language = "javascript"
		statement = "var X=3\nreturn X"
	}

	resource "snowflake_function" "test_funct_java" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		arguments {
			name = "arg1"
			type = "number"
		}
		comment = "Terraform acceptance test for java"
		return_type = "varchar"
		language = "java"
		handler = "CoolFunc.test"
		statement = "class CoolFunc {public static String test(int n) {return \"hello!\";}}"
	}

	resource "snowflake_function" "test_funct_complex" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		arguments {
			name = "arg1"
			type = "varchar"
		}
		arguments {
			name = "arg2"
			type = "DATE"
		}		
		comment = "Table func with 2 args"
		return_type = "table (x number, y number)"
		statement = <<EOT
select 1, 2
union all
select 3, 4
EOT
	}
	`, db, schema, name, name, name, name)
}
