package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_FileFormatCSV(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigCSV(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "CSV"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "record_delimiter", "r"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "field_delimiter", "f"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "file_extension", ".fsv"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_header", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_blank_lines", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "date_format", "YYY-MM-DD"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "time_format", "HH24:MI"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "timestamp_format", "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "binary_format", "UTF8"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "escape", "e"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "escape_unenclosed_field", "b"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "field_optionally_enclosed_by", "'"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if", "[\"NULL\"]"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "error_on_column_count_mismatch", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "validate_utf8", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "empty_field_as_null", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_byte_order_mark", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "encoding", "UTF-16"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func fileFormatConfigCSV(n string) string {
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

resource "snowflake_file_format" "test" {
	name = "%v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	format_type = "CSV"
	compression = "GZIP"
	record_delimiter = "r"
	field_delimiter = "f"
	file_extension = ".fsv"
	skip_header = 1
	skip_blank_lines = true
	date_format = "YYY-MM-DD"
	time_format = "HH24:MI"
	timestamp_format = "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"
	binary_format = "UTF8"
	escape = "e"
	escape_unenclosed_field = "b"
	trim_space = true
	field_optionally_enclosed_by = "'"
	null_if = ["NULL"]
	error_on_column_count_mismatch = true
	replace_invalid_characters = true
	validate_utf8 = false
	empty_field_as_null = false 
	skip_byte_order_mark = false
	encoding = "UTF-16"
	comment = "Terraform acceptance test"
}
`, n, n, n)
}
