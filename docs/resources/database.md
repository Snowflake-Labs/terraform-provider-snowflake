---
page_title: "snowflake_database Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_database`



## Example Usage

```terraform
resource "snowflake_database" "test" {
  name                        = "testing"
  comment                     = "test comment"
  data_retention_time_in_days = 3
}

resource "snowflake_database" "test2" {
  name    = "testing_2"
  comment = "test comment 2"
}
```

## Schema

### Required

- **name** (String, Required)

### Optional

- **comment** (String, Optional)
- **data_retention_time_in_days** (Number, Optional)
- **from_database** (String, Optional) Specify a database to create a clone from.
- **from_share** (Map of String, Optional) Specify a provider and a share in this map to create a database from a share.
- **id** (String, Optional) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_database.example name
```
