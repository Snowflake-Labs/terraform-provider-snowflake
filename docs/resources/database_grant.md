---
page_title: "snowflake_database_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_database_grant`



## Example Usage

```terraform
resource snowflake_database_grant grant {
  database_name = "db"

  privilege = "USAGE"
  roles     = ["role1", "role2"]
  shares    = ["share1", "share2"]

  with_grant_option = false
}
```

## Schema

### Required

- **database_name** (String, Required) The name of the database on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **privilege** (String, Optional) The privilege to grant on the database.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **shares** (Set of String, Optional) Grants privilege to these shares.
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is database name | | | privilege | true/false for with_grant_option
terraform import snowflake_database_grant.example 'databaseName|||USAGE|false'
```
