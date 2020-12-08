resource snowflake_task task {
  comment = "my task"

  database  = "db"
  schema    = "schema"
  warehouse = "warehouse"

  name          = "task"
  schedule      = "10"
  sql_statement = "select * from foo;"

  session_parameters = {
    "foo" : "bar",
  }

  user_task_timeout_ms = 10000
  after                = "preceding_task"
  when                 = "foo AND bar"
  enabled              = true
}
