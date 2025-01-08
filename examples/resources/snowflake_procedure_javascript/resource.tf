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
  procedure_definition = <<EOT
  if (x == 0) {
    return 1;
  } else {
    return 2;
  }
EOT
}
