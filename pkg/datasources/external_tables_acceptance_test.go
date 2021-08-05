package datasources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExternalTables(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_EXTERNAL_TABLE_TESTS"); ok {
		t.Skip("Skipping TestAccExternalTable")
	}

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	stageName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	externalTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: externalTables(databaseName, schemaName, stageName, externalTableName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_external_tables.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_external_tables.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_external_tables.t", "external_tables.#"),
					resource.TestCheckResourceAttr("data.snowflake_external_tables.t", "external_tables.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_external_tables.t", "external_tables.0.name", externalTableName),
				),
			},
		},
	})
}

func externalTables(databaseName string, schemaName string, stageName string, externalTableName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "test" {
		name = "%v"
	}

	resource snowflake_schema "test"{
		name 	 = "%v"
		database = snowflake_database.test.name
	}

	resource "snowflake_stage" "test" {
		name = "%v"
		url = "s3://snowflake-workshop-lab/weather-nyc"
		database = snowflake_database.test.name
		schema = snowflake_schema.test.name
		comment = "Terraform acceptance test"
	}

	resource "snowflake_external_table" "test_table" {
		database = snowflake_database.test.name
		schema   = snowflake_schema.test.name
		name     = "%v"
		comment  = "Terraform acceptance test"
		column {
			name = "column1"
			type = "STRING"
		as = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
		}
	    file_format = "TYPE = CSV"
	    location = "@${snowflake_database.test.name}.${snowflake_schema.test.name}.${snowflake_stage.test.name}"
	}

	data snowflake_external_tables "t" {
		database = snowflake_external_table.test_table.database
		schema = snowflake_external_table.test_table.schema
		depends_on = [snowflake_external_table.test_table]
	}
	`, databaseName, schemaName, stageName, externalTableName)
}
