# Minimal
resource "snowflake_function_javascript" "minimal" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = "my_javascript_function"
  arguments {
    arg_data_type = "VARIANT"
    arg_name      = "x"
  }
  return_type         = "VARIANT"
  function_definition = <<EOT
  if (x == 0) {
    return 1;
  } else {
    return 2;
  }
EOT
}
