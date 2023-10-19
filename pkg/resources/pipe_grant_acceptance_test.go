package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_PipeGrant(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: pipeGrantConfig(accName, "OPERATE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "pipe_name", accName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "privilege", "OPERATE"),
				),
			},
			{
				Config: pipeGrantConfig(accName, "ALL PRIVILEGES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "pipe_name", accName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			{
				ResourceName:      "snowflake_pipe_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_PipeGrantWithDefaultPrivilege(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: pipeGrantConfigWithDefaultPrivilege(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "pipe_name", accName),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_pipe_grant.test", "privilege", "OPERATE"),
				),
			},
			{
				ResourceName:      "snowflake_pipe_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func pipeGrantConfig(name, privilege string) string {
	s := `
resource "snowflake_table" "test" {
  database = "terraform_test_database"
  schema   = "terraform_test_schema"
  name     = "%s"
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
  name = "%s"
  database = "terraform_test_database"
  schema = "terraform_test_schema"
  comment = "Terraform acceptance test"
}

resource "snowflake_pipe_grant" "test" {
  pipe_name = snowflake_pipe.test.name
  database_name = "terraform_test_database"
  roles         = [snowflake_role.test.name]
  schema_name   = "terraform_test_schema"
  privilege 	  = "%s"
}

resource "snowflake_pipe" "test" {
  database       = "terraform_test_database"
  schema         = "terraform_test_schema"
  name           = "%s"
  comment        = "Terraform acceptance test"
  copy_statement = <<CMD
COPY INTO "${snowflake_table.test.database}"."${snowflake_table.test.schema}"."${snowflake_table.test.name}"
  FROM @"${snowflake_stage.test.database}"."${snowflake_stage.test.schema}"."${snowflake_stage.test.name}"
  FILE_FORMAT = (TYPE = CSV)
CMD
  auto_ingest    = false
}
`
	return fmt.Sprintf(s, name, name, name, privilege, name)
}

func pipeGrantConfigWithDefaultPrivilege(name string) string {
	s := `
resource "snowflake_table" "test" {
  database = "terraform_test_database"
  schema   = "terraform_test_schema"
  name     = "%s"
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
  name = "%s"
  database = "terraform_test_database"
  schema = "terraform_test_schema"
  comment = "Terraform acceptance test"
}

resource "snowflake_pipe_grant" "test" {
  pipe_name = snowflake_pipe.test.name
  database_name = "terraform_test_database"
  roles         = [snowflake_role.test.name]
  schema_name   = "terraform_test_schema"
}

resource "snowflake_pipe" "test" {
  database       = "terraform_test_database"
  schema         = "terraform_test_schema"
  name           = "%s"
  comment        = "Terraform acceptance test"
  copy_statement = <<CMD
COPY INTO "${snowflake_table.test.database}"."${snowflake_table.test.schema}"."${snowflake_table.test.name}"
  FROM @"${snowflake_stage.test.database}"."${snowflake_stage.test.schema}"."${snowflake_stage.test.name}"
  FILE_FORMAT = (TYPE = CSV)
CMD
  auto_ingest    = false
}
`
	return fmt.Sprintf(s, name, name, name, name)
}
