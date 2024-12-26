
resource "snowflake_stage" "example_stage" {
  name        = "EXAMPLE_STAGE"
  url         = "s3://com.example.bucket/prefix"
  database    = "EXAMPLE_DB"
  schema      = "EXAMPLE_SCHEMA"
  credentials = "AWS_KEY_ID='${var.example_aws_key_id}' AWS_SECRET_KEY='${var.example_aws_secret_key}'"
}
/*
  * Examples of usage for `file_format`:
    * with hardcoding value: `file_format="FORMAT_NAME = DB.SCHEMA.FORMATNAME"`
    * from dynamic value: `file_format = "FORMAT_NAME = ${snowflake_database.mydb.name}.${snowflake_schema.myschema.name}.${snowflake_file_format.myfileformat.name}"`
    * from expression: `file_format = format("FORMAT_NAME =%s.%s.MYFILEFORMAT", var.db_name, each.value.schema_name)`
  * Reference: [#265](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/265)
*/