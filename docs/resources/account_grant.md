---
page_title: "snowflake_account_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_account_grant`



## Example Usage

```terraform
resource snowflake_account_grant grant {
  roles             = ["role1", "role2"]
  privilege         = "CREATE ROLE"
  with_grant_option = false
}
```

## Schema

### Optional

- **id** (String, Optional) The ID of this resource.
- **privilege** (String, Optional) The privilege to grant on the account.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is account name | privilege | true/false for with_grant_option
terraform import snowflake_account_grant.example 'accountName|USAGE|true'
```
