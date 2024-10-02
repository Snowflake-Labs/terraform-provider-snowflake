# basic resource
resource "snowflake_secret_with_authorization_code_grant" "test" {
  name                            = "EXAMPLE_SECRET"
  database                        = "EXAMPLE_DB"
  schema                          = "EXAMPLE_SCHEMA"
  api_authentication              = "EXAMPLE_SECURITY_INTEGRATION_NAME"
  oauth_refresh_token             = "EXAMPLE_TOKEN"
  oauth_refresh_token_expiry_time = "2025-01-02 15:04:01"
  comment                         = "EXAMPLE_COMMENT"
}
