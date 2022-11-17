resource "snowflake_database" "db" {
  name = "db1"
}

resource "snowflake_failover_group" "source_failover_group" {
  name                      = "FG1"
  object_types              = ["WAREHOUSES", "DATABASES", "INTEGRATIONS", "ROLES"]
  allowed_accounts          = ["<account1>", "<account2>"]
  allowed_databases         = [snowflake_database.db.name]
  allowed_integration_types = ["SECURITY INTEGRATIONS"]
  replication_schedule {
    cron {
      expression = "0 0 10-20 * TUE,THU"
      time_zone  = "UTC"
    }
  }
}

provider "snowflake" {
  alias = "account2"
}

resource "snowflake_failover_group" "target_failover_group" {
  provider = snowflake.account2
  name     = "FG1"
  from_replica {
    organization_name   = "..."
    source_account_name = "..."
    name                = snowflake_failover_group.fg.name
  }
}
