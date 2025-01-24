# Test setup.
module "id" {
  source = "../id"
}

variable "resource_count" {
  type = number
}

terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "= 1.0.1"
    }
  }
}

locals {
  id_number_list = {
    for index, val in range(0, var.resource_count) :
    val => tostring(val)
  }
  test_prefix = format("PERFORMANCE_TESTS_%s", module.id.test_id)
}

resource "snowflake_database" "database" {
  count = var.resource_count > 0 ? 1 : 0
  name  = local.test_prefix
}

# basic resource
resource "snowflake_schema" "all_schemas" {
  database = snowflake_database.database[0].name
  for_each = local.id_number_list
  name     = format("perf_basic_%v", each.key)
}

# resource with all fields set (without dependencies)
resource "snowflake_schema" "schema" {
  database = snowflake_database.database[0].name
  for_each = local.id_number_list
  name     = format("perf_complete_%v", each.key)

  with_managed_access                           = true
  is_transient                                  = true
  comment                                       = "my schema"
  data_retention_time_in_days                   = 1
  max_data_extension_time_in_days               = 20
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
