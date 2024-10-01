resource "snowflake_table" "table" {
  database = "database"
  schema   = "schema"
  name     = "name"

  column {
    type = "NUMBER(38,0)"
    name = "id"
  }
}


# resource with more fields set
resource "snowflake_stream_on_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  copy_grants       = true
  table             = snowflake_table.table.fully_qualified_name
  append_only       = "true"
  show_initial_rows = "true"

  at {
    statement = "8e5d0ca9-005e-44e6-b858-a8f5b37c5726"
  }

  comment = "A stream."
}
