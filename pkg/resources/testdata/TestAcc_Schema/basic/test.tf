resource "snowflake_schema" "test" {
  name     = var.name
  database = var.database
}
