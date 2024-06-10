# 1. Preparing primary database
resource "snowflake_standard_database" "primary" {
  provider = primary_account # notice the provider fields
  name     = "database_name"
  replication {
    enable_to_account {
      account_identifier = "<secondary_account_organization_name>.<secondary_account_name>"
      with_failover      = true
    }
    ignore_edition_check = true
  }
}

# 2. Creating secondary database
## 2.1. Minimal version
resource "snowflake_secondary_database" "test" {
  provider      = secondary_account
  name          = snowflake_standard_database.primary.name # It's recommended to give a secondary database the same name as its primary database
  as_replica_of = "<primary_account_organization_name>.<primary_account_name>.${snowflake_standard_database.primary.name}"
}

## 2.2. Complete version (with every optional set)
resource "snowflake_secondary_database" "test" {
  provider      = secondary_account
  name          = snowflake_standard_database.primary.name # It's recommended to give a secondary database the same name as its primary database
  is_transient  = false
  as_replica_of = "<primary_account_organization_name>.<primary_account_name>.${snowflake_standard_database.primary.name}"
  comment       = "A secondary database"

  data_retention_time_in_days                   = 10
  max_data_extension_time_in_days               = 20
  external_volume                               = "<external_volume_name>"
  catalog                                       = "<catalog_name>"
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
}

# The snowflake_secondary_database resource doesn't refresh itself, as the best practice is to use tasks scheduled for a certain interval.
# To create the refresh tasks, use separate database and schema.

resource "snowflake_standard_database" "tasks" {
  name = "database_for_tasks"
}

resource "snowflake_schema" "tasks" {
  name     = "schema_for_tasks"
  database = snowflake_standard_database.tasks.name
}

resource "snowflake_task" "refresh_secondary_database" {
  database      = snowflake_standard_database.tasks.name
  name          = "refresh_secondary_database"
  schema        = snowflake_schema.tasks.name
  schedule      = "10 minute"
  sql_statement = "ALTER DATABASE ${snowflake_secondary_database.test.name} REFRESH"
}
