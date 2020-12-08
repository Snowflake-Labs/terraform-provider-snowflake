---
page_title: "snowflake_sequence_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_sequence_grant`





## Schema

### Required

- **database_name** (String, Required) The name of the database containing the current or future sequences on which to grant privileges.
- **schema_name** (String, Required) The name of the schema containing the current or future sequences on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true and a schema_name is provided, apply this grant on all future sequences in the given schema. When this is true and no schema_name is provided apply this grant on all future sequences in the given database. The sequence_name field must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the current or future sequence.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **sequence_name** (String, Optional) The name of the sequence on which to grant privileges immediately (only valid if on_future is false).
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.


