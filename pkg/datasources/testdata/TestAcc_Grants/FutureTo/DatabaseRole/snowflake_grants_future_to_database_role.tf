resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_database_role" "test" {
  name     = var.database_role
  database = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "\"${snowflake_database.test.name}\".\"${snowflake_database_role.test.name}\""
  privileges         = ["CREATE TABLE"]

  on_schema {
    future_schemas_in_database = snowflake_database.test.name
  }
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_privileges_to_database_role.test]

  future_grants_to {
    database_role = "\"${snowflake_database.test.name}\".\"${snowflake_database_role.test.name}\""
  }
}
