resource "snowflake_oauth_integration_for_partner_applications" "test" {
  name               = var.name
  oauth_client = var.oauth_client
  oauth_use_secondary_roles= var.oauth_use_secondary_roles
  blocked_roles_list = var.blocked_roles_list
}
