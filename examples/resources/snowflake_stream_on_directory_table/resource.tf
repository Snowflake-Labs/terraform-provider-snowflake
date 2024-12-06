# basic resource
resource "snowflake_stream_on_directory_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  stage = snowflake_stage.example.fully_qualified_name
}


# resource with more fields set
resource "snowflake_stream_on_directory_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  copy_grants = true
  stage       = snowflake_stage.example.fully_qualified_name

  comment = "A stream."
}
