package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestSamlIntegration(t *testing.T) {
	r := require.New(t)
	builder := snowflake.SamlIntegration("test_saml_integration")
	r.NotNil(builder)

	q := builder.Show()
	r.Equal("SHOW SECURITY INTEGRATIONS LIKE 'test_saml_integration'", q)

	q = builder.Describe()
	r.Equal("DESCRIBE SECURITY INTEGRATION \"test_saml_integration\"", q)

	c := builder.Create()
	c.SetRaw(`TYPE=SAML2`)
	c.SetString(`saml2_issuer`, "test_issuer")
	c.SetString(`saml2_sso_url`, "https://testsamlissuer.com")
	c.SetString(`saml2_provider`, "CUSTOM")
	c.SetString(`saml2_x509_cert`, "MIICdummybase64certificate")
	c.SetBool(`enabled`, true)
	q = c.Statement()
	r.Equal(`CREATE SECURITY INTEGRATION "test_saml_integration" TYPE=SAML2 SAML2_ISSUER='test_issuer' SAML2_PROVIDER='CUSTOM' SAML2_SSO_URL='https://testsamlissuer.com' SAML2_X509_CERT='MIICdummybase64certificate' ENABLED=true`, q)

	d := builder.Alter()
	d.SetRaw(`TYPE=SAML2`)
	d.SetString(`saml2_issuer`, "test_issuer")
	d.SetString(`saml2_sso_url`, "https://testsamlissuer.com")
	d.SetString(`saml2_provider`, "CUSTOM")
	d.SetString(`saml2_x509_cert`, "MIICdummybase64certificate")
	d.SetBool(`enabled`, false)
	q = d.Statement()
	r.Equal(`ALTER SECURITY INTEGRATION "test_saml_integration" SET TYPE=SAML2 SAML2_ISSUER='test_issuer' SAML2_PROVIDER='CUSTOM' SAML2_SSO_URL='https://testsamlissuer.com' SAML2_X509_CERT='MIICdummybase64certificate' ENABLED=false`, q)

	e := builder.Drop()
	r.Equal(`DROP SECURITY INTEGRATION "test_saml_integration"`, e)
}
