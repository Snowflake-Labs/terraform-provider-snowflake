# Minimal
resource "snowflake_function_python" "minimal" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "my_function_function"
  runtime_version = "3.8"
  arguments {
    arg_data_type = "NUMBER(36, 2)"
    arg_name      = "x"
  }
  function_definition = <<EOF
    def some_function(x):
      return x
  EOF
  handler             = "some_function"
  return_type         = "NUMBER(36, 2)"
}
