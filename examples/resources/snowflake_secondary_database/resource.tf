resource "snowflake_secondary_database" "test" {
  name          = "database_name"
  as_replica_of = "organization_name.account_name.primary_database_name"
  is_transient  = false

  data_retention_time_in_days {
    value = 10
  }

  max_data_extension_time_in_days {
    value = 20
  }

  external_volume       = "external_volume_name"
  catalog               = "catalog_name"
  default_ddl_collation = "en_US"
  log_level             = "OFF"
  trace_level           = "OFF"
  comment               = "A secondary database"
}
