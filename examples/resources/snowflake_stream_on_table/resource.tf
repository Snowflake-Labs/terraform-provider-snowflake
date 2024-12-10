# basic resource
resource "snowflake_stream_on_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  table = snowflake_table.example.fully_qualified_name
}

# resource with more fields set
resource "snowflake_stream_on_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  copy_grants       = true
  table             = snowflake_table.example.fully_qualified_name
  append_only       = "true"
  show_initial_rows = "true"

  at {
    statement = "8e5d0ca9-005e-44e6-b858-a8f5b37c5726"
  }

  comment = "A stream."
}
