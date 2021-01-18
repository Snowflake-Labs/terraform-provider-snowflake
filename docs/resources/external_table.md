---
page_title: "snowflake_external_table Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_external_table`



## Example Usage

```terraform
resource snowflake_external_table external_table {
  database = "db"
  schema   = "schema"
  name     = "external_table"
  comment  = "External table"

  column {
    name = "id"
    type = "int"
  }

  column {
    name = "data"
    type = "text"
  }
}
```

## Schema

### Required

- **column** (Block List, Min: 1) Definitions of a column to create in the external table. Minimum one required. (see [below for nested schema](#nestedblock--column))
- **database** (String, Required) The database in which to create the external table.
- **file_format** (String, Required) Specifies the file format for the external table.
- **location** (String, Required) Specifies a location for the external table.
- **name** (String, Required) Specifies the identifier for the external table; must be unique for the database and schema in which the externalTable is created.
- **schema** (String, Required) The schema in which to create the external table.

### Optional

- **auto_refresh** (Boolean, Optional) Specifies whether to automatically refresh the external table metadata once, immediately after the external table is created.
- **aws_sns_topic** (String, Optional) Specifies the aws sns topic for the external table.
- **comment** (String, Optional) Specifies a comment for the external table.
- **copy_grants** (Boolean, Optional) Specifies to retain the access permissions from the original table when an external table is recreated using the CREATE OR REPLACE TABLE variant
- **id** (String, Optional) The ID of this resource.
- **partition_by** (List of String, Optional) Specifies any partition columns to evaluate for the external table.
- **refresh_on_create** (Boolean, Optional) Specifies weather to refresh when an external table is created.

### Read-only

- **owner** (String, Read-only) Name of the role that owns the external table.

<a id="nestedblock--column"></a>
### Nested Schema for `column`

Required:

- **as** (String, Required) String that specifies the expression for the column. When queried, the column returns results derived from this expression.
- **name** (String, Required) Column name
- **type** (String, Required) Column type, e.g. VARIANT

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | external table name
terraform import snowflake_external_table.example 'dbName|schemaName|externalTableName'
```
