---
page_title: "snowflake_warehouse Resource - terraform-provider-snowflake"
subcategory: "Stable"
description: |-
  Resource used to manage warehouse objects. For more information, check warehouse documentation https://docs.snowflake.com/en/sql-reference/commands-warehouse.
---

<!-- TODO(SNOW-1844996): Remove this note.-->
-> **Note** Field `RESOURCE_CONSTRAINT` is currently missing. It will be added in the future.

<!-- TODO(SNOW-1642723): Remove or adjust this note.-->
-> **Note** Assigning resource monitors to warehouses requires ACCOUNTADMIN role. To do this, either manage the warehouse resource with ACCOUNTADMIN role, or use [execute](./execute) instead. See [this issue](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3019) for more details.

# snowflake_warehouse (Resource)

Resource used to manage warehouse objects. For more information, check [warehouse documentation](https://docs.snowflake.com/en/sql-reference/commands-warehouse).

## Example Usage

```terraform
# Resource with required fields
resource "snowflake_warehouse" "warehouse" {
  name = "WAREHOUSE"
}

# Resource with all fields
resource "snowflake_warehouse" "warehouse" {
  name                                = "WAREHOUSE"
  warehouse_type                      = "SNOWPARK-OPTIMIZED"
  warehouse_size                      = "MEDIUM"
  max_cluster_count                   = 4
  min_cluster_count                   = 2
  scaling_policy                      = "ECONOMY"
  auto_suspend                        = 1200
  auto_resume                         = false
  initially_suspended                 = false
  resource_monitor                    = snowflake_resource_monitor.monitor.fully_qualified_name
  comment                             = "An example warehouse."
  enable_query_acceleration           = true
  query_acceleration_max_scale_factor = 4

  max_concurrency_level               = 4
  statement_queued_timeout_in_seconds = 5
  statement_timeout_in_seconds        = 86400
}
```
-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

-> **Note** If a field has default value, it is displayed next to the type in the schema. If the default is computed from external sources (e.g., environment variables), it displays `Default is computed`.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_warehouse.example '"<warehouse_name>"'
```
