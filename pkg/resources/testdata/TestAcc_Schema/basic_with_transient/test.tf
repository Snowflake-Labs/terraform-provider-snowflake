resource "snowflake_schema" "test" {
  name         = var.name
  database     = var.database
  is_transient = var.is_transient
}
