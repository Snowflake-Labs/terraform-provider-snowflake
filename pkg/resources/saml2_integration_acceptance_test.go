package resources_test

import (
	"fmt"
	"maps"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Saml2Integration_basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	issuer, issuer2 := acc.TestClient().Ids.Alpha(), acc.TestClient().Ids.Alpha()
	cert, cert2 := random.GenerateX509(t), random.GenerateX509(t)
	validUrl, validUrl2 := "http://example.com", "http://example2.com"
	acsURL := acc.TestClient().Context.ACSURL(t)
	issuerURL := acc.TestClient().Context.IssuerURL(t)

	m := func(issuer, provider, ssoUrl, x509Cert string, complete bool, unset bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"name":            config.StringVariable(id.Name()),
			"saml2_issuer":    config.StringVariable(issuer),
			"saml2_provider":  config.StringVariable(provider),
			"saml2_sso_url":   config.StringVariable(ssoUrl),
			"saml2_x509_cert": config.StringVariable(x509Cert),
		}
		if complete {
			c["enabled"] = config.BoolVariable(true)
			c["comment"] = config.StringVariable("foo")
			c["saml2_enable_sp_initiated"] = config.BoolVariable(true)
			c["saml2_force_authn"] = config.BoolVariable(true)
			c["saml2_post_logout_redirect_url"] = config.StringVariable(validUrl)
			c["saml2_requested_nameid_format"] = config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified))
			c["saml2_sign_request"] = config.BoolVariable(true)
			// TODO(SNOW-1479617): set saml2_snowflake_x509_cert
			c["saml2_snowflake_acs_url"] = config.StringVariable(acsURL)
			c["saml2_snowflake_issuer_url"] = config.StringVariable(issuerURL)
			c["saml2_sp_initiated_login_page_label"] = config.StringVariable("foo")
			c["allowed_email_patterns"] = config.ListVariable(config.StringVariable("^(.+dev)@example.com$"))
			c["allowed_user_domains"] = config.ListVariable(config.StringVariable("example.com"))
		}
		// When unsetting, we have to keep those to prevent conditional force new being triggered
		if unset {
			c["saml2_snowflake_acs_url"] = config.StringVariable(acsURL)
			c["saml2_snowflake_issuer_url"] = config.StringVariable(issuerURL)
			c["saml2_sp_initiated_login_page_label"] = config.StringVariable("foo")
			c["allowed_email_patterns"] = config.ListVariable(config.StringVariable("^(.+dev)@example.com$"))
			c["allowed_user_domains"] = config.ListVariable(config.StringVariable("example.com"))
		}
		return c
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/basic"),
				ConfigVariables: m(issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert),
					resource.TestCheckNoResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_enable_sp_initiated", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sign_request", "unknown"),
					resource.TestCheckNoResourceAttr("snowflake_saml2_integration.test", "saml2_requested_nameid_format"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "unknown"),
					resource.TestCheckNoResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url"),
					resource.TestCheckNoResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url"),
					resource.TestCheckNoResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains"),
					resource.TestCheckNoResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_x509_cert.0.value", cert),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_enable_sp_initiated.0.value", "false"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sign_request.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_requested_nameid_format.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.comment.0.value", ""),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "show_output.0.created_on"),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/basic"),
				ConfigVariables: m(issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert, false, false),
				ResourceName:    "snowflake_saml2_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sso_url", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_x509_cert", cert),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sp_initiated_login_page_label", ""),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_enable_sp_initiated", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sign_request", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_requested_nameid_format", ""),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_post_logout_redirect_url", ""),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_force_authn", "false"),
					importchecks.TestCheckResourceAttrInstanceStateSet(id.Name(), "saml2_snowflake_issuer_url"),
					importchecks.TestCheckResourceAttrInstanceStateSet(id.Name(), "saml2_snowflake_acs_url"),
					importchecks.TestCheckResourceAttrNotInInstanceState(id.Name(), "allowed_user_domains"),
					importchecks.TestCheckResourceAttrNotInInstanceState(id.Name(), "allowed_email_patterns"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", ""),
				),
			},
			// set optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables: m(issuer2, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl2, cert2, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_enable_sp_initiated", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sign_request", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_issuer.0.value", issuer2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sso_url.0.value", validUrl2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_x509_cert.0.value", cert2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_enable_sp_initiated.0.value", "true"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sign_request.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_post_logout_redirect_url.0.value", "http://example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "true"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.comment.0.value", "foo"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "show_output.0.created_on"),
				),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables: m(issuer2, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl2, cert2, true, false),
				ResourceName:    "snowflake_saml2_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_issuer", issuer2),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sso_url", validUrl2),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_x509_cert", cert2),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sp_initiated_login_page_label", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_enable_sp_initiated", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sign_request", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_post_logout_redirect_url", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_force_authn", "true"),
					importchecks.TestCheckResourceAttrInstanceStateSet(id.Name(), "saml2_snowflake_issuer_url"),
					importchecks.TestCheckResourceAttrInstanceStateSet(id.Name(), "saml2_snowflake_acs_url"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_user_domains.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_user_domains.0", "example.com"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_email_patterns.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", "foo"),
				),
			},
			// change values externally
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables: m(issuer2, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl2, cert2, true, false),
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateSaml2(t, sdk.NewAlterSaml2SecurityIntegrationRequest(id).
						WithUnset(*sdk.NewSaml2IntegrationUnsetRequest().
							WithSaml2RequestedNameidFormat(true).
							WithSaml2PostLogoutRedirectUrl(true).
							WithSaml2ForceAuthn(true).
							WithComment(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift("snowflake_saml2_integration.test", "saml2_requested_nameid_format", sdk.String(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)), sdk.String(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress))),
						planchecks.ExpectDrift("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", sdk.String(validUrl), sdk.String("")),
						planchecks.ExpectDrift("snowflake_saml2_integration.test", "saml2_force_authn", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift("snowflake_saml2_integration.test", "comment", sdk.String("foo"), sdk.String("")),

						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_requested_nameid_format", tfjson.ActionUpdate, sdk.String(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress)), sdk.String(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified))),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", tfjson.ActionUpdate, sdk.String(""), sdk.String(validUrl)),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "comment", tfjson.ActionUpdate, sdk.String(""), sdk.String("foo")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_enable_sp_initiated", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sign_request", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_issuer.0.value", issuer2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sso_url.0.value", validUrl2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_x509_cert.0.value", cert2),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_enable_sp_initiated.0.value", "true"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sign_request.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_post_logout_redirect_url.0.value", "http://example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "true"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.comment.0.value", "foo"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "show_output.0.created_on"),
				),
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/recreates"),
				ConfigVariables: m(issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_enable_sp_initiated", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sign_request", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_requested_nameid_format", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_x509_cert.0.value", cert),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_enable_sp_initiated.0.value", "false"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sign_request.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.comment.0.value", ""),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "show_output.0.created_on"),
				),
			},
		},
	})
}

func saml2ConfigWithAuthn(name, issuer, provider, ssoUrl, x509Cert string, forceAuthn bool) string {
	return fmt.Sprintf(`
resource "snowflake_saml2_integration" "test" {
	name           = "%s"
	saml2_issuer = "%s"
	saml2_provider = "%s"
	saml2_sso_url = "%s"
	saml2_x509_cert = <<EOT
%s
EOT
	saml2_force_authn = %t
}
`, name, issuer, provider, ssoUrl, x509Cert, forceAuthn)
}

func saml2Config(name, issuer, provider, ssoUrl, x509Cert string) string {
	return fmt.Sprintf(`
resource "snowflake_saml2_integration" "test" {
	name           = "%s"
	saml2_issuer = "%s"
	saml2_provider = "%s"
	saml2_sso_url = "%s"
	saml2_x509_cert = <<EOT
%s
EOT
}
`, name, issuer, provider, ssoUrl, x509Cert)
}

func TestAcc_Saml2Integration_forceAuthn(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	issuer := acc.TestClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "http://example.com"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			// set up with concrete type
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_saml2_integration.test", "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionCreate, nil, sdk.String("true")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Config: saml2ConfigWithAuthn(id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "true"),
				),
			},
			// import when type in config
			{
				ResourceName: "snowflake_saml2_integration.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_force_authn", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "describe_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "describe_output.0.saml2_force_authn.0.value", "true"),
				),
			},
			// change type in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_saml2_integration.test", "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Config: saml2ConfigWithAuthn(id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// change back to non-default
			{
				Config: saml2ConfigWithAuthn(id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert, true),
			},
			// remove non-default type from config
			{
				Config: saml2Config(id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_saml2_integration.test", plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails("snowflake_saml2_integration.test", "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionUpdate, sdk.String("true"), sdk.String("unknown")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// add config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_saml2_integration.test", "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionUpdate, sdk.String("unknown"), sdk.String("unknown")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Config: saml2ConfigWithAuthn(id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// remove type from config but update externally to default (still expecting non-empty plan because we do not know the default)
			{
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateSaml2ForceAuthn(t, id, false)
				},
				Config: saml2Config(id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails("snowflake_saml2_integration.test", "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionUpdate, sdk.String("false"), sdk.String("unknown")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// change the size externally
			{
				PreConfig: func() {
					// we change the type to the type different from default, expecting action
					acc.TestClient().SecurityIntegration.UpdateSaml2ForceAuthn(t, id, true)
				},
				Config: saml2Config(id.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails("snowflake_saml2_integration.test", "saml2_force_authn", "describe_output"),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionUpdate, sdk.String("true"), sdk.String("unknown")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// import when no type in config
			{
				ResourceName: "snowflake_saml2_integration.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_force_authn", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "describe_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "describe_output.0.saml2_force_authn.0.value", "false"),
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
		CheckDestroy: acc.CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sso_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_x509_cert", cert),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_enable_sp_initiated", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sign_request", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_x509_cert.0.value", cert),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_enable_sp_initiated.0.value", "true"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sign_request.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_post_logout_redirect_url.0.value", "http://example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "true"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.comment.0.value", "foo"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_saml2_integration.test", "show_output.0.created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/complete"),
				ConfigVariables: m(),
				ResourceName:    "snowflake_saml2_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sso_url", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_provider", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_x509_cert", cert),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sp_initiated_login_page_label", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_enable_sp_initiated", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_sign_request", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_requested_nameid_format", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_post_logout_redirect_url", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_force_authn", "true"),
					importchecks.TestCheckResourceAttrInstanceStateSet(id.Name(), "saml2_snowflake_issuer_url"),
					importchecks.TestCheckResourceAttrInstanceStateSet(id.Name(), "saml2_snowflake_acs_url"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_user_domains.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_user_domains.0", "example.com"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_email_patterns.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "allowed_email_patterns.0", "^(.+dev)@example.com$"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", "foo"),
				),
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
			`Error: invalid Saml2SecurityIntegrationSaml2RequestedNameidFormatOption: invalid`,
			`Error: invalid Saml2SecurityIntegrationSaml2ProviderOption: INVALID`,
		}),
		CheckDestroy: acc.CheckDestroy(t, resources.Saml2SecurityIntegration),
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
		CheckDestroy: acc.CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_Saml2Integration_ForceNewIfEmpty(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	issuer := acc.TestClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "http://example.com"
	acsURL := acc.TestClient().Context.ACSURL(t)
	issuerURL := acc.TestClient().Context.IssuerURL(t)

	commonValues := map[string]config.Variable{
		"name":                                config.StringVariable(id.Name()),
		"saml2_issuer":                        config.StringVariable(issuer),
		"saml2_provider":                      config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
		"saml2_sso_url":                       config.StringVariable(validUrl),
		"saml2_x509_cert":                     config.StringVariable(cert),
		"saml2_snowflake_acs_url":             config.StringVariable(acsURL),
		"saml2_snowflake_issuer_url":          config.StringVariable(issuerURL),
		"saml2_sp_initiated_login_page_label": config.StringVariable("label"),
		"allowed_email_patterns":              config.ListVariable(config.StringVariable("^(.+dev)@example.com$")),
		"allowed_user_domains":                config.ListVariable(config.StringVariable("example.com")),
	}

	emptySpInitiatedLoginPageLabel := maps.Clone(commonValues)
	emptySpInitiatedLoginPageLabel["saml2_sp_initiated_login_page_label"] = config.StringVariable("")

	emptySnowflakeAcsUrl := maps.Clone(commonValues)
	emptySnowflakeAcsUrl["saml2_snowflake_acs_url"] = config.StringVariable("")

	emptySnowflakeIssuerUrl := maps.Clone(commonValues)
	emptySnowflakeIssuerUrl["saml2_snowflake_issuer_url"] = config.StringVariable("")

	emptyAllowedEmailPatterns := maps.Clone(commonValues)
	emptyAllowedEmailPatterns["allowed_email_patterns"] = config.ListVariable()

	emptyAllowedUserDomains := maps.Clone(commonValues)
	emptyAllowedUserDomains["allowed_user_domains"] = config.ListVariable()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/recreates"),
				ConfigVariables: commonValues,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/recreates"),
				ConfigVariables: emptySpInitiatedLoginPageLabel,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_saml2_integration.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.0", "^(.+dev)@example.com$"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/recreates"),
				ConfigVariables: emptySnowflakeIssuerUrl,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_saml2_integration.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.0", "^(.+dev)@example.com$"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value", strings.ToLower(issuerURL)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/recreates"),
				ConfigVariables: emptySnowflakeAcsUrl,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_saml2_integration.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.0", "^(.+dev)@example.com$"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value", strings.ToLower(acsURL)),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/recreates"),
				ConfigVariables: emptyAllowedEmailPatterns,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_saml2_integration.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.0", "example.com"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_email_patterns.#", "0"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_email_patterns.0.value", "[]"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/recreates"),
				ConfigVariables: emptyAllowedUserDomains,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_saml2_integration.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_sp_initiated_login_page_label", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_issuer_url", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_snowflake_acs_url", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "allowed_user_domains.#", "0"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_sp_initiated_login_page_label.0.value", "label"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_issuer_url.0.value", issuerURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_snowflake_acs_url.0.value", acsURL),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.allowed_user_domains.0.value", "[]"),
				),
			},
		},
	})
}

func TestAcc_Saml2Integration_DefaultValues(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	issuer := acc.TestClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "http://example.com"

	configVariables := config.Variables{
		"name":            config.StringVariable(id.Name()),
		"saml2_issuer":    config.StringVariable(issuer),
		"saml2_provider":  config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
		"saml2_sso_url":   config.StringVariable(validUrl),
		"saml2_x509_cert": config.StringVariable(cert),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// create with valid "zero" values
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/zero_values"),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange("snowflake_saml2_integration.test", "enabled", tfjson.ActionCreate, nil, sdk.String("false")),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionCreate, nil, sdk.String("false")),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", tfjson.ActionCreate, nil, sdk.String("")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "show_output", true),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", ""),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.enabled", "false"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// remove all from config (to validate that unset is run correctly)
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/basic"),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange("snowflake_saml2_integration.test", "enabled", tfjson.ActionUpdate, sdk.String("false"), sdk.String("unknown")),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionUpdate, sdk.String("false"), sdk.String("unknown")),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", tfjson.ActionUpdate, sdk.String(""), sdk.String("")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "show_output", true),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "unknown"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", ""),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.enabled", "false"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// add valid "zero" values again (to validate if set is run correctly)
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/zero_values"),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange("snowflake_saml2_integration.test", "enabled", tfjson.ActionUpdate, sdk.String("unknown"), sdk.String("unknown")),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_force_authn", tfjson.ActionUpdate, sdk.String("unknown"), sdk.String("unknown")),
						planchecks.ExpectChange("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", tfjson.ActionUpdate, sdk.String(""), sdk.String("")),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "show_output", true),
						planchecks.ExpectComputed("snowflake_saml2_integration.test", "describe_output", true),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_force_authn", "false"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "saml2_post_logout_redirect_url", ""),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "show_output.0.enabled", "false"),

					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_saml2_integration.test", "describe_output.0.saml2_force_authn.0.value", "false"),
				),
			},
			// import zero values
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Saml2Integration/zero_values"),
				ConfigVariables: configVariables,
				ImportState:     true,
				ResourceName:    "snowflake_saml2_integration.test",
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),

					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_force_authn", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "saml2_post_logout_redirect_url", ""),

					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.0.enabled", "false"),

					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "describe_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "describe_output.0.saml2_force_authn.0.value", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "describe_output.0.saml2_post_logout_redirect_url.0.value", ""),
				),
			},
		},
	})
}
