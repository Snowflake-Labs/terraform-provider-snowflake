data "snowflake_database_role" "db_role" {
  database = "MYDB"
  role     = "DBROLE"
}