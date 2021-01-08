---
page_title: "snowflake_integration_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_integration_grant`



## Example Usage

```terraform
resource snowflake_integration_grant grant {
  integration_name = "integration"

  privilege = "USAGE"
  roles     = ["role1", "role2"]

  with_grant_option = false
}
```

## Schema

### Required

- **integration_name** (String, Required) Identifier for the integration; must be unique for your account.

### Optional

- **id** (String, Optional) The ID of this resource.
- **privilege** (String, Optional) The privilege to grant on the integration.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is integration name | privilege | true/false for with_grant_option
terraform import snowflake_integration_grant.example 'intName|USAGE|true'
```
