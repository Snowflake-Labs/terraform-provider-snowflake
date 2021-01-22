---
page_title: "snowflake_schema_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_schema_grant`



## Example Usage

```terraform
resource snowflake_schema_grant grant {
  database_name = "db"
  schema_name   = "schema"

  privilege = "USAGE"
  roles     = ["role1", "role2"]
  shares    = ["share1", "share2"]

  on_future         = false
  with_grant_option = false
}
```

## Schema

### Required

- **database_name** (String, Required) The name of the database containing the schema on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true, apply this grant on all future schemas in the given database. The schema_name and shares fields must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the current or future schema. Note that if "OWNERSHIP" is specified, ensure that the role that terraform is using is granted access.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **schema_name** (String, Optional) The name of the schema on which to grant privileges.
- **shares** (Set of String, Optional) Grants privilege to these shares (only valid if on_future is unset).
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | | privilege | true/false for with_grant_option
terraform import snowflake_schema_grant.example 'databaseName|schemaName||MONITOR|false'
```
