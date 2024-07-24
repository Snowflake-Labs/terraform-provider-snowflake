# basic resource
resource "snowflake_schema" "schema" {
  name     = "schema_name"
  database = "database_name"
}

# resource with all fields set
resource "snowflake_schema" "schema" {
  name                = "schema_name"
  database            = "database_name"
  with_managed_access = true
  is_transient        = true
  comment             = "my schema"

  data_retention_time_in_days                   = 10
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
  pipe_execution_paused                         = false

}
