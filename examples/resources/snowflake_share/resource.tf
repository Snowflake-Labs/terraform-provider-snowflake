resource snowflake_share share {
  database_name = "db"
  schema_name   = "schema"
  stage_name    = "stage"

  privilege = "USAGE"
  roles     = ["role1", "role2"]
  shares    = ["share1", "share2"]

  with_grant_option = false
}
