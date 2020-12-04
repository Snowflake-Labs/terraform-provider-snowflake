---
page_title: "snowflake_task Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_task`





## Schema

### Required

- **database** (String, Required) The database in which to create the task.
- **name** (String, Required) Specifies the identifier for the task; must be unique for the database and schema in which the task is created.
- **schema** (String, Required) The schema in which to create the task.
- **sql_statement** (String, Required) Any single SQL statement, or a call to a stored procedure, executed when the task runs.
- **warehouse** (String, Required) The warehouse the task will use.

### Optional

- **after** (String, Optional) Specifies the predecessor task in the same database and schema of the current task. When a run of the predecessor task finishes successfully, it triggers this task (after a brief lag).
- **comment** (String, Optional) Specifies a comment for the task.
- **enabled** (Boolean, Optional) Specifies if the task should be started (enabled) after creation or should remain suspended (default).
- **id** (String, Optional) The ID of this resource.
- **schedule** (String, Optional) The schedule for periodically running the task. This can be a cron or interval in minutes.
- **session_parameters** (Map of String, Optional) Specifies session parameters to set for the session when the task runs. A task supports all session parameters.
- **user_task_timeout_ms** (Number, Optional) Specifies the time limit on a single run of the task before it times out (in milliseconds).
- **when** (String, Optional) Specifies a Boolean SQL expression; multiple conditions joined with AND/OR are supported.


