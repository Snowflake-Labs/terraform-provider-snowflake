resource "snowflake_stage_grant" "grant" {
  database_name = "database"
  schema_name   = "schema"
  stage_name    = "stage"

  privilege = "USAGE"

  roles = ["role1", "role2"]

  on_future         = false
  with_grant_option = false
}
