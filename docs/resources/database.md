---
page_title: "snowflake_database Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_database`





## Schema

### Required

- **name** (String, Required)

### Optional

- **comment** (String, Optional)
- **data_retention_time_in_days** (Number, Optional)
- **from_database** (String, Optional) Specify a database to create a clone from.
- **from_share** (Map of String, Optional) Specify a provider and a share in this map to create a database from a share.
- **id** (String, Optional) The ID of this resource.


