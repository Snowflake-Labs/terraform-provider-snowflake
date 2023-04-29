package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_TaskGrant(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskGrantConfig(name, 8, normal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "task_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "statement_timeout_in_seconds", "86400"),
				),
			},
			// UPDATE MAX_CONCURRENCY_LEVEL
			{
				Config: taskGrantConfig(name, 10, normal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "task_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "max_concurrency_level", "10"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "statement_timeout_in_seconds", "86400"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_task_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_TaskGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskGrantConfig(name, 8, onAll),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_task_grant.test", "task_name"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "statement_timeout_in_seconds", "86400"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_task_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_TaskGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskGrantConfig(name, 8, onFuture),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_task_grant.test", "task_name"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "statement_timeout_in_seconds", "86400"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_task_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func taskGrantConfig(name string, concurrency int32, grantType grantType) string {
	var taskNameConfig string
	switch grantType {
	case normal:
		taskNameConfig = "task_name \t= snowflake_task.test.name"
	case onFuture:
		taskNameConfig = "on_future = true"
	case onAll:
		taskNameConfig = "on_all = true"
	}

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
  name                         = snowflake_database.test.name
  max_concurrency_level        = %d
  statement_timeout_in_seconds = 86400
  query_acceleration_max_scale_factor = 0
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
  %s
  database_name = snowflake_database.test.name
  roles         = [snowflake_role.test.name]
  schema_name   = snowflake_schema.test.name
  privilege 	= "OPERATE"
}
`
	return fmt.Sprintf(s, name, name, concurrency, taskNameConfig)
}
