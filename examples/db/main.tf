resource "snowflake_database" "test" {
  name                        = "testing"
  comment                     = "test comment"
  data_retention_time_in_days = 3
}

resource "snowflake_database" "test2" {
  name    = "testing_2"
  comment = "test comment 2"
}
