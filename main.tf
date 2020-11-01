resource "snowflake_database" "test" {
  name                        = "testing"
  comment                     = "test comment"
  data_retention_time_in_days = 3
}
