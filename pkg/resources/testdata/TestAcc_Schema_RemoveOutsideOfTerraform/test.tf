resource "snowflake_schema" "test" {
  name     = var.schema_name
  database = var.database_name
}