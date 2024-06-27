# basic resource
resource "snowflake_api_authentication_integration_with_client_credentials" "test" {
  enabled             = true
  name                = "foo"
  oauth_client_id     = "foo"
  oauth_client_secret = "foo"
}
# resource with all fields set
resource "snowflake_api_authentication_integration_with_client_credentials" "test" {
  comment                     = "foo"
  enabled                     = true
  name                        = "foo"
  oauth_access_token_validity = 42
  oauth_allowed_scopes        = ["foo"]
  oauth_client_auth_method    = "CLIENT_SECRET_POST"
  oauth_client_id             = "foo"
  oauth_client_secret         = "foo"
  oauth_token_endpoint        = "https://example.com"
}
