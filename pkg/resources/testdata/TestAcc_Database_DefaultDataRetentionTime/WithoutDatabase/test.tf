resource "snowflake_account_parameter" "test" {
  key   = "DATA_RETENTION_TIME_IN_DAYS"
  value = var.account_data_retention_time
}

resource "snowflake_database" "test" {
  name                        = var.database
}
