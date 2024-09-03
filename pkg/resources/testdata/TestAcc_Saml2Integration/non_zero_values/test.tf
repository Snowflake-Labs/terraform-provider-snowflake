resource "snowflake_saml2_integration" "test" {
  name            = var.name
  saml2_issuer    = var.saml2_issuer
  saml2_sso_url   = var.saml2_sso_url
  saml2_provider  = var.saml2_provider
  saml2_x509_cert = var.saml2_x509_cert

  enabled                        = true
  saml2_force_authn              = true
  saml2_post_logout_redirect_url = "http://example.com"
}
