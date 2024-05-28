resource "snowflake_secondary_database" "test" {
  name          = var.name
  as_replica_of = var.as_replica_of
  is_transient  = var.transient

  data_retention_time_in_days {
    value = var.data_retention_time_in_days
  }

  max_data_extension_time_in_days {
    value = var.max_data_extension_time_in_days
  }

  external_volume              = var.external_volume
  catalog                      = var.catalog
  replace_invalid_characters   = var.replace_invalid_characters
  default_ddl_collation        = var.default_ddl_collation
  storage_serialization_policy = var.storage_serialization_policy
  log_level                    = var.log_level
  trace_level                  = var.trace_level
  comment                      = var.comment
}
