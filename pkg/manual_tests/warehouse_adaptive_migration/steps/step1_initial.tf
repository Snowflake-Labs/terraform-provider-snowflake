terraform {
  # removed blocks require 1.7 or higher (import blocks require 1.5)
  required_version = ">= 1.7"

  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

# The warehouse we are going to migrate. A standard warehouse with a size set,
# so we can observe Snowflake clearing the properties that do not apply to adaptive warehouses.
resource "snowflake_warehouse" "example" {
  name           = "ADAPTIVE_MIGRATION_TEST_WH"
  warehouse_size = "MEDIUM"
  comment        = "created as a standard warehouse"
}
