# Minimal
resource "snowflake_function_javascript" "minimal" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "my_javascript_function"
  arguments {
    arg_data_type = "VARIANT"
    arg_name      = "x"
  }
  function_definition = <<EOF
    return x;
  EOF
  return_type         = "VARIANT"
}
