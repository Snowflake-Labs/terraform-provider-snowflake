resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_database_role" "test" {
  name     = var.database_role
  database = snowflake_database.test.name
}

data "snowflake_grants" "test" {
  grants_to {
    database_role = "\"${snowflake_database.test.name}\".\"${snowflake_database_role.test.name}\""
  }
}
