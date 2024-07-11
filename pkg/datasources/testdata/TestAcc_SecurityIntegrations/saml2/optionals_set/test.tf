resource "snowflake_saml2_integration" "test" {
  allowed_email_patterns              = var.allowed_email_patterns
  allowed_user_domains                = var.allowed_user_domains
  comment                             = var.comment
  enabled                             = var.enabled
  name                                = var.name
  saml2_enable_sp_initiated           = var.saml2_enable_sp_initiated
  saml2_force_authn                   = var.saml2_force_authn
  saml2_issuer                        = var.saml2_issuer
  saml2_post_logout_redirect_url      = var.saml2_post_logout_redirect_url
  saml2_provider                      = var.saml2_provider
  saml2_requested_nameid_format       = var.saml2_requested_nameid_format
  saml2_sign_request                  = var.saml2_sign_request
  saml2_snowflake_acs_url             = var.saml2_snowflake_acs_url
  saml2_snowflake_issuer_url          = var.saml2_snowflake_issuer_url
  saml2_sp_initiated_login_page_label = var.saml2_sp_initiated_login_page_label
  saml2_sso_url                       = var.saml2_sso_url
  saml2_x509_cert                     = var.saml2_x509_cert
}

data "snowflake_security_integrations" "test" {
  depends_on = [snowflake_saml2_integration.test]

  like = var.name
}
