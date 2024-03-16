resource "snowflake_role" "test" {
  name = var.account_role_name
}

resource "snowflake_database" "test" {
  name = var.database_name
}

resource "snowflake_schema" "test" {
  name     = var.schema_name
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    future {
      object_type_plural = "TABLES"
      in_schema          = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
    }
  }
}
