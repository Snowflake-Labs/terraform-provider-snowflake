resource snowflake_view_grant grant {
  database_name = "db"
  schema_name   = "schema"
  view_name     = "view"

  privilege = "select"
  roles = [
    "role1",
    "role2",
  ]

  shares = [
    "share1",
    "share2",
  ]

  on_future         = false
  with_grant_option = false
}
