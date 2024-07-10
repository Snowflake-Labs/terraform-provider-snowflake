resource "snowflake_scim_integration" "test" {
  name        = var.name_1
  scim_client = var.scim_client
  run_as_role = var.run_as_role
  enabled     = var.enabled
}

resource "snowflake_saml2_integration" "test" {
  name            = var.name_2
  saml2_issuer    = var.saml2_issuer
  saml2_sso_url   = var.saml2_sso_url
  saml2_provider  = var.saml2_provider
  saml2_x509_cert = var.saml2_x509_cert
}

data "snowflake_security_integrations" "test" {
  depends_on = [snowflake_scim_integration.test, snowflake_saml2_integration.test]

  like = var.like
}
