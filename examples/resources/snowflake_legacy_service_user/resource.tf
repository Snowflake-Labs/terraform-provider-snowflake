resource "snowflake_legacy_service_user" "user" {
  name         = "Snowflake Legacy Service User"
  login_name   = "legacy_service_user"
  comment      = "A legacy service user of snowflake."
  password     = "secret"
  disabled     = false
  display_name = "Snowflake Legacy Service User"
  email        = "legacy.service.user@snowflake.example"

  default_warehouse       = "warehouse"
  default_secondary_roles = "ALL"
  default_role            = "role1"

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."

  must_change_password = true
}
