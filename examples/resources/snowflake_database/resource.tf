## Minimal
resource "snowflake_database" "primary" {
  name = "database_name"
}

## Complete (with every optional set)
resource "snowflake_database" "primary" {
  name         = "database_name"
  is_transient = false
  comment      = "my standard database"

  data_retention_time_in_days                   = 10
  data_retention_time_in_days_save              = 10
  max_data_extension_time_in_days               = 20
  external_volume                               = "<external_volume_name>"
  catalog                                       = "<catalog_name>"
  replace_invalid_characters                    = false
  default_ddl_collation                         = "en_US"
  storage_serialization_policy                  = "COMPATIBLE"
  log_level                                     = "INFO"
  trace_level                                   = "ALWAYS"
  suspend_task_after_num_failures               = 10
  task_auto_retry_attempts                      = 10
  user_task_managed_initial_warehouse_size      = "LARGE"
  user_task_timeout_ms                          = 3600000
  user_task_minimum_trigger_interval_in_seconds = 120
  quoted_identifiers_ignore_case                = false
  enable_console_output                         = false

  replication {
    enable_to_account {
      account_identifier = "<secondary_account_organization_name>.<secondary_account_name>"
      with_failover      = true
    }
    ignore_edition_check = true
  }
}

## Replication with for_each
locals {
  replication_configs = [
    {
      account_identifier = "<secondary_account_organization_name>.<secondary_account_name>"
      with_failover      = true
    },
    {
      account_identifier = "<secondary_account_organization_name>.<secondary_account_name>"
      with_failover      = true
    },
  ]
}

resource "snowflake_database" "primary" {
  name     = "database_name"
  for_each = local.replication_configs

  replication {
    enable_to_account    = each.value
    ignore_edition_check = true
  }
}
