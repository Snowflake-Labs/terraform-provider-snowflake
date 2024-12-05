# basic resource
resource "snowflake_stream_on_external_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  external_table = snowflake_external_table.example.fully_qualified_name
}


# resource with additional fields
resource "snowflake_stream_on_external_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  copy_grants    = true
  external_table = snowflake_external_table.example.fully_qualified_name
  insert_only    = "true"

  at {
    statement = "8e5d0ca9-005e-44e6-b858-a8f5b37c5726"
  }

  comment = "A stream."
}
