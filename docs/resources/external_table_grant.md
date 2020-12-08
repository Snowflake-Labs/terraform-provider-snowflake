---
page_title: "snowflake_external_table_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_external_table_grant`



## Example Usage

```terraform
resource snowflake_external_table_grant grant {
  database_name       = "db"
  schema_name         = "schema"
  external_table_name = "external_table"

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

- **database_name** (String, Required) The name of the database containing the current or future external tables on which to grant privileges.
- **schema_name** (String, Required) The name of the schema containing the current or future external tables on which to grant privileges.

### Optional

- **external_table_name** (String, Optional) The name of the external table on which to grant privileges immediately (only valid if on_future is false).
- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true and a schema_name is provided, apply this grant on all future external tables in the given schema. When this is true and no schema_name is provided apply this grant on all future external tables in the given database. The external_table_name and shares fields must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the current or future external table.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **shares** (Set of String, Optional) Grants privilege to these shares (only valid if on_future is false).
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | external table name | privilege | true/false for with_grant_option
terraform import snowflake_external_table_grant.example 'dbName|schemaName|externalTableName|SELECT|false'
```
