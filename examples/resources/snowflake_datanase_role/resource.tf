resource "snowflake_database_role" "db_role" {
  database = "database"
  name     = "role_1"
  comment  = "my db role"
}
