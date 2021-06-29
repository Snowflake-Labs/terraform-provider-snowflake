resource "snowflake_file_format" "example_file_format" {
  name        = "EXAMPLE_FILE_FORMAT"
  database    = "EXAMPLE_DB"
  schema      = "EXAMPLE_SCHEMA"
  format_type = "CSV"
}
