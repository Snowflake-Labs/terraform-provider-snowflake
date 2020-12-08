---
page_title: "snowflake_warehouse Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_warehouse`



## Example Usage

```terraform
resource snowflake_warehouse w {
  name           = "test"
  comment        = "foo"
  warehouse_size = "small"
}
```

## Schema

### Required

- **name** (String, Required)

### Optional

- **auto_resume** (Boolean, Optional) Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it.
- **auto_suspend** (Number, Optional) Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.
- **comment** (String, Optional)
- **id** (String, Optional) The ID of this resource.
- **initially_suspended** (Boolean, Optional) Specifies whether the warehouse is created initially in the ‘Suspended’ state.
- **max_cluster_count** (Number, Optional) Specifies the maximum number of server clusters for the warehouse.
- **min_cluster_count** (Number, Optional) Specifies the minimum number of server clusters for the warehouse (only applies to multi-cluster warehouses).
- **resource_monitor** (String, Optional) Specifies the name of a resource monitor that is explicitly assigned to the warehouse.
- **scaling_policy** (String, Optional) Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode.
- **statement_timeout_in_seconds** (Number, Optional) Specifies the time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system
- **wait_for_provisioning** (Boolean, Optional) Specifies whether the warehouse, after being resized, waits for all the servers to provision before executing any queued or new queries.
- **warehouse_size** (String, Optional)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_warehouse.example warehouseName
```
