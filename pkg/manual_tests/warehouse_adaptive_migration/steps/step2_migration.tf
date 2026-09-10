terraform {
  # removed blocks require 1.7 or higher (import blocks require 1.5)
  required_version = ">= 1.7"

  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

# 1. Forget the standard warehouse WITHOUT dropping it in Snowflake.
#    destroy = false is what makes this safe - omit it and Terraform drops the warehouse.
removed {
  from = snowflake_warehouse.example

  lifecycle {
    destroy = false
  }
}

# 2. The same warehouse, now described by the adaptive resource.
#    warehouse_size is intentionally absent - it does not apply to adaptive warehouses.
resource "snowflake_warehouse_adaptive" "example" {
  name                        = "ADAPTIVE_MIGRATION_TEST_WH"
  max_query_performance_level = "LARGE"
  comment                     = "created as a standard warehouse"
}

# 3. Adopt the existing warehouse into the adaptive resource.
import {
  to = snowflake_warehouse_adaptive.example
  id = "ADAPTIVE_MIGRATION_TEST_WH"
}
