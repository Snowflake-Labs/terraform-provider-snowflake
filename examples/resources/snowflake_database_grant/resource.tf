resource snowflake_database_grant grant {
  database_name = "db"

  privilege = "USAGE"
  roles     = ["role1", "role2"]
  shares    = ["share1", "share2"]

  with_grant_option = false
}
