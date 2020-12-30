---
page_title: "snowflake_masking_policy Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_masking_policy`





## Schema

### Required

- **database** (String, Required) The database in which to create the masking policy.
- **masking_expression** (String, Required) Specifies the SQL expression that transforms the data.
- **name** (String, Required) Specifies the identifier for the masking policy; must be unique for the database and schema in which the masking policy is created.
- **return_data_type** (String, Required) Specifies the data type to return.
- **schema** (String, Required) The schema in which to create the masking policy.
- **value_data_type** (String, Required) Specifies the data type to mask.

### Optional

- **comment** (String, Optional) Specifies a comment for the masking policy.
- **id** (String, Optional) The ID of this resource.


