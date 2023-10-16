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

func TestAcc_Procedure(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_PROCEDURE_TESTS"); ok {
		t.Skip("Skipping TestAcc_Procedure")
	}

	dbName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	procName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	expBody1 := "return \"Hi\"\n"
	expBody2 := "var X=3\nreturn X\n"
	expBody3 := "var X=1\nreturn X\n"

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: procedureConfig(dbName, schemaName, procName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "name", procName),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "statement", expBody2),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "arguments.#", "1"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "arguments.0.name", "ARG1"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "arguments.0.type", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc", "execute_as", "OWNER"),

					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_simple", "name", procName),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_simple", "comment", "user-defined procedure"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_simple", "statement", expBody1),

					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "name", procName),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "comment", "Proc with 2 args"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "statement", expBody3),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "execute_as", "CALLER"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "arguments.#", "2"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "arguments.1.name", "ARG2"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "arguments.1.type", "DATE"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "return_behavior", "IMMUTABLE"),
					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_complex", "null_input_behavior", "RETURNS NULL ON NULL INPUT"),

					resource.TestCheckResourceAttr("snowflake_procedure.test_proc_sql", "name", procName+"_sql"),
				),
			},
			{
				ResourceName:      "snowflake_procedure.test_proc_complex",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func procedureConfig(db, schema, name string) string {
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

	resource "snowflake_procedure" "test_proc_simple" {
		name = "%s"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		return_type = "varchar"
		language = "javascript"
		statement = <<-EOF
			return "Hi"
		EOF
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
		language = "javascript"
		return_type = "varchar"
		statement = <<-EOF
			var X=3
			return X
		EOF
	}

	resource "snowflake_procedure" "test_proc_complex" {
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
		comment = "Proc with 2 args"
		return_type = "VARCHAR"
		execute_as = "CALLER"
		return_behavior = "IMMUTABLE"
		null_input_behavior = "RETURNS NULL ON NULL INPUT"
		language = "javascript"
		statement = <<-EOF
			var X=1
			return X
		EOF
	}

	resource "snowflake_procedure" "test_proc_sql" {
		name = "%s_sql"
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		language = "SQL"
		return_type         = "INTEGER"
		execute_as          = "CALLER"
		return_behavior     = "IMMUTABLE"
		null_input_behavior = "RETURNS NULL ON NULL INPUT"
		statement           = <<EOT
	  declare
		x integer;
		y integer;
	  begin
		x := 3;
		y := x * x;
		return y;
	  end;
	  EOT
	  }
	`, db, schema, name, name, name, name)
}
