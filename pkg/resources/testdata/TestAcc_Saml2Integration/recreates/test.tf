resource "snowflake_saml2_integration" "test" {
  name                                = var.name
  saml2_issuer                        = var.saml2_issuer
  saml2_sso_url                       = var.saml2_sso_url
  saml2_provider                      = var.saml2_provider
  saml2_x509_cert                     = var.saml2_x509_cert
  saml2_sp_initiated_login_page_label = var.saml2_sp_initiated_login_page_label
  saml2_snowflake_issuer_url          = var.saml2_snowflake_issuer_url
  saml2_snowflake_acs_url             = var.saml2_snowflake_acs_url
  allowed_email_patterns              = var.allowed_email_patterns
  allowed_user_domains                = var.allowed_user_domains
}
