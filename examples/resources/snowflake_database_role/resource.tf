resource "snowflake_database" "test_database" {
  name = "database_name"
}

resource "snowflake_database_role" "test_database_role" {
  database = snowflake_database.test_database.fully_qualified_name
  name     = "database_role_name"
  comment  = "my database role"
}
