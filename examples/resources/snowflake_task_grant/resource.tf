resource "snowflake_task_grant" "grant" {
  database_name = "database"
  schema_name   = "schema"
  task_name     = "task"

  privilege = "OPERATE"
  roles     = ["role1", "role2"]

  on_future         = false
  with_grant_option = false
}
