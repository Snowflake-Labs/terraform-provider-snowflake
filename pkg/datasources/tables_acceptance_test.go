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
	stageName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	externalTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: tables(databaseName, schemaName, tableName, stageName, externalTableName),
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

func tables(databaseName string, schemaName string, tableName string, stageName string, externalTableName string) string {
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

	resource "snowflake_stage" "s" {
		name = "%v"
		url = "s3://snowflake-workshop-lab/weather-nyc"
		database = snowflake_database.d.name
		schema = snowflake_schema.s.name
	}

	resource "snowflake_external_table" "et" {
		database = snowflake_database.d.name
		schema   = snowflake_schema.s.name
		name     = "%v"
		column {
			name = "column1"
			type = "STRING"
		as = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
		}
	    file_format = "TYPE = CSV"
	    location = "@${snowflake_database.d.name}.${snowflake_schema.s.name}.${snowflake_stage.s.name}"
	}

	data snowflake_tables "t" {
		database = snowflake_table.t.database
		schema = snowflake_table.t.schema
		depends_on = [snowflake_table.t, snowflake_external_table.et]
	}
	`, databaseName, schemaName, tableName, stageName, externalTableName)
}
