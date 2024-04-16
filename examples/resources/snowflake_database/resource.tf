resource "snowflake_database" "simple" {
  name                        = "testing"
  comment                     = "test comment"
  data_retention_time_in_days = 3
}

resource "snowflake_database" "with_replication" {
  name    = "testing_2"
  comment = "test comment 2"
  replication_configuration {
    accounts             = ["test_account1", "test_account_2"]
    ignore_edition_check = true
  }
}

resource "snowflake_database" "from_replica" {
  name                        = "testing_3"
  comment                     = "test comment"
  data_retention_time_in_days = 3
  from_replica                = "\"org1\".\"account1\".\"primary_db_name\""
}

resource "snowflake_database" "from_share" {
  name    = "testing_4"
  comment = "test comment"
  from_share = {
    provider = "account1_locator"
    share    = "share1"
  }
}
