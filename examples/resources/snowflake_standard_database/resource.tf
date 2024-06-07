resource "snowflake_standard_database" "primary" {
  name         = "database_name"
  is_transient = false
  comment      = "my standard database"

  data_retention_time_in_days                   = 10
  max_data_extension_time_in_days               = 20
  external_volume                               = "<external_volume_name>"
  catalog                                       = "<external_volume_name>"
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
    enable_for_account {
      account_identifier = "<secondary_account_organization_name>.<secondary_account_name>"
      with_failover      = true
    }
    ignore_edition_check = true
  }
}
