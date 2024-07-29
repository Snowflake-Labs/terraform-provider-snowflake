resource "snowflake_schema" "test" {
  name                  = var.name
  database              = var.database
  pipe_execution_paused = var.pipe_execution_paused
}
