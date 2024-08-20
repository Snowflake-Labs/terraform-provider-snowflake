resource "snowflake_schema" "test" {
  name                = var.name
  database            = var.database
  with_managed_access = var.with_managed_access
  is_transient        = var.is_transient
  comment             = var.comment

  data_retention_time_in_days                   = var.data_retention_time_in_days
  max_data_extension_time_in_days               = var.max_data_extension_time_in_days
  external_volume                               = var.external_volume
  catalog                                       = var.catalog
  replace_invalid_characters                    = var.replace_invalid_characters
  default_ddl_collation                         = var.default_ddl_collation
  storage_serialization_policy                  = var.storage_serialization_policy
  log_level                                     = var.log_level
  trace_level                                   = var.trace_level
  suspend_task_after_num_failures               = var.suspend_task_after_num_failures
  task_auto_retry_attempts                      = var.task_auto_retry_attempts
  user_task_managed_initial_warehouse_size      = var.user_task_managed_initial_warehouse_size
  user_task_timeout_ms                          = var.user_task_timeout_ms
  user_task_minimum_trigger_interval_in_seconds = var.user_task_minimum_trigger_interval_in_seconds
  quoted_identifiers_ignore_case                = var.quoted_identifiers_ignore_case
  enable_console_output                         = var.enable_console_output
  pipe_execution_paused                         = var.pipe_execution_paused

}
