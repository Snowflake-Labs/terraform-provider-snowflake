package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_TaskGrant(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: taskGrantConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "task_name", accName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
				),
			},
		},
	})
}

func taskGrantConfig(name string) string {
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

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_warehouse" "test" {
  name = snowflake_database.test.name
}

resource "snowflake_task" "test" {
  name     	    = snowflake_schema.test.name
  database  	= snowflake_database.test.name
  schema   	  	= snowflake_schema.test.name
  warehouse 	= snowflake_warehouse.test.name
  sql_statement = "SHOW FUNCTIONS"
  enabled  	  	= true
  schedule 	  	= "15 MINUTES"
  lifecycle {
    ignore_changes = [session_parameters]
  }
}

resource "snowflake_task_grant" "test" {
  task_name 	= snowflake_task.test.name
  database_name = snowflake_database.test.name
  roles         = [snowflake_role.test.name]
  schema_name   = snowflake_schema.test.name
  privilege 	= "OPERATE"
}
`
	return fmt.Sprintf(s, name, name)
}
