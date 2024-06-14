resource "snowflake_database_old" "db" {
  name    = var.db
  comment = "test comment"
}
