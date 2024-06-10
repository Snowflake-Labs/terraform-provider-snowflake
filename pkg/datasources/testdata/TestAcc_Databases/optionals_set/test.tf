resource "snowflake_standard_database" "test" {
  name    = var.name
  comment = var.comment
  replication {
    enable_to_account {
      account_identifier = var.account_identifier
      with_failover      = true
    }
    ignore_edition_check = true
  }
}

data "snowflake_databases" "test" {
  depends_on  = [snowflake_standard_database.test]
  like        = var.name
  starts_with = var.name
  limit {
    rows = 1
  }
}
