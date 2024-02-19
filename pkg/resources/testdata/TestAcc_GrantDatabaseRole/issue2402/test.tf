resource "snowflake_database" "database" {
  name = var.database
}

resource "snowflake_database_role" "database_role" {
  database = snowflake_database.database.name
  name     = var.database_role_name
}

resource "snowflake_database_role" "parent_database_role" {
  database = snowflake_database.database.name
  name     = var.parent_database_role_name
}

resource "snowflake_grant_database_role" "g" {
  database_role_name        = "\"${snowflake_database.database.name}\".\"${snowflake_database_role.database_role.name}\""
  parent_database_role_name = "\"${snowflake_database.database.name}\".\"${snowflake_database_role.parent_database_role.name}\""
}
