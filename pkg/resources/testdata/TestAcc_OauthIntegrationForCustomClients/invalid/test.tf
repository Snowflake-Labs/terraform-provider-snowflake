resource "snowflake_oauth_integration_for_custom_clients" "test" {
  name                      = var.name
  oauth_client_type         = var.oauth_client_type
  oauth_redirect_uri        = var.oauth_redirect_uri
  oauth_use_secondary_roles = var.oauth_use_secondary_roles
  blocked_roles_list        = var.blocked_roles_list
}
