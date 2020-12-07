---
page_title: "snowflake_stage_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_stage_grant`





## Schema

### Required

- **database_name** (String, Required) The name of the database containing the current stage on which to grant privileges.
- **schema_name** (String, Required) The name of the schema containing the current stage on which to grant privileges.
- **stage_name** (String, Required) The name of the stage on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **privilege** (String, Optional) The privilege to grant on the stage.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **shares** (Set of String, Optional) Grants privilege to these shares.
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is stage name | privilege | true/false for with_grant_option
terraform import snowflake_stage_grant.example 'stageName|USAGE|true'
```
