# Test setup.
module "id" {
  source = "../id"
}

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
  test_prefix = format("PERFORMANCE_TESTS_%s", module.id.test_id)
}

resource "snowflake_resource_monitor" "monitor" {
  count = var.resource_count > 0 ? 1 : 0
  name  = local.test_prefix
}

# Resource with required fields
resource "snowflake_warehouse" "basic" {
  for_each = local.id_number_list
  name     = format("%s_BASIC_%v", local.test_prefix, each.key)
}

# Resource with all fields
resource "snowflake_warehouse" "complete" {
  for_each                            = local.id_number_list
  name                                = format("%s_COMPLETE_%v", local.test_prefix, each.key)
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
