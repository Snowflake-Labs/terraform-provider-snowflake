# basic resource
resource "snowflake_secret_with_basic_authentication" "test" {
  name     = "EXAMPLE_SECRET"
  database = "EXAMPLE_DB"
  schema   = "EXAMPLE_SCHEMA"
  username = "EXAMPLE_USERNAME"
  password = "EXAMPLE_PASSWORD"
  comment  = "EXAMPLE_COMMENT"
}
