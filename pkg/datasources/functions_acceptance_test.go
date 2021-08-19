package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFunctions(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	functionName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	functionWithArgumentsName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: functions(databaseName, schemaName, functionName, functionWithArgumentsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_functions.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_functions.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_functions.t", "functions.#"),
					resource.TestCheckResourceAttr("data.snowflake_functions.t", "functions.#", "2"),
				),
			},
		},
	})
}

func functions(databaseName string, schemaName string, functionName string, functionWithArgumentsName string) string {
	s := `
resource "snowflake_database" "test_database" {
	name 	  = "%v"
	comment = "Terraform acceptance test"
}
resource "snowflake_schema" "test_schema" {
	name 	   = "%v"
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

data snowflake_functions "t" {
	database = snowflake_database.test_database.name
	schema = snowflake_schema.test_schema.name
	depends_on = [snowflake_function.test_funct_simple, snowflake_function.test_funct]
}
`
	return fmt.Sprintf(s, databaseName, schemaName, functionName, functionWithArgumentsName)
}
