package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalTables(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	stageName := acc.TestClient().Ids.Alpha()
	externalTableName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
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
