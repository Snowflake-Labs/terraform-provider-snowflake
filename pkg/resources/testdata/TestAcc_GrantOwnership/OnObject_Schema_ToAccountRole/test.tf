resource "snowflake_account_role" "test" {
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
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "SCHEMA"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
  }
}
