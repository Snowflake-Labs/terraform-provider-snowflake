resource snowflake_stream stream {
  comment = "A stream."

  database = "db"
  schema   = "schema"
  name     = "stream"

  on_table    = "table"
  append_only = false

  owner = "role1"
}
