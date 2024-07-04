# basic resource
resource "snowflake_oauth_integration_for_partner_applications" "test" {
  name               = "example"
  oauth_client       = "LOOKER"
  oauth_redirect_uri = "http://example.com"
  blocked_roles_list = ["ACCOUNTADMIN", "SECURITYADMIN"]
}

# resource with all fields set
resource "snowflake_oauth_integration_for_partner_applications" "test" {
  name                         = "example"
  oauth_client                 = "TABLEAU_DESKTOP"
  enabled                      = "true"
  oauth_issue_refresh_tokens   = "true"
  oauth_refresh_token_validity = 3600
  oauth_use_secondary_roles    = "IMPLICIT"
  blocked_roles_list           = ["ACCOUNTADMIN", "SECURITYADMIN", "role_id1", "role_id2"]
  comment                      = "example oauth integration for partner applications"
}
