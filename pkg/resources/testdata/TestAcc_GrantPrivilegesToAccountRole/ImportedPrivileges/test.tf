resource "snowflake_shared_database" "test" {
  name       = var.shared_database_name
  from_share = var.external_share_name
}

resource "snowflake_account_role" "test" {
  name = var.role_name
}

resource "snowflake_grant_privileges_to_account_role" "test" {
  depends_on        = [snowflake_shared_database.test, snowflake_account_role.test]
  account_role_name = "\"${var.role_name}\""
  privileges        = var.privileges
  on_account_object {
    object_type = "DATABASE"
    object_name = "\"${var.shared_database_name}\""
  }
}
