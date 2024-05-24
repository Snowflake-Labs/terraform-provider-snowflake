resource "snowflake_secondary_database" "test" {
  name                  = var.name
  as_replica_of         = var.as_replica_of
  is_transient          = var.transient
  external_volume       = var.external_volume
  catalog               = var.catalog
  default_ddl_collation = var.default_ddl_collation
  log_level             = var.log_level
  trace_level           = var.trace_level
  comment               = var.comment
}
