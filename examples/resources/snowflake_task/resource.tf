resource "snowflake_task" "task" {
  comment = "my task"

  database  = "database"
  schema    = "schema"
  warehouse = "warehouse"

  name          = "task"
  schedule      = "10 MINUTE"
  sql_statement = "select * from foo;"

  session_parameters = {
    "foo" : "bar",
  }

  user_task_timeout_ms = 10000
  after                = "preceding_task"
  when                 = "foo AND bar"
  enabled              = true
}

resource "snowflake_task" "serverless_task" {
  comment = "my serverless task"

  database = "db"
  schema   = "schema"

  name          = "serverless_task"
  schedule      = "10 MINUTE"
  sql_statement = "select * from foo;"

  session_parameters = {
    "foo" : "bar",
  }

  user_task_timeout_ms                     = 10000
  user_task_managed_initial_warehouse_size = "XSMALL"
  after                                    = [snowflake_task.task.name]
  when                                     = "foo AND bar"
  enabled                                  = true
}

resource "snowflake_task" "test_task" {
  comment = "task with allow_overlapping_execution"

  database = "database"
  schema   = "schema"

  name          = "test_task"
  sql_statement = "select 1 as c;"

  allow_overlapping_execution = true
  enabled                     = true
}
