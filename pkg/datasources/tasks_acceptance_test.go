package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tasks(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	taskName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
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
		name                         = snowflake_database.test.name
		max_concurrency_level        = 8
		statement_timeout_in_seconds = 172800
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
