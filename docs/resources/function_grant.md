---
page_title: "snowflake_function_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_function_grant`



## Example Usage

```terraform
resource snowflake_function_grant grant {
  database_name   = "db"
  schema_name     = "schema"
  function_name  = "function"

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

- **database_name** (String, Required) The name of the database containing the current or future functions on which to grant privileges.
- **schema_name** (String, Required) The name of the schema containing the current or future functions on which to grant privileges.

### Optional

- **arguments** (Block List) List of the arguments for the function (must be present if function_name is present) (see [below for nested schema](#nestedblock--arguments))
- **function_name** (String, Optional) The name of the function on which to grant privileges immediately (only valid if on_future is false).
- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true and a schema_name is provided, apply this grant on all future functions in the given schema. When this is true and no schema_name is provided apply this grant on all future functions in the given database. The function_name, arguments, return_type, and shares fields must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the current or future function.
- **return_type** (String, Optional) The return type of the function (must be present if function_name is present)
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
# format is database name | schema name | function signature | privilege | true/false for with_grant_option
terraform import snowflake_function_grant.example 'dbName|schemaName|functionName(ARG1 ARG1TYPE, ARG2 ARG2TYPE):RETURNTYPE|USAGE|false'
```
