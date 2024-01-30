resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_stage" "test" {
  name                = var.name
  schema              = snowflake_schema.test.name
  database            = snowflake_database.test.name
  comment             = var.comment
  url                 = var.url
  storage_integration = var.storage_integration
  credentials         = var.credentials
  encryption          = var.encryption
  file_format         = var.file_format
}
