resource "snowflake_role" "test" {
  name = var.account_role_name
}

resource "snowflake_database" "test" {
  name = var.database_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    object_type = ""
    object_name = snowflake_database.test.name
  }
}
