resource "snowflake_oauth_integration_for_partner_applications" "test" {
  blocked_roles_list           = var.blocked_roles_list
  comment                      = var.comment
  enabled                      = var.enabled
  name                         = var.name
  oauth_client                 = var.oauth_client
  oauth_issue_refresh_tokens   = var.oauth_issue_refresh_tokens
  oauth_refresh_token_validity = var.oauth_refresh_token_validity
  oauth_use_secondary_roles    = var.oauth_use_secondary_roles
}
