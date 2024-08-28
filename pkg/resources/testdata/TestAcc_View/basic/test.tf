resource "snowflake_view" "test" {
  name      = var.name
  database  = var.database
  schema    = var.schema
  statement = var.statement
}
