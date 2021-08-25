package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Stream(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: streamConfig(accName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "insert_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				Config: streamConfig(accName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", true),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				Config: externalTableStreamConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_EXTERNAL_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "on_external_table", true),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
		},
	})
}

func streamConfig(name string, append_only bool) string {
	append_only_config := ""
	if append_only {
		append_only_config = "append_only = true"
	}

	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_stream_on_table" {
	database        = snowflake_database.test_database.name
	schema          = snowflake_schema.test_schema.name
	name            = "STREAM_ON_TABLE"
	comment         = "Terraform acceptance test"
	change_tracking = true

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
}

resource "snowflake_stream" "test_stream" {
	database    = snowflake_database.test_database.name
	schema      = snowflake_schema.test_schema.name
	name        = "%s"
	comment     = "Terraform acceptance test"
	on_table    = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_table.test_stream_on_table.name}"
	%s
}
`
	return fmt.Sprintf(s, name, name, name, append_only_config)
}

func externalTableStreamConfig(name string) string {
	// Refer to external_table_acceptance_test.go for the original source on
	// external table resources and dependents (modified slightly here).
	locations := []string{"s3://com.example.bucket/prefix"}
	s := `
resource "snowflake_database" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_stage" "test" {
	name = "%v"
	url = "s3://com.example.bucket/prefix"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	comment = "Terraform acceptance test"
	storage_integration = snowflake_storage_integration.external_table_stream_integration.name
}

resource "snowflake_storage_integration" "external_table_stream_integration" {
	name = "%v"
	storage_allowed_locations = %q
	storage_provider = "S3"
	storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
}

resource "snowflake_external_table" "test_external_stream_table" {
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
	name     = "STREAM_ON_EXTERNAL_TABLE"
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "STRING"
    as = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
    as = "($1:'CreatedDate'::timestamp)"
	}
  file_format = "TYPE = CSV"
  location = "@${snowflake_database.test.name}.${snowflake_schema.test.name}.${snowflake_stage.test.name}"
}

resource "snowflake_stream" "test_external_table_stream" {
	database          = snowflake_database.test_database.name
	schema            = snowflake_schema.test_schema.name
	name              = "%s"
	comment           = "Terraform acceptance test"
	on_external_table = true
	on_table          = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_external_table.test_external_stream_table.name}"
}
`

	return fmt.Sprintf(s, name, name, name, name, locations, name, name)
}
