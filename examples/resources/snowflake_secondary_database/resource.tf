# 1. Preparing primary database
resource "snowflake_database" "primary" {
  provider = primary_account # notice the provider fields
  name     = "database_name"
  replication_configuration {
    accounts             = ["<secondary_account_organization_name>.<secondary_account_name>"]
    ignore_edition_check = true
  }
}

# 2. Creating secondary database
resource "snowflake_secondary_database" "test" {
  provider      = secondary_account
  name          = snowflake_database.primary.name # It's recommended to give a secondary database the same name as its primary database
  as_replica_of = "<primary_account_organization_name>.<primary_account_name>.${snowflake_database.primary.name}"
  is_transient  = false

  data_retention_time_in_days {
    value = 10
  }

  max_data_extension_time_in_days {
    value = 20
  }

  external_volume              = "external_volume_name"
  catalog                      = "catalog_name"
  replace_invalid_characters   = false
  default_ddl_collation        = "en_US"
  storage_serialization_policy = "OPTIMIZED"
  log_level                    = "OFF"
  trace_level                  = "OFF"
  comment                      = "A secondary database"
}
