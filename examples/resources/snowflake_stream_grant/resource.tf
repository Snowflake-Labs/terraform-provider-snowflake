resource "snowflake_stream_grant" "grant" {
  database_name = "database"
  schema_name   = "schema"
  stream_name   = "view"

  privilege = "SELECT"
  roles     = ["role1", "role2"]

  on_future         = false
  with_grant_option = false
}
