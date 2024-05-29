resource "snowflake_shared_database" "test" {
  name                         = var.name
  from_share                   = var.from_share
  external_volume              = var.external_volume
  catalog                      = var.catalog
  replace_invalid_characters   = var.replace_invalid_characters
  default_ddl_collation        = var.default_ddl_collation
  storage_serialization_policy = var.storage_serialization_policy
  log_level                    = var.log_level
  trace_level                  = var.trace_level
  comment                      = var.comment
}
