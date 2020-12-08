---
page_title: "snowflake_stream_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_stream_grant`





## Schema

### Required

- **database_name** (String, Required) The name of the database containing the current or future streams on which to grant privileges.
- **schema_name** (String, Required) The name of the schema containing the current or future streams on which to grant privileges.

### Optional

- **id** (String, Optional) The ID of this resource.
- **on_future** (Boolean, Optional) When this is set to true and a schema_name is provided, apply this grant on all future streams in the given schema. When this is true and no schema_name is provided apply this grant on all future streams in the given database. The stream_name field must be unset in order to use on_future.
- **privilege** (String, Optional) The privilege to grant on the current or future stream.
- **roles** (Set of String, Optional) Grants privilege to these roles.
- **stream_name** (String, Optional) The name of the stream on which to grant privileges immediately (only valid if on_future is false).
- **with_grant_option** (Boolean, Optional) When this is set to true, allows the recipient role to grant the privileges to other roles.


