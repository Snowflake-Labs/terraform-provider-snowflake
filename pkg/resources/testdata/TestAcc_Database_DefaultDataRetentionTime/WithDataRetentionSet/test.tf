resource "snowflake_database" "test" {
  name                        = var.database
  data_retention_time_in_days = var.database_data_retention_time
}
