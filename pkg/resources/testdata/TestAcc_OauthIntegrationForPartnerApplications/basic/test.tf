resource "snowflake_oauth_integration_for_partner_applications" "test" {
  name         = var.name
  oauth_client = var.oauth_client
}
