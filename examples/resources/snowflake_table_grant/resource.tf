resource snowflake_table_grant grant {
  database_name = "database"
  schema_name   = "schema"
  table_name    = "table"

  privilege = "SELECT"
  roles     = ["role1"]
  shares    = ["share1"]

  on_future         = false
  with_grant_option = false
}
