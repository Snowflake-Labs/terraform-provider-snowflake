# basic resource
resource "snowflake_api_authentication_integration_with_client_credentials" "test" {
  enabled             = true
  name                = "test"
  oauth_client_id     = "sn-oauth-134o9erqfedlc"
  oauth_client_secret = "eb9vaXsrcEvrFdfcvCaoijhilj4fc"
}
# resource with all fields set
resource "snowflake_api_authentication_integration_with_client_credentials" "test" {
  comment                     = "comment"
  enabled                     = true
  name                        = "test"
  oauth_access_token_validity = 42
  oauth_allowed_scopes        = ["useraccount"]
  oauth_client_auth_method    = "CLIENT_SECRET_POST"
  oauth_client_id             = "sn-oauth-134o9erqfedlc"
  oauth_client_secret         = "eb9vaXsrcEvrFdfcvCaoijhilj4fc"
  oauth_token_endpoint        = "https://example.com"
}
