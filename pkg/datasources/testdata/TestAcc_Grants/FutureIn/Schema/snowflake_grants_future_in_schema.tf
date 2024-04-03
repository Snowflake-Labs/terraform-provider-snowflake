data "snowflake_current_role" "test" {}

resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = data.snowflake_current_role.test.name
  privileges        = ["INSERT"]

  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_schema          = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
    }
  }
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_privileges_to_account_role.test]

  future_grants_in {
    schema = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
  }
}
