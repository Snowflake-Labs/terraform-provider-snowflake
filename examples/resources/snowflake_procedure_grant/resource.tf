resource "snowflake_procedure_grant" "grant" {
  database_name       = "database"
  schema_name         = "schema"
  procedure_name      = "procedure"
  argument_data_types = ["array", "string"]
  privilege           = "USAGE"
  roles               = ["role1", "role2"]
  shares              = ["share1", "share2"]
  on_future           = false
  with_grant_option   = false
}
