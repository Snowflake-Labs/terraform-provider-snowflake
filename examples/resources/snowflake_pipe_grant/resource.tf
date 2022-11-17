resource "snowflake_pipe_grant" "grant" {
  database_name = "database"
  schema_name   = "schema"
  pipe_name     = "pipe"

  privilege = "OPERATE"
  roles     = ["role1", "role2"]

  on_future         = false
  with_grant_option = false
}
