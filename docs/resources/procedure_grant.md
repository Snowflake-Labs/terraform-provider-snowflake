---
page_title: "snowflake_procedure_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_procedure_grant`



## Example Usage

```terraform
resource snowflake_procedure_grant grant {
  database_name   = "db"
  schema_name     = "schema"
  procedure_name  = "procedure"

  arguments   = [
    {
      "name": "a",
      "type": "array"
    },
    {
      "name": "b",
      "type": "string"
    }
  ]
  return_type = "string"

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

- **database_name** (String, Required) The name of the database containing the current or future procedures on which to grant privileges.
- **schema_name** (String, Required) The name of the schema containing the current or future procedures on which to grant privileges.

### Optional

- **arguments** (Block List) List of the arguments for the procedure (must be present if procedure_name is present) (see [below for nested schema](#nestedblock--arguments))
- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true and a schema_name is provided, apply this grant on all future procedures in the given schema. When this is true and no schema_name is provided apply this grant on all future procedures in the given database. The procedure_name and shares fields must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the current or future procedure.
- **procedure_name** (String, Optional) The name of the procedure on which to grant privileges immediately (only valid if on_future is false).
- **return_type** (String, Optional) The return type of the procedure (must be present if procedure_name is present)
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **shares** (Set of String, Optional) Grants privilege to these shares (only valid if on_future is false).
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

<a id="nestedblock--arguments"></a>
### Nested Schema for `arguments`

Required:

- **name** (String, Required) The argument name
- **type** (String, Required) The argument type

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | procedure signature | privilege | true/false for with_grant_option
terraform import snowflake_procedure_grant.example 'dbName|schemaName|procedureName(ARG1 ARG1TYPE, ARG2 ARG2TYPE):RETURNTYPE|USAGE|false'
```
