resource "snowflake_database" "db" {
  name = "db1"
}

resource "snowflake_failover_group" "source_failover_group" {
  name                      = "FG1"
  object_types              = ["WAREHOUSES", "DATABASES", "INTEGRATIONS", "ROLES"]
  allowed_accounts          = ["<org_name>.<target_account_name1>", "<org_name>.<target_account_name2>"]
  allowed_databases         = [snowflake_database.db.name]
  allowed_integration_types = ["SECURITY INTEGRATIONS"]
  replication_schedule {
    cron {
      expression = "0 0 10-20 * TUE,THU"
      time_zone  = "UTC"
    }

    // replication_schedule could also be specified with interval instead of cron
    // interval = 10
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
    name                = snowflake_failover_group.source_failover_group.name
  }
}
