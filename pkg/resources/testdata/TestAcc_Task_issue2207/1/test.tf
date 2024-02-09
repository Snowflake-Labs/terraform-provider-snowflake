resource "snowflake_task" "root_task" {
  name          = var.root_name
  database      = var.database
  schema        = var.schema
  warehouse     = var.warehouse
  sql_statement = "SELECT 1"
  enabled       = true
  schedule      = "5 MINUTE"
}

resource "snowflake_task" "child_task" {
  name          = var.child_name
  database      = snowflake_task.root_task.database
  schema        = snowflake_task.root_task.schema
  warehouse     = snowflake_task.root_task.warehouse
  sql_statement = "SELECT 1"
  enabled       = true
  after         = [snowflake_task.root_task.name]
  comment       = var.comment
}
