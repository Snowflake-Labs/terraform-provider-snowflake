---
page_title: "snowflake_schema Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_schema`



## Example Usage

```terraform
resource snowflake_schema schema {
  database = "db"
  name     = "schema"
  comment  = "A schema."

  is_transient        = false
  is_managed          = false
  data_retention_days = 1
}
```

## Schema

### Required

- **database** (String, Required) The database in which to create the schema.
- **name** (String, Required) Specifies the identifier for the schema; must be unique for the database in which the schema is created.

### Optional

- **comment** (String, Optional) Specifies a comment for the schema.
- **data_retention_days** (Number, Optional) Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the schema, as well as specifying the default Time Travel retention time for all tables created in the schema.
- **id** (String, Optional) The ID of this resource.
- **is_managed** (Boolean, Optional) Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner.
- **is_transient** (Boolean, Optional) Specifies a schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.

## Import

Import is supported using the following syntax:

```shell
# format is schema name | privilege | true/false for with_grant_option
terraform import snowflake_schema.example 'schemaName|USAGE|true'
```
