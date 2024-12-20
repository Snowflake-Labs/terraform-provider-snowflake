resource "snowflake_procedure_scala" "w" {
  database = "Database"
  schema   = "Schema"
  name     = "Name"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  handler              = "TestFunc.echoVarchar"
  procedure_definition = <<EOT
  import com.snowflake.snowpark_java.Session
  class TestFunc {
    def echoVarchar(session : Session, x : String): String = {
      return x
    }
  }
EOT
  runtime_version      = "2.12"
  snowpark_package     = "1.14.0"
}
