---
page_title: "snowflake_stage Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_stage`



## Example Usage

```terraform
resource "snowflake_stage" "example_stage" {
  name        = "EXAMPLE_STAGE"
  url         = "s3://com.example.bucket/prefix"
  database    = "EXAMPLE_DB"
  schema      = "EXAMPLE_SCHEMA"
  credentials = "AWS_KEY_ID='${var.example_aws_key_id}' AWS_SECRET_KEY='${var.example_aws_secret_key}'"
}

resource "snowflake_stage_grant" "grant_example_stage" {
  database_name = snowflake_stage.example_stage.database
  schema_name   = snowflake_stage.example_stage.schema
  roles         = ["LOADER"]
  privilege     = "OWNERSHIP"
  stage_name    = snowflake_stage.example_stage.name
}
```

## Schema

### Required

- **database** (String, Required) The database in which to create the stage.
- **name** (String, Required) Specifies the identifier for the stage; must be unique for the database and schema in which the stage is created.
- **schema** (String, Required) The schema in which to create the stage.

### Optional

- **aws_external_id** (String, Optional)
- **comment** (String, Optional) Specifies a comment for the stage.
- **copy_options** (String, Optional) Specifies the copy options for the stage.
- **credentials** (String, Optional) Specifies the credentials for the stage.
- **encryption** (String, Optional) Specifies the encryption settings for the stage.
- **file_format** (String, Optional) Specifies the file format for the stage.
- **id** (String, Optional) The ID of this resource.
- **snowflake_iam_user** (String, Optional)
- **storage_integration** (String, Optional) Specifies the name of the storage integration used to delegate authentication responsibility for external cloud storage to a Snowflake identity and access management (IAM) entity.
- **url** (String, Optional) Specifies the URL for the stage.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | stage name
terraform import snowflake_stage.example 'dbName|schemaName|stageName'
```
