resource "snowflake_grant_privileges_to_account_role" "test" {
  role_name      = var.name
  all_privileges = var.all_privileges
  always_apply   = var.always_apply
  on_account_object {
    object_type = "DATABASE"
    object_name = var.database
  }
}
