resource "snowflake_shared_database" "test" {
  name                  = "shared_database"
  from_share            = "organization_name.account_name.share_name"
  is_transient          = false
  external_volume       = "external_volume_name"
  catalog               = "catalog_name"
  default_ddl_collation = "en_US"
  log_level             = "OFF"
  trace_level           = "OFF"
  comment               = "A shared database"
}
