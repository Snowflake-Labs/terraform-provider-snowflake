resource "snowflake_external_table_grant" "grant" {
  database_name       = "database"
  schema_name         = "schema"
  external_table_name = "external_table"

  privilege = "SELECT"
  roles     = ["role1", "role2"]

  shares = ["share1", "share2"]

  on_future         = false
  with_grant_option = false
}
