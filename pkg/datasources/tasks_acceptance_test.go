package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTasks(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	taskName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: tasks(databaseName, schemaName, taskName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tasks.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_tasks.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_tasks.t", "tasks.#"),
					resource.TestCheckResourceAttr("data.snowflake_tasks.t", "tasks.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_tasks.t", "tasks.0.name", taskName),
				),
			},
		},
	})
}

func tasks(databaseName string, schemaName string, taskName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "test" {
	   name = "%v"
	}

	resource snowflake_schema "test"{
		name 	 = "%v"
		database = snowflake_database.test.name
	}

	resource snowflake_warehouse "test" {
		name = snowflake_database.test.name
	}

	resource snowflake_task "test" {
		name     	  = "%v"
		database  	  = snowflake_database.test.name
		schema   	  = snowflake_schema.test.name
		warehouse 	  = snowflake_warehouse.test.name
		sql_statement = "SHOW FUNCTIONS"
		enabled  	  = true
		schedule 	  = "15 MINUTES"
		lifecycle {
		  ignore_changes = [session_parameters]
		}
	  }

	data snowflake_tasks "t" {
		database = snowflake_task.test.database
		schema = snowflake_task.test.schema
		depends_on = [snowflake_task.test]
	}
	`, databaseName, schemaName, taskName)
}
