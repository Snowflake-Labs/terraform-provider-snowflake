package datasources_test

import (
	"maps"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecurityIntegrations_Saml2(t *testing.T) {
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
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/saml2/optionals_set"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_x509_cert.0.value", cert),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_enable_sp_initiated.0.value", "true"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_sign_request.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_post_logout_redirect_url.0.value", "http://example.com"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_force_authn.0.value", "true"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", "foo"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/saml2/optionals_unset"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_Scim(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)
	configVariables := config.Variables{
		"name":           config.StringVariable(id.Name()),
		"comment":        config.StringVariable(comment),
		"network_policy": config.StringVariable(networkPolicy.Name),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/optionals_set"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.network_policy.0.value", sdk.NewAccountObjectIdentifier(networkPolicy.Name).Name()), // TODO(SNOW-999049): Fix during identifiers rework
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.run_as_role.0.value", "GENERIC_SCIM_PROVISIONER"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.sync_password.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", comment),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/optionals_unset"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	commonVariables := config.Variables{
		"name_1": config.StringVariable(idOne.Name()),
		"name_2": config.StringVariable(idTwo.Name()),
		"name_3": config.StringVariable(idThree.Name()),
	}

	likeConfig := config.Variables{
		"like": config.StringVariable(idOne.Name()),
	}
	maps.Copy(likeConfig, commonVariables)

	likeConfig2 := config.Variables{
		"like": config.StringVariable(prefix + "%"),
	}
	maps.Copy(likeConfig2, commonVariables)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/like"),
				ConfigVariables: likeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/like"),
				ConfigVariables: likeConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "2"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_SecurityIntegrationNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one security integration"),
			},
		},
	})
}
