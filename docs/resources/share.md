---
page_title: "snowflake_share Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_share`





## Schema

### Required

- **name** (String, Required) Specifies the identifier for the share; must be unique for the account in which the share is created.

### Optional

- **accounts** (List of String, Optional) A list of accounts to be added to the share.
- **comment** (String, Optional) Specifies a comment for the managed account.
- **id** (String, Optional) The ID of this resource.


