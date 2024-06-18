package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Saml2Integration_basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	issuer, issuer2 := acc.TestClient().Ids.Alpha(), acc.TestClient().Ids.Alpha()
	cert, cert2 := random.GenerateX509(t), random.GenerateX509(t)
	validUrl, validUrl2 := "http://example.com", "http://example2.com"
	acsURL := acc.TestClient().Context.ACSURL(t)
	issuerURL := acc.TestClient().Context.IssuerURL(t)

	m := func(issuer, provider, ssoUrl, x509Cert string, complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"name":            config.StringVariable(id.Name()),
			"saml2_issuer":    config.StringVariable(issuer),
			"saml2_provider":  config.StringVariable(provider),
			"saml2_sso_url":   config.StringVariable(ssoUrl),
			"saml2_x509_cert": config.StringVariable(x509Cert),
		}
		if complete {
			c["enabled"] = config.BoolVariable(true)
			c["allowed_email_patterns"] = config.ListVariable(config.StringVariable("^(.+dev)@example.com$"))
			c["allowed_user_domains"] = config.ListVariable(config.StringVariable("example.com"))
			c["comment"] = config.StringVariable("foo")
			c["saml2_enable_sp_initiated"] = config.BoolVariable(true)
			c["saml2_force_authn"] = config.BoolVariable(true)
			c["saml2_post_logout_redirect_url"] = config.StringVariable(validUrl)
			c["saml2_requested_nameid_format"] = config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified))
			c["saml2_sign_request"] = config.BoolVariable(true)
			c["saml2_snowflake_acs_url"] = config.StringVariable(acsURL)
			c["saml2_snowflake_issuer_url"] = config.StringVariable(issuerURL)
			c["saml2_sp_initiated_login_page_label"] = config.StringVariable("foo")
			// TODO(SNOW-1479617): set saml2_snowflake_x509_cert
		}
		return c
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/basic"),
				ConfigVariables: m(issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_acs_url"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_metadata"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_digest_methods_used"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_signature_methods_used"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_x509_cert"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables: m(issuer2, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl2, cert2, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_digest_methods_used"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_enable_sp_initiated", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sign_request", "true"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_signature_methods_used"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_metadata"),
					// TODO(SNOW-1479617): assert a proper value
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_x509_cert"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert2),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "created_on"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables:         m(issuer2, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl2, cert2, true),
				ResourceName:            "snowflake_saml2_integration.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"saml2_snowflake_metadata"},
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/basic"),
				ConfigVariables: m(issuer2, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl2, cert2, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "0"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "0"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_digest_methods_used"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_enable_sp_initiated", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_requested_nameid_format", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sign_request", "false"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_signature_methods_used"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_acs_url"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_metadata"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_x509_cert"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert2),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "created_on"),
				),
			},
		},
	})
}

func TestAcc_Saml2Integration_complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	issuer := acc.TestClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "http://example.com"
	acsURL := acc.TestClient().Context.ACSURL(t)
	issuerURL := acc.TestClient().Context.IssuerURL(t)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"allowed_email_patterns":              config.ListVariable(config.StringVariable("^(.+dev)@example.com$")),
			"allowed_user_domains":                config.ListVariable(config.StringVariable("example.com")),
			"comment":                             config.StringVariable("foo"),
			"enabled":                             config.BoolVariable(true),
			"name":                                config.StringVariable(id.Name()),
			"saml2_enable_sp_initiated":           config.BoolVariable(true),
			"saml2_force_authn":                   config.BoolVariable(true),
			"saml2_issuer":                        config.StringVariable(issuer),
			"saml2_post_logout_redirect_url":      config.StringVariable(validUrl),
			"saml2_provider":                      config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
			"saml2_requested_nameid_format":       config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
			"saml2_sign_request":                  config.BoolVariable(true),
			"saml2_snowflake_acs_url":             config.StringVariable(acsURL),
			"saml2_snowflake_issuer_url":          config.StringVariable(issuerURL),
			"saml2_sp_initiated_login_page_label": config.StringVariable("foo"),
			"saml2_sso_url":                       config.StringVariable(validUrl),
			"saml2_x509_cert":                     config.StringVariable(cert),
			// TODO(SNOW-1479617): set saml2_snowflake_x509_cert
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_digest_methods_used"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_enable_sp_initiated", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sign_request", "true"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_signature_methods_used"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_metadata"),
					// TODO(SNOW-1479617): assert a proper value
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "saml2_snowflake_x509_cert"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "created_on"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables:         m(),
				ResourceName:            "snowflake_saml2_integration.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"saml2_snowflake_metadata"},
			},
		},
	})
}

func TestAcc_Saml2Integration_invalid(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"allowed_email_patterns":              config.ListVariable(config.StringVariable("foo")),
			"allowed_user_domains":                config.ListVariable(config.StringVariable("foo")),
			"comment":                             config.StringVariable("foo"),
			"enabled":                             config.BoolVariable(true),
			"name":                                config.StringVariable("foo"),
			"saml2_enable_sp_initiated":           config.BoolVariable(true),
			"saml2_force_authn":                   config.BoolVariable(true),
			"saml2_issuer":                        config.StringVariable("foo"),
			"saml2_post_logout_redirect_url":      config.StringVariable("foo"),
			"saml2_provider":                      config.StringVariable("invalid"),
			"saml2_requested_nameid_format":       config.StringVariable("invalid"),
			"saml2_sign_request":                  config.BoolVariable(true),
			"saml2_snowflake_acs_url":             config.StringVariable("foo"),
			"saml2_snowflake_issuer_url":          config.StringVariable("foo"),
			"saml2_snowflake_x509_cert":           config.StringVariable("foo"),
			"saml2_sp_initiated_login_page_label": config.StringVariable("foo"),
			"saml2_sso_url":                       config.StringVariable("foo"),
			"saml2_x509_cert":                     config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			fmt.Sprintf(`expected saml2_provider to be one of %q, got invalid`, sdk.AsStringList(sdk.AllSaml2SecurityIntegrationSaml2Providers)),
			fmt.Sprintf(`expected saml2_requested_nameid_format to be one of %q, got invalid`, sdk.AsStringList(sdk.AllSaml2SecurityIntegrationSaml2RequestedNameidFormats)),
		}),

		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_Saml2Integration_InvalidIncomplete(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name": config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			`The argument "saml2_issuer" is required, but no definition was found.`,
			`The argument "saml2_provider" is required, but no definition was found.`,
			`The argument "saml2_sso_url" is required, but no definition was found.`,
			`The argument "saml2_x509_cert" is required, but no definition was found.`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}
