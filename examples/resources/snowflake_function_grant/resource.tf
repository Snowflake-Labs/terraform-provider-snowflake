resource "snowflake_function_grant" "grant" {
  database_name       = "database"
  schema_name         = "schema"
  function_name       = "function"
  argument_data_types = ["array", "string"]
  privilege           = "USAGE"
  roles               = ["role1", "role2"]
  shares              = ["share1", "share2"]
  on_future           = false
  with_grant_option   = false
}
