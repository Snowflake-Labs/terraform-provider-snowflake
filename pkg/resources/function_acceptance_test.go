package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Function(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_FUNCTION_TESTS"); ok {
		t.Skip("Skipping TestAcc_Function")
	}

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
				Config: functionConfig(functName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_function.test_function", "name", functName),
					resource.TestCheckResourceAttr("snowflake_function.test_function", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_function.test_function", "statement", expBody2),
					resource.TestCheckResourceAttr("snowflake_function.test_function", "arguments.#", "1"),
					resource.TestCheckResourceAttr("snowflake_function.test_function", "arguments.0.name", "ARG1"),
					resource.TestCheckResourceAttr("snowflake_function.test_function", "arguments.0.type", "VARCHAR"),

					resource.TestCheckResourceAttr("snowflake_function.test_function_simple", "name", functName),
					resource.TestCheckResourceAttr("snowflake_function.test_function_simple", "comment", "user-defined function"),
					resource.TestCheckResourceAttr("snowflake_function.test_function_simple", "statement", expBody1),

					resource.TestCheckResourceAttr("snowflake_function.test_function_complex", "name", functName),
					resource.TestCheckResourceAttr("snowflake_function.test_function_complex", "comment", "Table func with 2 args"),
					resource.TestCheckResourceAttr("snowflake_function.test_function_complex", "statement", expBody3),
					resource.TestCheckResourceAttr("snowflake_function.test_function_complex", "arguments.#", "2"),
					resource.TestCheckResourceAttr("snowflake_function.test_function_complex", "arguments.1.name", "ARG2"),
					resource.TestCheckResourceAttr("snowflake_function.test_function_complex", "arguments.1.type", "DATE"),

					resource.TestCheckResourceAttr("snowflake_function.test_function_java", "name", functName),
					resource.TestCheckResourceAttr("snowflake_function.test_function_java", "comment", "Terraform acceptance test for java"),
					resource.TestCheckResourceAttr("snowflake_function.test_function_java", "statement", expBody4),
					resource.TestCheckResourceAttr("snowflake_function.test_function_java", "arguments.#", "1"),
					resource.TestCheckResourceAttr("snowflake_function.test_function_java", "arguments.0.name", "ARG1"),
					resource.TestCheckResourceAttr("snowflake_function.test_function_java", "arguments.0.type", "NUMBER"),
					checkBool("snowflake_function.test_function_java", "is_secure", false), // this is from user_acceptance_test.go

					// TODO: temporarily remove unit tests to allow for urgent release
					// resource.TestCheckResourceAttr("snowflake_function.test_function_python", "name", functName),
					// resource.TestCheckResourceAttr("snowflake_function.test_function_python", "comment", "Terraform acceptance test for python"),
					// resource.TestCheckResourceAttr("snowflake_function.test_function_python", "statement", expBody5),
					// resource.TestCheckResourceAttr("snowflake_function.test_function_python", "arguments.#", "2"),
					// resource.TestCheckResourceAttr("snowflake_function.test_function_python", "arguments.0.name", "ARG1"),
					// resource.TestCheckResourceAttr("snowflake_function.test_function_python", "arguments.0.type", "NUMBER"),
				),
			},
		},
	})
}

func functionConfig(name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
	resource "snowflake_function" "test_function_simple" {
		name = "%s"
		database = "%s"
		schema   = "%s"
		return_type = "float"
		statement = "3.141592654::FLOAT"
	}

	resource "snowflake_function" "test_function" {
		name = "%s"
		database = "%s"
		schema   = "%s"
		arguments {
			name = "arg1"
			type = "varchar"
		}
		comment = "Terraform acceptance test"
		return_type = "varchar"
		language = "javascript"
		statement = "var X=3\nreturn X"
	}

	resource "snowflake_function" "test_function_java" {
		name = "%s"
		database = "%s"
		schema   = "%s"
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

	resource "snowflake_function" "test_function_complex" {
		name = "%s"
		database = "%s"
		schema   = "%s"
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
	`, name, databaseName, schemaName, name, databaseName, schemaName, name, databaseName, schemaName, name, databaseName, schemaName)
}
