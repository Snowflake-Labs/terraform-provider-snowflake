---
page_title: "snowflake_role_grants Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_role_grants`





## Schema

### Required

- **role_name** (String, Required) The name of the role we are granting.

### Optional

- **id** (String, Optional) The ID of this resource.
- **roles** (Set of String, Optional) Grants role to this specified role.
- **users** (Set of String, Optional) Grants role to this specified user.


