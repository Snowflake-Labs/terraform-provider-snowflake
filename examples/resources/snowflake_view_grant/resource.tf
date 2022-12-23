resource "snowflake_view_grant" "grant" {
  database_name = "database"
  schema_name   = "schema"
  view_name     = "view"

  privilege = "SELECT"
  roles     = ["role1", "role2"]

  shares = ["share1", "share2"]

  on_future         = false
  with_grant_option = false
}

/*
Snowflake view grant is an object level grant, not a schema level grant. To add schema level
grants, use the `snowflake_schema_grant` resource
*/

resource "snowflake_schema_grant" "grant" {
  database_name = "database"
  schema_name   = "schema"

  privilege = "USAGE"
  roles     = ["role1", "role2"]
}
