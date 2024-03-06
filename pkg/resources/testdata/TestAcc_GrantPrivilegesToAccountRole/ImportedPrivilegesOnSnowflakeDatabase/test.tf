resource "snowflake_role" "test" {
  name = var.role_name
}

resource "snowflake_grant_privileges_to_account_role" "test" {
  depends_on        = [snowflake_role.test]
  account_role_name = "\"${var.role_name}\""
  privileges        = var.privileges
  on_account_object {
    object_type = "DATABASE"
    object_name = "\"SNOWFLAKE\""
  }
}
