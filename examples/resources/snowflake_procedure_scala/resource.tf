resource "snowflake_procedure_scala" "w" {
  database = "Database"
  schema   = "Schema"
  name     = "Name"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name = "x"
  }
  return_type = "VARCHAR(100)"
  handler = "TestFunc.echoVarchar"
  procedure_definition = "\n\timport com.snowflake.snowpark_java.Session\n\n\tclass TestFunc {\n\t\tdef echoVarchar(session : Session, x : String): String = {\n\t\t\treturn x\n\t\t}\n\t}\n"
  runtime_version = "2.12"
  snowpark_package = "1.14.0"
}
