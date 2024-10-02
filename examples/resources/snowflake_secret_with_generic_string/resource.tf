# basic resource
resource "snowflake_secret_with_client_credentials" "test" {
  name          = "EXAMPLE_SECRET"
  database      = "EXAMPLE_DB"
  schema        = "EXAMPLE_SCHEMA"
  secret_string = "EXAMPLE_SECRET_STRING"
  comment       = "EXAMPLE_COMMENT"
}
