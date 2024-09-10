resource "snowflake_schema" "test" {
  name                        = var.schema
  database                    = var.database
  data_retention_time_in_days = var.schema_data_retention_time
}
