package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_MaterializedView(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	warehouseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: materializedViewConfig(warehouseName, tableName, viewName, fmt.Sprintf("SELECT ID, DATA FROM \\\"%s\\\";", tableName), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", warehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_materialized_view.test", "is_secure", true),
				),
			},
		},
	})
}

func TestAcc_MaterializedView2(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	warehouseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: materializedViewConfig(warehouseName, tableName, viewName, fmt.Sprintf("SELECT ID, DATA FROM \\\"%s\\\" WHERE ID LIKE 'foo%%';", tableName), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", warehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_materialized_view.test", "is_secure", true),
				),
			},
		},
	})
}

func materializedViewConfig(warehouseName string, tableName string, viewName string, q string, databaseName string, schemaName string) string {
	// convert the cluster from string slice to string
	return fmt.Sprintf(`
resource "snowflake_warehouse" "wh" {
	name = "%s"
}
resource "snowflake_table" "test" {
	name     = "%s"
	database = "%s"
	schema   = "%s"

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
	database  = "%s"
	schema    = "%s"
	warehouse = snowflake_warehouse.wh.name
	is_secure = true
	or_replace = false
	statement = "%s"

	depends_on = [
  		snowflake_table.test
  	]
}
`, warehouseName, tableName, databaseName, schemaName, viewName, databaseName, schemaName, q)
}
