data "snowflake_current_role" "test" {}

resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_database_role" "test" {
  name     = var.database_role
  database = snowflake_database.test.name
}

resource "snowflake_grant_database_role" "test" {
  database_role_name = "\"${snowflake_database.test.name}\".\"${snowflake_database_role.test.name}\""
  parent_role_name   = data.snowflake_current_role.test.name
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_database_role.test]

  grants_of {
    database_role = "\"${snowflake_database.test.name}\".\"${snowflake_database_role.test.name}\""
  }
}
