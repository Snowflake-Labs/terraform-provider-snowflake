resource "snowflake_service_user" "service_user" {
  name         = "Snowflake Service User"
  login_name   = "service_user"
  comment      = "A service user of snowflake."
  disabled     = false
  display_name = "Snowflake Service User"
  email        = "service_user@snowflake.example"

  default_warehouse       = "warehouse"
  default_secondary_roles = "ALL"
  default_role            = "role1"

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."
}
