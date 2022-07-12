resource snowflake_pipe_grant grant {
  database_name = "db"
  schema_name   = "schema"
  pipe_name     = "pipe"

  privilege = "operate"
  roles = [
    "role1",
    "role2",
  ]

  on_future         = false
  with_grant_option = false
}
