resource "snowflake_role" "test" {
  name = var.account_role_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name

  on {
    all {
      object_type_plural = "PIPES"
      in_database        = var.database
    }
  }
}
