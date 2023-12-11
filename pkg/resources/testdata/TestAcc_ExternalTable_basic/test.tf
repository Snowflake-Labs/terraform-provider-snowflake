resource "snowflake_storage_integration" "i" {
  name = var.name
  storage_allowed_locations = [var.location]
  storage_provider = "S3"
  storage_aws_role_arn = var.aws_arn
}

resource "snowflake_stage" "test" {
  name = var.name
  url = var.location
  database = var.database
  schema = var.schema
  storage_integration = snowflake_storage_integration.i.name
}

resource "snowflake_external_table" "test_table" {
  name = var.name
  database = var.database
  schema = var.schema
  comment  = "Terraform acceptance test"
  column {
    name = "column1"
    type = "STRING"
    as   = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
  }
  column {
    name = "column2"
    type = "TIMESTAMP_NTZ(9)"
    as   = "($1:\"CreatedDate\"::timestamp)"
  }
  file_format = "TYPE = CSV"
  location = "@\"${var.database}\".\"${var.schema}\".\"${snowflake_stage.test.name}\""
}
