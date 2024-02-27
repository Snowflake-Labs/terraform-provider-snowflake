resource "snowflake_database" "test" {
  name                        = var.database
  data_retention_time_in_days = var.database_data_retention_time
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}
