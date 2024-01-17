resource "snowflake_grant_privileges_to_role" "test" {
  privileges = ["CREATE SCHEMA"]
  role_name  = "\"${var.name}\""
  on_account_object {
    object_type = "DATABASE"
    object_name = var.database
  }
}