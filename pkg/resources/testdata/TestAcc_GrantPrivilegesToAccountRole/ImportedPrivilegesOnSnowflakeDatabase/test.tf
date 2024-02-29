resource "snowflake_database" "test" {
  name                        = var.shared_database_name
  data_retention_time_in_days = 0
  from_share = {
    provider = var.account_name
    share    = var.share_name
  }
}

resource "snowflake_role" "test" {
  name = var.role_name
}

resource "snowflake_grant_privileges_to_account_role" "test" {
  depends_on        = [snowflake_database.test, snowflake_role.test]
  account_role_name = "\"${var.role_name}\""
  privileges        = var.privileges
  on_account_object {
    object_type = "APPLICATION"
    object_name = "\"${var.shared_database_name}\""
  }
}
