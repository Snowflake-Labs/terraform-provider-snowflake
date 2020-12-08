---
page_title: "snowflake_system_get_aws_sns_iam_policy Data Source - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Data Source `snowflake_system_get_aws_sns_iam_policy`





## Schema

### Required

- **aws_sns_topic_arn** (String, Required) Amazon Resource Name (ARN) of the SNS topic for your S3 bucket

### Optional

- **id** (String, Optional) The ID of this resource.

### Read-only

- **aws_sns_topic_policy_json** (String, Read-only) IAM policy for Snowflake’s SQS queue to subscribe to this topic


