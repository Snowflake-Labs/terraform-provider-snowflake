resource "snowflake_oauth_integration_for_partner_applications" "test" {
  name                         = var.name
  oauth_client                 = var.oauth_client
  blocked_roles_list           = var.blocked_roles_list
  enabled                      = var.enabled
  oauth_issue_refresh_tokens   = var.oauth_issue_refresh_tokens
  oauth_refresh_token_validity = var.oauth_refresh_token_validity
  oauth_use_secondary_roles    = var.oauth_use_secondary_roles
  comment                      = var.comment
}

data "snowflake_security_integrations" "test" {
  depends_on = [snowflake_oauth_integration_for_partner_applications.test]

  with_describe = false
  like          = var.name
}
