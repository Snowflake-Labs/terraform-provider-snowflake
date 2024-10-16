resource "snowflake_stage" "example_stage" {
  name        = "EXAMPLE_STAGE"
  url         = "s3://com.example.bucket/prefix"
  database    = "EXAMPLE_DB"
  schema      = "EXAMPLE_SCHEMA"
  credentials = "AWS_KEY_ID='${var.example_aws_key_id}' AWS_SECRET_KEY='${var.example_aws_secret_key}'"
}

# basic resource
resource "snowflake_stream_on_directory_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  stage = snowflake_stage.stage.fully_qualified_name
}


# resource with more fields set
resource "snowflake_stream_on_directory_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  copy_grants = true
  stage       = snowflake_stage.stage.fully_qualified_name

  at {
    statement = "8e5d0ca9-005e-44e6-b858-a8f5b37c5726"
  }

  comment = "A stream."
}
