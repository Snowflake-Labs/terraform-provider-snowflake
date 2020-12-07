
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
