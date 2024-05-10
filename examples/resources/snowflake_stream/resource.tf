resource "snowflake_table" "table" {
  database = "database"
  schema   = "schema"
  name     = "name"

  column {
    type = "NUMBER(38,0)"
    name = "id"
  }
}

resource "snowflake_stream" "stream" {
  comment = "A stream."

  database = "database"
  schema   = "schema"
  name     = "stream"

  on_table    = snowflake_table.table.qualified_name
  append_only = false
  insert_only = false

  owner = "role1"
}
