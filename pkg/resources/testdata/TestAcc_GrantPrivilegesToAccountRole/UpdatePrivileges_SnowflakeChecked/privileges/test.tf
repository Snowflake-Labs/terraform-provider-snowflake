resource "snowflake_grant_privileges_to_account_role" "test" {
 role_name = var.name
  privileges         = var.privileges
 on_account_object {
  object_type = "DATABASE"
  object_name = var.database
 }
}
