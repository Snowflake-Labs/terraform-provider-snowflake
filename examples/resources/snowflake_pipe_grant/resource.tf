resource snowflake_pipe_grant grant {
  database_name = "db"
  schema_name   = "schema"
  sequence_name = "sequence"

  privilege = "operate"
  roles = [
    "role1",
    "role2",
  ]

  on_future         = false
  with_grant_option = false
}
