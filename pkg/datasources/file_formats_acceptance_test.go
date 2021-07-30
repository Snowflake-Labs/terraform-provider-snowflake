package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFileFormats(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	fileFormatName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: fileFormats(databaseName, schemaName, fileFormatName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_file_formats.t", "file_formats.#"),
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "file_formats.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_file_formats.t", "file_formats.0.name", fileFormatName),
				),
			},
		},
	})
}

func fileFormats(databaseName string, schemaName string, fileFormatName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "d" {
		name = "%v"
	}

	resource snowflake_schema "s"{
		name 	 = "%v"
		database = snowflake_database.d.name
	}

	resource snowflake_file_format "t"{
		name 	 	= "%v"
		database 	= snowflake_schema.s.database
		schema 	 	= snowflake_schema.s.name
		format_type = "CSV"
		compression = "GZIP"
		record_delimiter = "\r"
		field_delimiter = ";"
		file_extension = ".ssv"
		skip_header = 1
		skip_blank_lines = true
		date_format = "YYY-MM-DD"
		time_format = "HH24:MI"
		timestamp_format = "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"
		binary_format = "UTF8"
		escape = "\\"
		escape_unenclosed_field = "!"
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

	data snowflake_file_formats "t" {
		database = snowflake_file_format.t.database
		schema = snowflake_file_format.t.schema
		depends_on = [snowflake_file_format.t]
	}
	`, databaseName, schemaName, fileFormatName)
}
