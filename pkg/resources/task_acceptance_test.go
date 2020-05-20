package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTask(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{{
			Config: taskConfig(accName),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("snowflake_task.test", "name", accName),
				resource.TestCheckResourceAttr("snowflake_task.test", "database", accName),
				resource.TestCheckResourceAttr("snowflake_task.test", "schema", accName),
				resource.TestCheckResourceAttr("snowflake_task.test", "warehouse", accName),
				resource.TestCheckResourceAttr("snowflake_task.test", "schedule", "60 minute"),
				resource.TestCheckResourceAttr("snowflake_task.test", "sql_statement", "SELECT 1"),
				resource.TestCheckResourceAttr("snowflake_task.test", "comment", "Terraform acceptance test"),
			),
		}},
	})
}

func taskConfig(n string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_warehouse" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_task" "test" {
	name = "%v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	warehouse = snowflake_warehouse.test.name
	comment = "Terraform acceptance test"
	schedule = "60 minute"
	sql_statement = "SELECT 1"
}
`, n, n, n, n)
}
