resource "snowflake_external_table" "external_table" {
  database    = "db"
  schema      = "schema"
  name        = "external_table"
  comment     = "External table"
  file_format = "TYPE = CSV FIELD_DELIMITER = '|'"

  column {
    name = "id"
    type = "int"
  }

  column {
    name = "data"
    type = "text"
  }
}
