resource "snowflake_procedure_python" "w" {
  database = "Database"
  schema   = "Schema"
  name     = "Name"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  handler              = "echoVarchar"
  procedure_definition = "\ndef echoVarchar(x):\n\tresult = \"\"\n\tfor a in range(5):\n\t\tresult += x\n\treturn result\n"
  runtime_version      = "3.8"
  snowpark_package     = "1.14.0"
}
