resource "snowflake_external_oauth_integration" "test" {
  name                                            = var.name
  type                                            = var.type
  enabled                                         = var.enabled
  external_oauth_issuer                           = var.external_oauth_issuer
  external_oauth_scope_mapping_attribute          = var.external_oauth_scope_mapping_attribute
  external_oauth_snowflake_user_mapping_attribute = var.external_oauth_snowflake_user_mapping_attribute
}
