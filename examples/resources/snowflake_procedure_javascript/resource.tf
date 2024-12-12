# basic
resource "snowflake_procedure_javascript" "basic" {
  database = "Database"
  schema   = "Schema"
  name     = "Name"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  procedure_definition = "\n\tif (x \u003c= 0) {\n\t\treturn 1;\n\t} else {\n\t\tvar result = 1;\n\t\tfor (var i = 2; i \u003c= x; i++) {\n\t\t\tresult = result * i;\n\t\t}\n\t\treturn result;\n\t}\n"
}
