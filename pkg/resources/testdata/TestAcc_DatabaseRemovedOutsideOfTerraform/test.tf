resource "snowflake_database" "db" {
  name    = var.db
  comment = "test comment"
}
