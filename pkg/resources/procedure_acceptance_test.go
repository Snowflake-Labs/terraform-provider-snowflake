package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Procedure(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_PROCEDURE_TESTS"); ok {
		t.Skip("Skipping TestAcc_Procedure")
	}

	dbName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	procName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	procBody := "var X=3\\nreturn X"
	expBody := "var X=3\nreturn X"
	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: procedureConfig(dbName, schemaName, procName, procBody),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "name", procName),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "statement", expBody),
					// resource.TestCheckResourceAttrSet("snowflake_external_function.test_func", "created_on"),
				),
			},
		},
	})
}

func procedureConfig(db, schema, name, stmnt string) string {
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


	resource "snowflake_procedure" "test_proc" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		arguments {
			name = "arg1"
			type = "varchar"
		}
		comment = "Terraform acceptance test"
		return_type = "varchar"
		statement = "%s"
	}
	`, db, schema, name, stmnt)
}
