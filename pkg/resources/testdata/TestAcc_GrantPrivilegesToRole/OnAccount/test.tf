resource "snowflake_grant_privileges_to_role" "test" {
  role_name         = var.name
  privileges        = var.privileges
  on_account        = true
  with_grant_option = var.with_grant_option
}
