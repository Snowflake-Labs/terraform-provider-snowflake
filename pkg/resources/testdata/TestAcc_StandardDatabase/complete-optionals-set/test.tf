resource "snowflake_standard_database" "test" {
  name          = var.name
  comment = var.comment
  is_transient  = var.transient

  data_retention_time_in_days {
    value = var.data_retention_time_in_days
  }
  max_data_extension_time_in_days {
    value = var.max_data_extension_time_in_days
  }
  external_volume              {
    value= var.external_volume
  }
  catalog                      {
    value = var.catalog
  }
  replace_invalid_characters   {
    value = var.replace_invalid_characters
  }
  default_ddl_collation        {
    value = var.default_ddl_collation
  }
  storage_serialization_policy {
    value = var.storage_serialization_policy
  }
  log_level                    {
    value = var.log_level
  }
  trace_level                  {
    value = var.trace_level
  }
}
