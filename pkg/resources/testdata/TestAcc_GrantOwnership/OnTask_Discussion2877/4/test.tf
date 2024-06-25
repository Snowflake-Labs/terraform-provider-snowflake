resource "snowflake_role" "test" {
  name = var.account_role_name
}

resource "snowflake_task" "test" {
  database      = var.database
  schema        = var.schema
  name          = var.task
  warehouse     = var.warehouse
  sql_statement = "SELECT CURRENT_TIMESTAMP"
}

resource "snowflake_task" "child" {
  database      = var.database
  schema        = var.schema
  name          = var.child
  warehouse     = var.warehouse
  after         = [snowflake_task.test.name]
  sql_statement = "SELECT CURRENT_TIMESTAMP"
}

resource "snowflake_grant_ownership" "test" {
  depends_on        = [snowflake_task.test, snowflake_task.child]
  account_role_name = snowflake_role.test.name

  on {
    all {
      object_type_plural = "TASKS"
      in_schema          = "\"${var.database}\".\"${var.schema}\""
    }
  }
}
