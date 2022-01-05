resource "snowflake_saml_integration" "saml_integration" {
  name = "saml_integration"
  saml2_provider = "CUSTOM"
  saml2_issuer = "test_issuer"
  saml2_sso_url = "https://testsamlissuer.com"
  saml2_x509_cert = "MIICdummybase64certificate"
  enabled = true
}