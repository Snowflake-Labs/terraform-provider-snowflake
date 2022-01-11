resource snowflake_task task {
  comment = "my task"

  database  = "db"
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

resource snowflake_task serverless_task {
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
  after                                    = "preceding_task"
  when                                     = "foo AND bar"
  enabled                                  = true
}
