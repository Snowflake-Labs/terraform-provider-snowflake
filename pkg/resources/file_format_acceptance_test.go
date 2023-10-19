package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_FileFormatCSV(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigCSV(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "CSV"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "record_delimiter", "\r"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "field_delimiter", ";"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "file_extension", ".ssv"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_header", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_blank_lines", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "date_format", "YYY-MM-DD"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "time_format", "HH24:MI"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "timestamp_format", "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "binary_format", "UTF8"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "escape", "\\"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "escape_unenclosed_field", "!"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "field_optionally_enclosed_by", "'"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "2"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.1", ""),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "error_on_column_count_mismatch", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "empty_field_as_null", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_byte_order_mark", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "encoding", "UTF-16"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_FileFormatJSON(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigJSON(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "JSON"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "date_format", "YYY-MM-DD"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "time_format", "HH24:MI"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "timestamp_format", "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "binary_format", "UTF8"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "file_extension", ".jsn"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "enable_octal", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "allow_duplicate", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "strip_outer_array", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "strip_null_values", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "ignore_utf8_errors", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_byte_order_mark", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_FileFormatAvro(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigAvro(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "AVRO"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_FileFormatORC(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigORC(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "ORC"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_FileFormatParquet(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigParquet(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "PARQUET"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "SNAPPY"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "binary_as_text", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "trim_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.#", "1"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "null_if.0", "NULL"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAcc_FileFormatXML(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigXML(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "XML"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "compression", "GZIP"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "preserve_space", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "strip_outer_element", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "disable_snowflake_data", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "disable_auto_convert", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "skip_byte_order_mark", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

// The following tests check that Terraform will accept the default values generated at creation and not drift.
// See https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1706
func TestAcc_FileFormatCSVDefaults(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(accName, "CSV"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "CSV"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatJSONDefaults(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(accName, "JSON"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "JSON"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatAVRODefaults(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(accName, "AVRO"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "AVRO"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatORCDefaults(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(accName, "ORC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "ORC"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatPARQUETDefaults(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(accName, "PARQUET"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "PARQUET"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_FileFormatXMLDefaults(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfigFullDefaults(accName, "XML"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "XML"),
				),
			},
			{
				ResourceName:      "snowflake_file_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func fileFormatConfigCSV(n string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "terraform_test_database"
	schema = "terraform_test_schema"
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
	null_if = ["NULL", ""]
	error_on_column_count_mismatch = true
	replace_invalid_characters = true
	empty_field_as_null = false
	skip_byte_order_mark = false
	encoding = "UTF-16"
	comment = "Terraform acceptance test"
}
`, n)
}

func fileFormatConfigJSON(n string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "terraform_test_database"
	schema = "terraform_test_schema"
	format_type = "JSON"
	compression = "GZIP"
	date_format = "YYY-MM-DD"
	time_format = "HH24:MI"
	timestamp_format = "YYYY-MM-DD HH24:MI:SS.FFTZH:TZM"
	binary_format = "UTF8"
	trim_space = true
	null_if = ["NULL"]
	file_extension = ".jsn"
	enable_octal = true
	allow_duplicate = true
	strip_outer_array = true
	strip_null_values = true
	ignore_utf8_errors = true
	skip_byte_order_mark = false
	comment = "Terraform acceptance test"
}
`, n)
}

func fileFormatConfigAvro(n string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "terraform_test_database"
	schema = "terraform_test_schema"
	format_type = "AVRO"
	compression = "GZIP"
	trim_space = true
	null_if = ["NULL"]
	comment = "Terraform acceptance test"
}
`, n)
}

func fileFormatConfigORC(n string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "terraform_test_database"
	schema = "terraform_test_schema"
	format_type = "ORC"
	trim_space = true
	null_if = ["NULL"]
	comment = "Terraform acceptance test"
}
`, n)
}

func fileFormatConfigParquet(n string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "terraform_test_database"
	schema = "terraform_test_schema"
	format_type = "PARQUET"
	compression = "SNAPPY"
	binary_as_text = true
	trim_space = true
	null_if = ["NULL"]
	comment = "Terraform acceptance test"
}
`, n)
}

func fileFormatConfigXML(n string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "terraform_test_database"
	schema = "terraform_test_schema"
	format_type = "XML"
	compression = "GZIP"
	ignore_utf8_errors = true
	preserve_space = true
	strip_outer_element = true
	disable_snowflake_data =  true
	disable_auto_convert =  true
	skip_byte_order_mark = false
	comment = "Terraform acceptance test"
}
`, n)
}

func fileFormatConfigFullDefaults(n, formatType string) string {
	return fmt.Sprintf(`
resource "snowflake_file_format" "test" {
	name = "%v"
	database = "terraform_test_database"
	schema = "terraform_test_schema"
	format_type = "%s"
}
`, n, formatType)
}
