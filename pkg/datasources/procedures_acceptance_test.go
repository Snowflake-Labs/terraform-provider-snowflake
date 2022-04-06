package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProcedures(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	procedureName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	procedureWithArgumentsName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: procedures(databaseName, schemaName, procedureName, procedureWithArgumentsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_procedures.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_procedures.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_procedures.t", "procedures.#"),
					resource.TestCheckResourceAttr("data.snowflake_procedures.t", "procedures.#", "2"),
				),
			},
		},
	})
}

func procedures(databaseName string, schemaName string, procedureName string, procedureWithArgumentsName string) string {
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

resource "snowflake_procedure" "test_proc_simple" {
	name = "%v"
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	return_type = "VARCHAR"
	statement = "return \"Hi\""
}

resource "snowflake_procedure" "test_proc" {
	name = "%v"
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	arguments {
		name = "arg1"
		type = "varchar"
	}
	comment = "Terraform acceptance test"
	return_type = "varchar"
	statement = "var X=3\nreturn X"
}

data snowflake_procedures "t" {
	database = snowflake_database.test_database.name
	schema = snowflake_schema.test_schema.name
	depends_on = [snowflake_procedure.test_proc_simple, snowflake_procedure.test_proc]
}
`
	return fmt.Sprintf(s, databaseName, schemaName, procedureName, procedureWithArgumentsName)
}
