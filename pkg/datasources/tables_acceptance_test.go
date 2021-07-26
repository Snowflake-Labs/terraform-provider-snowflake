package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTables(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: tables(databaseName, schemaName, tableName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tables.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_tables.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_tables.t", "tables.#"),
					resource.TestCheckResourceAttr("data.snowflake_tables.t", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_tables.t", "tables.0.name", tableName),
				),
			},
		},
	})
}

func tables(databaseName string, schemaName string, tableName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "d" {
		name = "%v"
	}

	resource snowflake_schema "s"{
		name 	 = "%v"
		database = snowflake_database.d.name
	}

	resource snowflake_table "t"{
		name 	 = "%v"
		database = snowflake_schema.s.database
		schema 	 = snowflake_schema.s.name
		column {
			name = "column2"
			type = "VARCHAR(16)"
		}
	}

	data snowflake_tables "t" {
		database = snowflake_table.t.database
		schema = snowflake_table.t.schema
	}
	`, databaseName, schemaName, tableName)
}
