---
page_title: "snowflake_view_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_view_grant`



## Example Usage

```terraform
resource snowflake_view_grant grant {
  database_name = "db"
  schema_name   = "schema"
  view_name     = "view"

  privilege = "select"
  roles = [
    "role1",
    "role2",
  ]

  shares = [
    "share1",
    "share2",
  ]

  on_future         = false
  with_grant_option = false
}
```

## Schema

### Required

- **database_name** (String, Required) The name of the database containing the current or future views on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true and a schema_name is provided, apply this grant on all future views in the given schema. When this is true and no schema_name is provided apply this grant on all future views in the given database. The view_name and shares fields must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the current or future view.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **schema_name** (String, Optional) The name of the schema containing the current or future views on which to grant privileges.
- **shares** (Set of String, Optional) Grants privilege to these shares (only valid if on_future is unset).
- **view_name** (String, Optional) The name of the view on which to grant privileges immediately (only valid if on_future is unset).
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | view name | privilege | true/false for with_grant_option
terraform import snowflake_view_grant.example 'dbName|schemaName|viewName|USAGE|false'
```
