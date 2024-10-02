# basic resource
resource "snowflake_secret_with_client_credentials" "test" {
  name               = "EXAMPLE_SECRET"
  database           = "EXAMPLE_DB"
  schema             = "EXAMPLE_SCHEMA"
  api_authentication = "EXAMPLE_SECURITY_INTEGRATION_NAME"
  oauth_scopes       = ["useraccount", "testscope"]
  comment            = "EXAMPLE_COMMENT"
}
