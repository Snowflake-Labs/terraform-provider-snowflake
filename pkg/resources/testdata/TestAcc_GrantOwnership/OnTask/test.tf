resource "snowflake_role" "test" {
  name = var.account_role_name
}

resource "snowflake_task" "test" {
  database       = var.database
  schema         = var.schema
  name           = var.task
  sql_statement = "SELECT CURRENT_TIMESTAMP"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name

  on {
    object_type = "TASK"
    object_name = "\"${snowflake_task.test.database}\".\"${snowflake_task.test.schema}\".\"${snowflake_task.test.name}\""
  }
}
