package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStreams(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	streamName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: streams(databaseName, schemaName, tableName, streamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_streams.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_streams.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_streams.t", "streams.#"),
					resource.TestCheckResourceAttr("data.snowflake_streams.t", "streams.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_streams.t", "streams.0.name", streamName),
				),
			},
		},
	})
}

func streams(databaseName string, schemaName string, tableName string, streamName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "test_database" {
		name = "%v"
	}

	resource snowflake_schema "test_schema" {
		name 	 = "%v"
		database = snowflake_database.test_database.name
	}

	resource snowflake_table "test_stream_on_table" {
		database 	    = snowflake_database.test_database.name
		schema   	    = snowflake_schema.test_schema.name
		change_tracking = true
		name     	    = "%v"
		comment  	    = "Terraform acceptance test"
		column {
			name = "column1"
			type = "VARIANT"
		}
		column {
			name = "column2"
			type = "VARCHAR(16777216)"
		}
	}
	
	resource snowflake_stream "test_stream" {
		database = snowflake_database.test_database.name
		schema   = snowflake_schema.test_schema.name
		name     = "%v"
		comment  = "Terraform acceptance test"
		on_table = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_table.test_stream_on_table.name}"
	}

	data snowflake_streams "t" {
		database = snowflake_stream.test_stream.database
		schema = snowflake_stream.test_stream.schema
		depends_on = [snowflake_stream.test_stream]
	}
	`, databaseName, schemaName, tableName, streamName)
}
