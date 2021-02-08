---
page_title: "snowflake_warehouse_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_warehouse_grant`



## Example Usage

```terraform
resource snowflake_warehouse_grant grant {
  warehouse_name = "wh"
  privilege      = "MODIFY"

  roles = [
    "role1",
  ]

  with_grant_option = false
}
```

## Schema

### Required

- **warehouse_name** (String, Required) The name of the warehouse on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **privilege** (String, Optional) The privilege to grant on the warehouse.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is warehouse name | | | privilege | true/false for with_grant_option
terraform import snowflake_warehouse_grant.example 'warehouseName|||MODIFY|true'
```
