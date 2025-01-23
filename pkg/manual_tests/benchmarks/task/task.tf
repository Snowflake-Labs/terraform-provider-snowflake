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

resource "snowflake_schema" "schema" {
  database = snowflake_database.database[0].name
  count    = var.resource_count > 0 ? 1 : 0
  name     = "PERFORMANCE_TESTS"
}

# Basic standalone task
resource "snowflake_task" "task" {
  database  = snowflake_database.database[0].name
  schema    = snowflake_schema.schema[0].name
  for_each  = local.id_number_list
  name      = format("perf_basic_%v", each.key)
  warehouse = "SNOWFLAKE"
  started   = true
  schedule {
    minutes = 5
  }
  sql_statement = "select 1"
}

# Complete standalone task
resource "snowflake_task" "test" {
  database      = snowflake_database.database[0].name
  schema        = snowflake_schema.schema[0].name
  for_each      = local.id_number_list
  name          = format("perf_complete_%v", each.key)
  warehouse     = "SNOWFLAKE"
  started       = true
  sql_statement = "select 1"

  config                      = "{\"key\":\"value\"}"
  allow_overlapping_execution = true
  comment                     = "complete task"

  schedule {
    minutes = 10
  }

  # Session Parameters
  suspend_task_after_num_failures               = 10
  task_auto_retry_attempts                      = 0
  user_task_minimum_trigger_interval_in_seconds = 30
  user_task_timeout_ms                          = 3600000
  abort_detached_query                          = false
  autocommit                                    = true
  binary_input_format                           = "HEX"
  binary_output_format                          = "HEX"
  client_memory_limit                           = 1536
  client_metadata_request_use_connection_ctx    = false
  client_prefetch_threads                       = 4
  client_result_chunk_size                      = 160
  client_result_column_case_insensitive         = false
  client_session_keep_alive                     = false
  client_session_keep_alive_heartbeat_frequency = 3600
  client_timestamp_type_mapping                 = "TIMESTAMP_LTZ"
  date_input_format                             = "AUTO"
  date_output_format                            = "YYYY-MM-DD"
  enable_unload_physical_type_optimization      = true
  error_on_nondeterministic_merge               = true
  error_on_nondeterministic_update              = false
  geography_output_format                       = "GeoJSON"
  geometry_output_format                        = "GeoJSON"
  jdbc_use_session_timezone                     = true
  json_indent                                   = 2
  lock_timeout                                  = 43200
  log_level                                     = "OFF"
  multi_statement_count                         = 1
  noorder_sequence_as_default                   = true
  odbc_treat_decimal_as_int                     = false
  query_tag                                     = ""
  quoted_identifiers_ignore_case                = false
  rows_per_resultset                            = 0
  s3_stage_vpce_dns_name                        = ""
  search_path                                   = "$current, $public"
  statement_queued_timeout_in_seconds           = 0
  statement_timeout_in_seconds                  = 172800
  strict_json_output                            = false
  timestamp_day_is_always_24h                   = false
  timestamp_input_format                        = "AUTO"
  timestamp_ltz_output_format                   = ""
  timestamp_ntz_output_format                   = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_output_format                       = "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM"
  timestamp_type_mapping                        = "TIMESTAMP_NTZ"
  timestamp_tz_output_format                    = ""
  timezone                                      = "America/Los_Angeles"
  time_input_format                             = "AUTO"
  time_output_format                            = "HH24:MI:SS"
  trace_level                                   = "OFF"
  transaction_abort_on_error                    = false
  transaction_default_isolation_level           = "READ COMMITTED"
  two_digit_century_start                       = 1970
  unsupported_ddl_action                        = "ignore"
  use_cached_result                             = true
  week_of_year_policy                           = 0
  week_start                                    = 0
}
