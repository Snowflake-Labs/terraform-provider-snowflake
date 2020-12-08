---
page_title: "snowflake_share Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_share`



## Example Usage

```terraform
resource snowflake_share share {
  database_name = "db"
  schema_name   = "schema"
  stage_name    = "stage"

  privilege = "USAGE"
  roles     = ["role1", "role2"]
  shares    = ["share1", "share2"]

  with_grant_option = false
}
```

## Schema

### Required

- **name** (String, Required) Specifies the identifier for the share; must be unique for the account in which the share is created.

### Optional

- **accounts** (List of String, Optional) A list of accounts to be added to the share.
- **comment** (String, Optional) Specifies a comment for the managed account.
- **id** (String, Optional) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_share.example name
```
