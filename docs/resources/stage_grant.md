---
page_title: "snowflake_stage_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_stage_grant`



## Example Usage

```terraform
resource snowflake_stage_grant grant {
  database_name = "db"
  schema_name   = "schema"
  stage_name    = "stage"

  privilege = "USAGE"

  roles  = ["role1", "role2"]
  shares = ["share1", "share2"]

  on_future         = false
  with_grant_option = false
}
```

## Schema

### Required

- **database_name** (String, Required) The name of the database containing the current stage on which to grant privileges.
- **schema_name** (String, Required) The name of the schema containing the current stage on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true and a schema_name is provided, apply this grant on all future stages in the given schema. When this is true and no schema_name is provided apply this grant on all future stages in the given database. The stage_name and shares fields must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the stage.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **shares** (Set of String, Optional) Grants privilege to these shares (only valid if on_future is false).
- **stage_name** (String, Optional) The name of the stage on which to grant privilege (only valid if on_future is false).
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | stage name | privilege | true/false for with_grant_option
terraform import snowflake_stage_grant.example 'databaseName|schemaName|stageName|USAGE|true'
```
