data "snowflake_database_role" "db_role" {
  database = "MYDB"
  name     = "DBROLE"
}