resource snowflake_sequence_grant grant {
  database_name = "db"
  schema_name   = "schema"
  sequence_name = "sequence"

  privilege = "select"
  roles = [
    "role1",
    "role2",
  ]

  on_future         = false
  with_grant_option = false
}
