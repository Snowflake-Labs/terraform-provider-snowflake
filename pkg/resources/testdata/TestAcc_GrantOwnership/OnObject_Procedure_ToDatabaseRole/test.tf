resource "snowflake_database" "test" {
  name = var.database_name
}

resource "snowflake_database_role" "test" {
  name     = var.database_role_name
  database = snowflake_database.test.name
}

resource "snowflake_schema" "test" {
  name     = var.schema_name
  database = snowflake_database.test.name
}

resource "snowflake_procedure_javascript" "test" {
  name                = var.procedure_name
  database            = snowflake_database.test.name
  schema              = snowflake_schema.test.name
  return_type         = "FLOAT"
  execute_as          = "CALLER"
  return_behavior     = "VOLATILE"
  null_input_behavior = "RETURNS NULL ON NULL INPUT"
  statement           = <<EOT
var X=1
return X
EOT
}

resource "snowflake_grant_ownership" "test" {
  database_role_name = snowflake_database_role.test.fully_qualified_name
  on {
    object_type = "PROCEDURE"
    object_name = snowflake_procedure_javascript.test.fully_qualified_name
  }
}
