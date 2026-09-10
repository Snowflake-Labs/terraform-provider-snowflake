terraform {
  # removed blocks require 1.7 or higher (import blocks require 1.5)
  required_version = ">= 1.7"

  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

# Steady state: the removed and import blocks are gone. A plan here must be empty,
# which proves the standard resource is really out of the state and the adaptive one is stable.
resource "snowflake_warehouse_adaptive" "example" {
  name                        = "ADAPTIVE_MIGRATION_TEST_WH"
  max_query_performance_level = "LARGE"
  comment                     = "created as a standard warehouse"
}
