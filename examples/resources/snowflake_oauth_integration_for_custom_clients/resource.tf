resource "snowflake_oauth_integration_for_custom_clients" "tableau_desktop" {
  name                         = "TABLEAU_DESKTOP"
  oauth_client                 = "TABLEAU_DESKTOP"
  enabled                      = true
  oauth_issue_refresh_tokens   = true
  oauth_refresh_token_validity = 3600
  blocked_roles_list           = ["SYSADMIN"]
}
