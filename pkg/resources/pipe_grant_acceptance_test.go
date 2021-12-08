package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_PipeGrant(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: pipeGrantConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "pipe_name", accName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "privilege", "OPERATE"),
				),
			},
		},
	})
}

func pipeGrantConfig(name string) string {
	s := `
resource "snowflake_database" "test" {
  name = "%v"
  comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
  name = snowflake_database.test.name
  database = snowflake_database.test.name
  comment = "Terraform acceptance test"
}

resource "snowflake_table" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = snowflake_schema.test.name
  column {
	name = "id"
	type = "NUMBER(5,0)"
  }
  column {
    name = "data"
	type = "VARCHAR(16)"
  }
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_stage" "test" {
  name = snowflake_schema.test.name
  database = snowflake_database.test.name
  schema = snowflake_schema.test.name
  comment = "Terraform acceptance test"
}

resource "snowflake_pipe_grant" "test" {
  pipe_name = snowflake_pipe.test.name
  database_name = snowflake_database.test.name
  roles         = [snowflake_role.test.name]
  schema_name   = snowflake_schema.test.name
  privilege 	  = "OPERATE"
}

resource "snowflake_pipe" "test" {
  database       = snowflake_database.test.name
  schema         = snowflake_schema.test.name
  name           = snowflake_schema.test.name
  comment        = "Terraform acceptance test"
  copy_statement = <<CMD
COPY INTO "${snowflake_table.test.database}"."${snowflake_table.test.schema}"."${snowflake_table.test.name}"
  FROM @"${snowflake_stage.test.database}"."${snowflake_stage.test.schema}"."${snowflake_stage.test.name}"
  FILE_FORMAT = (TYPE = CSV)
CMD
  auto_ingest    = false
}
`
	return fmt.Sprintf(s, name, name)
}
