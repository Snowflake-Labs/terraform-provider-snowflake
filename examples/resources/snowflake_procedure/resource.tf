resource "snowflake_schema" "db" {
  name                = "MYDB"
  data_retention_days = 1
}

resource "snowflake_schema" "schema" {
  database            = snowflake_database.db.name
  name                = "MYSCHEMA"
  data_retention_days = 1
}

resource "snowflake_procedure" "proc" {
  name     = "SAMPLEPROC"
  database = snowflake_database.db.name
  schema   = snowflake_schema.schema.name
  arguments {
    name = "arg1"
    type = "varchar"
  }
  arguments {
    name = "arg2"
    type = "DATE"
  }
  comment             = "Procedure with 2 arguments"
  return_type         = "VARCHAR"
  execute_as          = "CALLER"
  return_behavior     = "IMMUTABLE"
  null_input_behavior = "RETURNS NULL ON NULL INPUT"
  statement           = <<EOT
var X=1
return X
EOT
}
