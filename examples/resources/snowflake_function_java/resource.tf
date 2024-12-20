resource "snowflake_function_java" "w" {
  database = "Database"
  schema   = "Schema"
  name     = "Name"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type         = "VARCHAR(100)"
  handler             = "TestFunc.echoVarchar"
  function_definition = <<EOT
  class TestFunc {
    public static String echoVarchar(String x) {
      return x;
    }
  }
EOT
}
