# minimal
resource "snowflake_service_user" "minimal" {
  name = "Snowflake Service User - minimal"
}

# with all attributes set
resource "snowflake_service_user" "service_user" {
  name         = "Snowflake Service User"
  login_name   = "service_user"
  comment      = "A service user of snowflake."
  disabled     = "false"
  display_name = "Snowflake Service User"
  email        = "service_user@snowflake.example"

  default_warehouse              = "warehouse"
  default_secondary_roles_option = "ALL"
  default_role                   = "role1"
  default_namespace              = "some.namespace"

  mins_to_unlock = 9
  days_to_expiry = 8

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."
}

# all parameters set on the resource level
resource "snowflake_service_user" "u" {
  name = "Snowflake Service User with all parameters"

  abort_detached_query                          = true
  autocommit                                    = false
  binary_input_format                           = "UTF8"
  binary_output_format                          = "BASE64"
  client_memory_limit                           = 1024
  client_metadata_request_use_connection_ctx    = true
  client_prefetch_threads                       = 2
  client_result_chunk_size                      = 48
  client_result_column_case_insensitive         = true
  client_session_keep_alive                     = true
  client_session_keep_alive_heartbeat_frequency = 2400
  client_timestamp_type_mapping                 = "TIMESTAMP_NTZ"
  date_input_format                             = "YYYY-MM-DD"
  date_output_format                            = "YY-MM-DD"
  enable_unload_physical_type_optimization      = false
  enable_unredacted_query_syntax_error          = true
  error_on_nondeterministic_merge               = false
  error_on_nondeterministic_update              = true
  geography_output_format                       = "WKB"
  geometry_output_format                        = "WKB"
  jdbc_treat_decimal_as_int                     = false
  jdbc_treat_timestamp_ntz_as_utc               = true
  jdbc_use_session_timezone                     = false
  json_indent                                   = 4
  lock_timeout                                  = 21222
  log_level                                     = "ERROR"
  multi_statement_count                         = 0
  network_policy                                = "BVYDGRAT_0D5E3DD1_F644_03DE_318A_1179886518A7"
  noorder_sequence_as_default                   = false
  odbc_treat_decimal_as_int                     = true
  prevent_unload_to_internal_stages             = true
  query_tag                                     = "some_tag"
  quoted_identifiers_ignore_case                = true
  rows_per_resultset                            = 2
  search_path                                   = "$public, $current"
  simulated_data_sharing_consumer               = "some_consumer"
  statement_queued_timeout_in_seconds           = 10
  statement_timeout_in_seconds                  = 10
  strict_json_output                            = true
  s3_stage_vpce_dns_name                        = "vpce-id.s3.region.vpce.amazonaws.com"
  time_input_format                             = "HH24:MI"
  time_output_format                            = "HH24:MI"
  timestamp_day_is_always_24h                   = true
  timestamp_input_format                        = "YYYY-MM-DD"
  timestamp_ltz_output_format                   = "YYYY-MM-DD HH24:MI:SS"
  timestamp_ntz_output_format                   = "YYYY-MM-DD HH24:MI:SS"
  timestamp_output_format                       = "YYYY-MM-DD HH24:MI:SS"
  timestamp_type_mapping                        = "TIMESTAMP_LTZ"
  timestamp_tz_output_format                    = "YYYY-MM-DD HH24:MI:SS"
  timezone                                      = "Europe/Warsaw"
  trace_level                                   = "ON_EVENT"
  transaction_abort_on_error                    = true
  transaction_default_isolation_level           = "READ COMMITTED"
  two_digit_century_start                       = 1980
  unsupported_ddl_action                        = "FAIL"
  use_cached_result                             = false
  week_of_year_policy                           = 1
  week_start                                    = 1
}
