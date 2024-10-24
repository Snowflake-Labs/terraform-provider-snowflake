resource "snowflake_view" "view" {
  database  = "database"
  schema    = "schema"
  name      = "view"
  statement = <<-SQL
    select * from foo;
SQL
}

# basic resource
resource "snowflake_stream_on_view" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  view = snowflake_view.view.fully_qualified_name
}

# resource with additional fields
resource "snowflake_stream_on_view" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  copy_grants       = true
  view              = snowflake_view.view.fully_qualified_name
  append_only       = "true"
  show_initial_rows = "true"

  at {
    statement = "8e5d0ca9-005e-44e6-b858-a8f5b37c5726"
  }

  comment = "A stream."
}
