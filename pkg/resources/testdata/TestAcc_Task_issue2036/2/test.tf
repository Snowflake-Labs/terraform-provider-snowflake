resource "snowflake_task" "test_task" {
  name          = var.name
  database      = var.database
  schema        = var.schema
  warehouse     = var.warehouse
  sql_statement = "SELECT 1"
  enabled       = true
  schedule      = "5 MINUTE"
  when          = "TRUE"
}
