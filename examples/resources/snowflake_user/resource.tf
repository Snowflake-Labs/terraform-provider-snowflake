resource "snowflake_user" "user" {
  name         = "Snowflake User"
  login_name   = "snowflake_user"
  comment      = "A user of snowflake."
  password     = "secret"
  disabled     = false
  display_name = "Snowflake User"
  email        = "user@snowflake.example"
  first_name   = "Snowflake"
  last_name    = "User"

  default_warehouse       = "warehouse"
  default_secondary_roles = "ALL"
  default_role            = "role1"

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."

  must_change_password = false
}
