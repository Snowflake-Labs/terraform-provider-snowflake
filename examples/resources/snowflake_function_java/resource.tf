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
  function_definition = "\n\tclass TestFunc {\n\t\tpublic static String echoVarchar(String x) {\n\t\t\treturn x;\n\t\t}\n\t}\n"
}
