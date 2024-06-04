resource "snowflake_standard_database" "primary" {
  name         = "database_name"
  is_transient = false
  comment      = "my standard database"

  data_retention_time_in_days {
    value = 10
  }
  max_data_extension_time_in_days {
    value = 20
  }
  external_volume {
    value = "<external_volume_name>"
  }
  catalog {
    value = "<external_volume_name>"
  }
  replace_invalid_characters {
    value = false
  }
  default_ddl_collation {
    value = "en_US"
  }
  storage_serialization_policy {
    value = "COMPATIBLE"
  }
  log_level {
    value = "INFO"
  }
  trace_level {
    value = "ALWAYS"
  }

  replication {
    enable_for_account {
      account_identifier = "<secondary_account_organization_name>.<secondary_account_name>"
      with_failover      = true
    }
    ignore_edition_check = true
  }
}
