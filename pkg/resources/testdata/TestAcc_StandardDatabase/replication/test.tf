resource "snowflake_standard_database" "test" {
  name = var.name

  replication {
    enable_for_account {
      account_identifier = var.account_identifier
      with_failover      = var.with_failover
    }

    ignore_edition_check = var.ignore_edition_check
  }
}
