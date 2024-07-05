# basic resource
resource "snowflake_oauth_integration_for_custom_clients" "basic" {
  name               = "saml_integration"
  oauth_client_type  = "CONFIDENTIAL"
  oauth_redirect_uri = "https://example.com"
  blocked_roles_list = ["ACCOUNTADMIN", "SECURITYADMIN"]
}

# resource with all fields set
resource "snowflake_oauth_integration_for_custom_clients" "complete" {
  name                             = "saml_integration"
  oauth_client_type                = "CONFIDENTIAL"
  oauth_redirect_uri               = "https://example.com"
  enabled                          = "true"
  oauth_allow_non_tls_redirect_uri = "true"
  oauth_enforce_pkce               = "true"
  oauth_use_secondary_roles        = "NONE"
  pre_authorized_roles_list        = ["role_id1", "role_id2"]
  blocked_roles_list               = ["ACCOUNTADMIN", "SECURITYADMIN", "role_id1", "role_id2"]
  oauth_issue_refresh_tokens       = "true"
  oauth_refresh_token_validity     = 87600
  network_policy                   = "network_policy_id"
  oauth_client_rsa_public_key      = file("rsa.pub")
  oauth_client_rsa_public_key_2    = file("rsa2.pub")
  comment                          = "my oauth integration"
}
