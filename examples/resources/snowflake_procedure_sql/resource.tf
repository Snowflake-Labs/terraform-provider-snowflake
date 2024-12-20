resource "snowflake_procedure_sql" "w" {
  database = "Database"
  schema   = "Schema"
  name     = "Name"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  procedure_definition = <<EOT
BEGIN
  RETURN message;
END;
EOT
}
