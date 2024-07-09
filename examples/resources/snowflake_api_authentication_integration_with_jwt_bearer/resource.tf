# basic resource
resource "snowflake_api_authentication_integration_with_jwt_bearer" "test" {
  enabled                = true
  name                   = "test"
  oauth_client_id        = "sn-oauth-134o9erqfedlc"
  oauth_client_secret    = "eb9vaXsrcEvrFdfcvCaoijhilj4fc"
  oauth_assertion_issuer = "issuer"
}
# resource with all fields set
resource "snowflake_api_authentication_integration_with_jwt_bearer" "test" {
  comment                      = "comment"
  enabled                      = true
  name                         = "test"
  oauth_access_token_validity  = 42
  oauth_authorization_endpoint = "https://example.com"
  oauth_client_auth_method     = "CLIENT_SECRET_POST"
  oauth_client_id              = "sn-oauth-134o9erqfedlc"
  oauth_client_secret          = "eb9vaXsrcEvrFdfcvCaoijhilj4fc"
  oauth_refresh_token_validity = 42
  oauth_token_endpoint         = "https://example.com"
  oauth_assertion_issuer       = "issuer"
}
