package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_SamlIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SAML_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAccSamlIntegration")
	}

	samlIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: samlIntegrationConfig(samlIntName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "name", samlIntName),
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "saml2_issuer", "test_issuer"),
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "saml2_sso_url", "https://testsamlissuer.com"),
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "saml2_provider", "CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "saml2_x509_cert", "MIICYzCCAcygAwIBAgIBADANBgkqhkiG9w0BAQUFADAuMQswCQYDVQQGEwJVUzEMMAoGA1UEChMDSUJNMREwDwYDVQQLEwhMb2NhbCBDQTAeFw05OTEyMjIwNTAwMDBaFw0wMDEyMjMwNDU5NTlaMC4xCzAJBgNVBAYTAlVTMQwwCgYDVQQKEwNJQk0xETAPBgNVBAsTCExvY2FsIENBMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD2bZEo7xGaX2/0GHkrNFZvlxBou9v1Jmt/PDiTMPve8r9FeJAQ0QdvFST/0JPQYD20rH0bimdDLgNdNynmyRoS2S/IInfpmf69iyc2G0TPyRvmHIiOZbdCd+YBHQi1adkj17NDcWj6S14tVurFX73zx0sNoMS79q3tuXKrDsxeuwIDAQABo4GQMIGNMEsGCVUdDwGG+EIBDQQ+EzxHZW5lcmF0ZWQgYnkgdGhlIFNlY3VyZVdheSBTZWN1cml0eSBTZXJ2ZXIgZm9yIE9TLzM5MCAoUkFDRikwDgYDVR0PAQH/BAQDAgAGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFJ3+ocRyCTJw067dLSwr/nalx6YMMA0GCSqGSIb3DQEBBQUAA4GBAMaQzt+zaj1GU77yzlr8iiMBXgdQrwsZZWJo5exnAucJAEYQZmOfyLiMD6oYq+ZnfvM0n8G/Y79q8nhwvuxpYOnRSAXFp6xSkrIOeZtJMY1h00LKp/JX3Ng1svZ2agE126JHsQ0bhzN5TKsYfbwfTwfjdWAGy6Vf1nYi/rO+ryMO"),
					resource.TestCheckResourceAttrSet("snowflake_saml_integration.test_saml_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_saml_integration.test_saml_int", "saml2_snowflake_x509_cert"),
					resource.TestCheckResourceAttrSet("snowflake_saml_integration.test_saml_int", "saml2_snowflake_acs_url"),
					resource.TestCheckResourceAttrSet("snowflake_saml_integration.test_saml_int", "saml2_snowflake_issuer_url"),
					resource.TestCheckResourceAttrSet("snowflake_saml_integration.test_saml_int", "saml2_snowflake_metadata"),
					resource.TestCheckResourceAttrSet("snowflake_saml_integration.test_saml_int", "saml2_digest_methods_used"),
					resource.TestCheckResourceAttrSet("snowflake_saml_integration.test_saml_int", "saml2_signature_methods_used"),
				),
			},
		},
	})
}

func samlIntegrationConfig(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_saml_integration" "test_saml_int" {
		name = "%s"
		saml2_issuer = "test_issuer"
		saml2_sso_url = "https://testsamlissuer.com"
		saml2_provider = "CUSTOM"
		saml2_x509_cert = "MIICYzCCAcygAwIBAgIBADANBgkqhkiG9w0BAQUFADAuMQswCQYDVQQGEwJVUzEMMAoGA1UEChMDSUJNMREwDwYDVQQLEwhMb2NhbCBDQTAeFw05OTEyMjIwNTAwMDBaFw0wMDEyMjMwNDU5NTlaMC4xCzAJBgNVBAYTAlVTMQwwCgYDVQQKEwNJQk0xETAPBgNVBAsTCExvY2FsIENBMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD2bZEo7xGaX2/0GHkrNFZvlxBou9v1Jmt/PDiTMPve8r9FeJAQ0QdvFST/0JPQYD20rH0bimdDLgNdNynmyRoS2S/IInfpmf69iyc2G0TPyRvmHIiOZbdCd+YBHQi1adkj17NDcWj6S14tVurFX73zx0sNoMS79q3tuXKrDsxeuwIDAQABo4GQMIGNMEsGCVUdDwGG+EIBDQQ+EzxHZW5lcmF0ZWQgYnkgdGhlIFNlY3VyZVdheSBTZWN1cml0eSBTZXJ2ZXIgZm9yIE9TLzM5MCAoUkFDRikwDgYDVR0PAQH/BAQDAgAGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFJ3+ocRyCTJw067dLSwr/nalx6YMMA0GCSqGSIb3DQEBBQUAA4GBAMaQzt+zaj1GU77yzlr8iiMBXgdQrwsZZWJo5exnAucJAEYQZmOfyLiMD6oYq+ZnfvM0n8G/Y79q8nhwvuxpYOnRSAXFp6xSkrIOeZtJMY1h00LKp/JX3Ng1svZ2agE126JHsQ0bhzN5TKsYfbwfTwfjdWAGy6Vf1nYi/rO+ryMO"
		enabled = false
	}
	`, name)
}
