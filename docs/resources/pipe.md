---
page_title: "snowflake_pipe Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_pipe`





## Schema

### Required

- **copy_statement** (String, Required) Specifies the copy statement for the pipe.
- **database** (String, Required) The database in which to create the pipe.
- **name** (String, Required) Specifies the identifier for the pipe; must be unique for the database and schema in which the pipe is created.
- **schema** (String, Required) The schema in which to create the pipe.

### Optional

- **auto_ingest** (Boolean, Optional) Specifies a auto_ingest param for the pipe.
- **aws_sns_topic_arn** (String, Optional) Specifies the Amazon Resource Name (ARN) for the SNS topic for your S3 bucket.
- **comment** (String, Optional) Specifies a comment for the pipe.
- **id** (String, Optional) The ID of this resource.

### Read-only

- **notification_channel** (String, Read-only) Amazon Resource Name of the Amazon SQS queue for the stage named in the DEFINITION column.
- **owner** (String, Read-only) Name of the role that owns the pipe.


