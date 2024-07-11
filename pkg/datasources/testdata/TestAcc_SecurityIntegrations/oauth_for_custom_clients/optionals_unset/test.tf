resource "snowflake_oauth_integration_for_custom_clients" "test" {
  blocked_roles_list               = var.blocked_roles_list
  comment                          = var.comment
  enabled                          = var.enabled
  name                             = var.name
  network_policy                   = var.network_policy
  oauth_allow_non_tls_redirect_uri = var.oauth_allow_non_tls_redirect_uri
  oauth_client_rsa_public_key      = var.oauth_client_rsa_public_key
  oauth_client_rsa_public_key_2    = var.oauth_client_rsa_public_key_2
  oauth_client_type                = var.oauth_client_type
  oauth_enforce_pkce               = var.oauth_enforce_pkce
  oauth_issue_refresh_tokens       = var.oauth_issue_refresh_tokens
  oauth_redirect_uri               = var.oauth_redirect_uri
  oauth_refresh_token_validity     = var.oauth_refresh_token_validity
  oauth_use_secondary_roles        = var.oauth_use_secondary_roles
  pre_authorized_roles_list        = var.pre_authorized_roles_list
}


data "snowflake_security_integrations" "test" {
  depends_on = [snowflake_oauth_integration_for_custom_clients.test]

  with_describe = false
  like          = var.name
}
