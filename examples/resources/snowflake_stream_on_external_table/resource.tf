resource "snowflake_external_table" "external_table" {
  database    = "db"
  schema      = "schema"
  name        = "external_table"
  comment     = "External table"
  file_format = "TYPE = CSV FIELD_DELIMITER = '|'"
  location    = "@stage/directory/"

  column {
    name = "id"
    type = "int"
  }

  column {
    name = "data"
    type = "text"
  }
}


resource "snowflake_stream_on_external_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  copy_grants    = true
  external_table = snowflake_external_table.external_table.fully_qualified_name
  insert_only    = "true"

  at {
    statement = "8e5d0ca9-005e-44e6-b858-a8f5b37c5726"
  }

  comment = "A stream."
}
