resource "snowflake_oauth_integration_for_partner_applications" "test" {
  name                         = var.name
  oauth_client                 = var.oauth_client
  oauth_redirect_uri           = var.oauth_redirect_uri
  blocked_roles_list           = var.blocked_roles_list
  enabled                      = var.enabled
  oauth_issue_refresh_tokens   = var.oauth_issue_refresh_tokens
  oauth_refresh_token_validity = var.oauth_refresh_token_validity
  oauth_use_secondary_roles    = var.oauth_use_secondary_roles
  comment                      = var.comment
}
