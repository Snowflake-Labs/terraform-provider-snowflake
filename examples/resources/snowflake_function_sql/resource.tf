# Minimal
resource "snowflake_function_sql" "minimal" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "my_sql_function"
  arguments {
    arg_data_type = "FLOAT"
    arg_name      = "arg_name"
  }
  function_definition = <<EOF
    arg_name
  EOF
  return_type         = "FLOAT"
}

# Complete
resource "snowflake_function_sql" "complete" {
  database  = snowflake_database.test.name
  schema    = snowflake_schema.test.name
  name      = "my_sql_function"
  is_secure = "false"
  arguments {
    arg_data_type = "FLOAT"
    arg_name      = "arg_name"
  }
  function_definition     = <<EOF
    arg_name
  EOF
  return_type             = "FLOAT"
  return_results_behavior = "VOLATILE"
  comment                 = "some comment"
}
