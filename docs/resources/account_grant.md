---
page_title: "snowflake_account_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_account_grant`





## Schema

### Optional

- **id** (String, Optional) The ID of this resource.
- **privilege** (String, Optional) The privilege to grant on the schema.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is account name | privilege | true/false for with_grant_option
terraform import snowflake_account_grant.example 'accountName|USAGE|true'
```
