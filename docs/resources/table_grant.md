---
page_title: "snowflake_table_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_table_grant`



## Example Usage

```terraform
resource snowflake_table_grant grant {
  database_name = "database"
  schema_name   = "schema"
  table_name    = "table"

  privilege = "SELECT"
  roles     = ["role1"]
  shares    = ["share1"]

  on_future         = false
  with_grant_option = false
}
```

## Schema

### Required

- **database_name** (String, Required) The name of the database containing the current or future tables on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true and a schema_name is provided, apply this grant on all future tables in the given schema. When this is true and no schema_name is provided apply this grant on all future tables in the given database. The table_name and shares fields must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the current or future table.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **schema_name** (String, Optional) The name of the schema containing the current or future tables on which to grant privileges.
- **shares** (Set of String, Optional) Grants privilege to these shares (only valid if on_future is unset).
- **table_name** (String, Optional) The name of the table on which to grant privileges immediately (only valid if on_future is unset).
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | table name | privilege | true/false for with_grant_option
terraform import snowflake_table_grant.example 'databaseName|schemaName|tableName|MODIFY|true'
```
