---
page_title: "snowflake_stream Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_stream`





## Schema

### Required

- **database** (String, Required) The database in which to create the stream.
- **name** (String, Required) Specifies the identifier for the stream; must be unique for the database and schema in which the stream is created.
- **schema** (String, Required) The schema in which to create the stream.

### Optional

- **append_only** (Boolean, Optional) Type of the stream that will be created.
- **comment** (String, Optional) Specifies a comment for the stream.
- **id** (String, Optional) The ID of this resource.
- **on_table** (String, Optional) Name of the table the stream will monitor.

### Read-only

- **owner** (String, Read-only) Name of the role that owns the stream.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | on table name
terraform import snowflake_stream.example 'dbName|schemaName|tableName'
```
