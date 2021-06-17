package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExternalTable(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_EXTERNAL_TABLE_TESTS"); ok {
		t.Skip("Skipping TestAccExternalTable")
	}
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: externalTableConfig(accName, []string{"s3://com.example.bucket/prefix"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "database", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func externalTableConfig(name string, locations []string) string {
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
	storage_integration = snowflake_storage_integration.i.name
}

resource "snowflake_storage_integration" "i" {
	name = "%v"
	storage_allowed_locations = %q
	storage_provider = "S3"
	storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
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
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
    as = "($1:'CreatedDate'::timestamp)"
	}
  file_format = "TYPE = CSV"
  location = "@${snowflake_database.test.name}.${snowflake_schema.test.name}.${snowflake_stage.test.name}"
}
`
	return fmt.Sprintf(s, name, name, name, name, locations, name)
}
