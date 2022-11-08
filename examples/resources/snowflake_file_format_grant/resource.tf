resource "snowflake_file_format_grant" "grant" {
  database_name    = "database"
  schema_name      = "schema"
  file_format_name = "file_format"

  privilege = "SELECT"
  roles     = ["role1", "role2"]

  on_future         = false
  with_grant_option = false
}
