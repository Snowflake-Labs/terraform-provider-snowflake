# Test setup.
variable "resource_count" {
  type = number
}

terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "= 1.0.1"
    }
  }
}

locals {
  id_number_list = {
    for index, val in range(0, var.resource_count) :
    val => tostring(val)
  }
}

resource "snowflake_resource_monitor" "monitor" {
  count = var.resource_count > 0 ? 1 : 0
  name  = "perf_resource_monitor"
}

# Resource with required fields
resource "snowflake_warehouse" "basic" {
  for_each = local.id_number_list
  name     = format("perf_basic_%v", each.key)
}

# Resource with all fields
resource "snowflake_warehouse" "complete" {
  for_each                            = local.id_number_list
  name                                = format("perf_complete_%v", each.key)
  warehouse_type                      = "SNOWPARK-OPTIMIZED"
  warehouse_size                      = "MEDIUM"
  max_cluster_count                   = 4
  min_cluster_count                   = 2
  scaling_policy                      = "ECONOMY"
  auto_suspend                        = 1200
  auto_resume                         = false
  initially_suspended                 = false
  resource_monitor                    = snowflake_resource_monitor.monitor[0].fully_qualified_name
  comment                             = "An example warehouse."
  enable_query_acceleration           = true
  query_acceleration_max_scale_factor = 4
  max_concurrency_level               = 4
  statement_queued_timeout_in_seconds = 5
  statement_timeout_in_seconds        = 86400
}
