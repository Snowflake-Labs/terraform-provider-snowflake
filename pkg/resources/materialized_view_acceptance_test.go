package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_MaterializedView(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_MATERIALIZED_VIEW_TESTS"); ok {
		t.Skip("Skipping TestAcc_MaterializedView")
	}
	dbName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	warehouseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: materializedViewConfig(warehouseName, dbName, schemaName, tableName, viewName, fmt.Sprintf("SELECT ID, DATA FROM \\\"%s\\\";", tableName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", dbName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", warehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_materialized_view.test", "is_secure", true),
				),
			},
		},
	})
}

func TestAcc_MaterializedView2(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_MATERIALIZED_VIEW_TESTS"); ok {
		t.Skip("Skipping TestAcc_MaterializedView2")
	}
	dbName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	warehouseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: materializedViewConfig(warehouseName, dbName, schemaName, tableName, viewName, fmt.Sprintf("SELECT ID, DATA FROM \\\"%s\\\" WHERE ID LIKE 'foo%%';", tableName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", dbName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", warehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_materialized_view.test", "is_secure", true),
				),
			},
		},
	})
}

func materializedViewConfig(warehouseName string, dbName string, schemaName string, tableName string, viewName string, q string) string {
	// convert the cluster from string slice to string
	return fmt.Sprintf(`
resource "snowflake_warehouse" "test" {
	name = "%s"
	initially_suspended = false
}

resource "snowflake_database" "test" {
	name = "%s"
}

resource "snowflake_schema" "test" {
	database  = snowflake_database.test.name
	name      = "%s"
}

resource "snowflake_table" "test" {
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
	name     = "%s"

	column {
		name = "ID"
		type = "NUMBER(38,0)"
	}

	column {
		name = "DATA"
		type = "VARCHAR(16777216)"
	}
}

resource "snowflake_materialized_view" "test" {
	name      = "%s"
	comment   = "Terraform test resource"
	database  = snowflake_database.test.name
	schema    = snowflake_schema.test.name
	warehouse = snowflake_warehouse.test.name
	is_secure = true
	or_replace = false
	statement = "%s"

	depends_on = [
		snowflake_warehouse.test,
  		snowflake_table.test
  	]
}
`, warehouseName, dbName, schemaName, tableName, viewName, q)
}
