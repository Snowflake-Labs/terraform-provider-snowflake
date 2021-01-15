---
page_title: "snowflake_view Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_view`



## Example Usage

```terraform
resource snowflake_view view {
  database = "db"
  schema   = "schema"
  name     = "view"

  comment = "comment"

  statement  = <<-SQL
    select * from foo;
SQL
  or_replace = false
  is_secure  = false
}
```

## Schema

### Required

- **database** (String, Required) The database in which to create the view. Don't use the | character.
- **name** (String, Required) Specifies the identifier for the view; must be unique for the schema in which the view is created. Don't use the | character.
- **schema** (String, Required) The schema in which to create the view. Don't use the | character.
- **statement** (String, Required) Specifies the query used to create the view.

### Optional

- **comment** (String, Optional) Specifies a comment for the view.
- **id** (String, Optional) The ID of this resource.
- **is_secure** (Boolean, Optional) Specifies that the view is secure.
- **or_replace** (Boolean, Optional) Overwrites the View if it exists.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | view name
terraform import snowflake_view.example 'dbName|schemaName|viewName'
```
