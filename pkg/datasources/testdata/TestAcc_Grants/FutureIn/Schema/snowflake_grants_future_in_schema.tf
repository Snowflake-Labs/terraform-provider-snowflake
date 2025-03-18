data "snowflake_current_role" "test" {}

resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = data.snowflake_current_role.test.name
  privileges        = ["INSERT"]

  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_schema          = "\"${var.database}\".\"${var.schema}\""
    }
  }
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_privileges_to_account_role.test]

  future_grants_in {
    schema = "\"${var.database}\".\"${var.schema}\""
  }
}
