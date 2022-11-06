resource "snowflake_stream" "stream" {
  comment = "A stream."

  database = "database"
  schema   = "schema"
  name     = "stream"

  on_table    = "table"
  append_only = false
  insert_only = false

  owner = "role1"
}
