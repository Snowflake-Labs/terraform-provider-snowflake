resource "snowflake_sequence_grant" "grant" {
  database_name = "database"
  schema_name   = "schema"
  sequence_name = "sequence"

  privilege = "SELECT"
  roles     = ["role1", "role2"]

  on_future         = false
  with_grant_option = false
}
