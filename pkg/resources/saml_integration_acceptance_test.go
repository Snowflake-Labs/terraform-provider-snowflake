package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SamlIntegration(t *testing.T) {
	// TODO [SNOW-926148]: unskip
	testenvs.SkipTestIfSet(t, testenvs.SkipSamlIntegrationTest, "because was skipped earlier")

	samlIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: samlIntegrationConfig(samlIntName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "name", samlIntName),
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "saml2_issuer", "test_issuer"),
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "saml2_sso_url", "https://samltest.id/saml/sp"),
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "saml2_provider", "CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_saml_integration.test_saml_int", "saml2_x509_cert", "MIIERTCCAq2gAwIBAgIJAKmtzjCD1+tqMA0GCSqGSIb3DQEBCwUAMDUxMzAxBgNVBAMTKmlwLTE3Mi0zMS0yOC02NC51cy13ZXN0LTIuY29tcHV0ZS5pbnRlcm5hbDAeFw0xODA4MTgyMzI0MjNaFw0yODA4MTUyMzI0MjNaMDUxMzAxBgNVBAMTKmlwLTE3Mi0zMS0yOC02NC51cy13ZXN0LTIuY29tcHV0ZS5pbnRlcm5hbDCCAaIwDQYJKoZIhvcNAQEBBQADggGPADCCAYoCggGBALhUlY3SkIOze+l8y6dBzM6p7B8OykJWlwizszU16Lih8D7KLhNJfahoVxbPxB3YFM/81PJLOeK2krvJ5zY6CJyQY3sPQAkZKI7I8qq9lmZ2g4QPqybNstXS6YUXJNUt/ixbbK/N97+LKTiSutbD1J7AoFnouMuLjlhN5VRZ43jez4xLSHVZaYuUFKn01Y9oLKbj46LQnZnJCAGpTgPqEQJr6GpVGw43bKyUpGoaPrdDRgRgtPMUWgFDkgcI3QiV1lsKfBs1t1E2UA7ACFnlJZpEuBtwgivzo3VeitiSaF3Jxh25EY5/vABpcgQQRz3RH2l8MMKdRsxb8VT3yh2S+CX55s+cN67LiCPr6f2u+KS1iKfB9mWN6o2S4lcmo82HIBbsuXJV0oA1HrGMyyc4Y9nng/I8iuAp8or1JrWRHQ+8NzO85DWK0rtvtLPxkvw0HK32glyuOP/9F05Z7+tiVIgn67buC0EdoUm1RSpibqmB1ST2PikslOlVbJuy4Ah93wIDAQABo1gwVjA1BgNVHREELjAsgippcC0xNzItMzEtMjgtNjQudXMtd2VzdC0yLmNvbXB1dGUuaW50ZXJuYWwwHQYDVR0OBBYEFAdsTxYfulJ5yunYtgYJHC9IcevzMA0GCSqGSIb3DQEBCwUAA4IBgQB3J6i7KreiHL8NPMglfWLHk1PZOgvIEEpKL+GRebvcbyqgcuc3VVPylq70VvGqhJxp1q/mzLfraUiypzfWFGm9zfwIg0H5TqRZYEPTvgIhIICjaDWRwZBDJG8D5G/KoV60DlUG0crPBlIuCCr/SRa5ZoDQqvucTfr3Rx4Ha6koXFSjoSXllR+jn4GnInhm/WH137a+v35PUcffNxfuehoGn6i4YeXF3cwJK4e35cOFW+dLbnaLk+Ty7HOGvpw86h979C6mJ9qEHYgq9rQyzlSPbLZGZSgVcIezunOaOsWm81BsXRNNJjzHGCqKf8RMhd8oZP55+2/SVRBwnkGyUNCuDPrJcymC95ZT2NW/KeWkz28HF2i31xQmecT2r3lQRSM8acvOXQsNEDCDvJvCzJT9c2AnsnO24r6arPXs/UWAxOI+MjclXPLkLD6uTHV+Oo8XZ7bOjegD5hL6/bKUWnNMurQNGrmi/jvqsCFLDKftl7ajuxKjtodnSuwhoY7NQy8="),
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
		saml2_sso_url = "https://samltest.id/saml/sp"
		saml2_provider = "CUSTOM"
		saml2_x509_cert = "MIIERTCCAq2gAwIBAgIJAKmtzjCD1+tqMA0GCSqGSIb3DQEBCwUAMDUxMzAxBgNVBAMTKmlwLTE3Mi0zMS0yOC02NC51cy13ZXN0LTIuY29tcHV0ZS5pbnRlcm5hbDAeFw0xODA4MTgyMzI0MjNaFw0yODA4MTUyMzI0MjNaMDUxMzAxBgNVBAMTKmlwLTE3Mi0zMS0yOC02NC51cy13ZXN0LTIuY29tcHV0ZS5pbnRlcm5hbDCCAaIwDQYJKoZIhvcNAQEBBQADggGPADCCAYoCggGBALhUlY3SkIOze+l8y6dBzM6p7B8OykJWlwizszU16Lih8D7KLhNJfahoVxbPxB3YFM/81PJLOeK2krvJ5zY6CJyQY3sPQAkZKI7I8qq9lmZ2g4QPqybNstXS6YUXJNUt/ixbbK/N97+LKTiSutbD1J7AoFnouMuLjlhN5VRZ43jez4xLSHVZaYuUFKn01Y9oLKbj46LQnZnJCAGpTgPqEQJr6GpVGw43bKyUpGoaPrdDRgRgtPMUWgFDkgcI3QiV1lsKfBs1t1E2UA7ACFnlJZpEuBtwgivzo3VeitiSaF3Jxh25EY5/vABpcgQQRz3RH2l8MMKdRsxb8VT3yh2S+CX55s+cN67LiCPr6f2u+KS1iKfB9mWN6o2S4lcmo82HIBbsuXJV0oA1HrGMyyc4Y9nng/I8iuAp8or1JrWRHQ+8NzO85DWK0rtvtLPxkvw0HK32glyuOP/9F05Z7+tiVIgn67buC0EdoUm1RSpibqmB1ST2PikslOlVbJuy4Ah93wIDAQABo1gwVjA1BgNVHREELjAsgippcC0xNzItMzEtMjgtNjQudXMtd2VzdC0yLmNvbXB1dGUuaW50ZXJuYWwwHQYDVR0OBBYEFAdsTxYfulJ5yunYtgYJHC9IcevzMA0GCSqGSIb3DQEBCwUAA4IBgQB3J6i7KreiHL8NPMglfWLHk1PZOgvIEEpKL+GRebvcbyqgcuc3VVPylq70VvGqhJxp1q/mzLfraUiypzfWFGm9zfwIg0H5TqRZYEPTvgIhIICjaDWRwZBDJG8D5G/KoV60DlUG0crPBlIuCCr/SRa5ZoDQqvucTfr3Rx4Ha6koXFSjoSXllR+jn4GnInhm/WH137a+v35PUcffNxfuehoGn6i4YeXF3cwJK4e35cOFW+dLbnaLk+Ty7HOGvpw86h979C6mJ9qEHYgq9rQyzlSPbLZGZSgVcIezunOaOsWm81BsXRNNJjzHGCqKf8RMhd8oZP55+2/SVRBwnkGyUNCuDPrJcymC95ZT2NW/KeWkz28HF2i31xQmecT2r3lQRSM8acvOXQsNEDCDvJvCzJT9c2AnsnO24r6arPXs/UWAxOI+MjclXPLkLD6uTHV+Oo8XZ7bOjegD5hL6/bKUWnNMurQNGrmi/jvqsCFLDKftl7ajuxKjtodnSuwhoY7NQy8="
		enabled = false
	}
	`, name)
}
