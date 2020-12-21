---
page_title: "snowflake_resource_monitor_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_resource_monitor_grant`



## Example Usage

```terraform
resource snowflake_monitor_grant grant {
  monitor_name      = "monitor"
  privilege         = "MODIFY"
  roles             = ["role1"]
  with_grant_option = false
}
```

## Schema

### Required

- **monitor_name** (String, Required) Identifier for the resource monitor; must be unique for your account.

### Optional

- **id** (String, Optional) The ID of this resource.
- **privilege** (String, Optional) The privilege to grant on the resource monitor.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_resource_monitor_grant.example name
```
