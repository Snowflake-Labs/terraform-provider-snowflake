package datasources_test

import (
	"maps"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecurityIntegrations_MultipleTypes(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	issuer := acc.TestClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "http://example.com"

	role := snowflakeroles.GenericScimProvisioner
	configVariables := config.Variables{
		// saml2
		"name_1":          config.StringVariable(idOne.Name()),
		"saml2_issuer":    config.StringVariable(issuer),
		"saml2_provider":  config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
		"saml2_sso_url":   config.StringVariable(validUrl),
		"saml2_x509_cert": config.StringVariable(cert),
		// scim
		"name_2":      config.StringVariable(idTwo.Name()),
		"scim_client": config.StringVariable(string(sdk.ScimSecurityIntegrationScimClientGeneric)),
		"run_as_role": config.StringVariable(role.Name()),
		"enabled":     config.BoolVariable(true),

		"like": config.StringVariable(prefix + "%"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/multiple_types"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "2"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", idOne.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.run_as_role.0.value", "GENERIC_SCIM_PROVISIONER"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.sync_password.0.value", "true"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.show_output.0.name", idTwo.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.1.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.describe_output.0.saml2_provider.0.value", "CUSTOM"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.1.describe_output.0.saml2_x509_cert.0.value", cert),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_ApiAuthenticationWithAuthorizationCodeGrant(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                      config.StringVariable("foo"),
			"enabled":                      config.BoolVariable(true),
			"name":                         config.StringVariable(id.Name()),
			"oauth_access_token_validity":  config.IntegerVariable(42),
			"oauth_authorization_endpoint": config.StringVariable("https://example.com"),
			"oauth_client_auth_method":     config.StringVariable(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
			"oauth_client_id":              config.StringVariable("foo"),
			"oauth_client_secret":          config.StringVariable("foo"),
			"oauth_refresh_token_validity": config.IntegerVariable(12345),
			"oauth_token_endpoint":         config.StringVariable("https://example.com"),
			"oauth_allowed_scopes":         config.SetVariable(config.StringVariable("foo")),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/api_authentication_with_authorization_code_grant/optionals_set"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "API_AUTHENTICATION"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_access_token_validity.0.value", "42"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "12345"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_id.0.value", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_auth_method.0.value", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allowed_scopes.0.value", "[foo]"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.parent_integration.0.value", ""),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.auth_type.0.value", "OAUTH2"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", "foo")),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/api_authentication_with_authorization_code_grant/optionals_unset"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "API_AUTHENTICATION"),
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

func TestAcc_SecurityIntegrations_ApiAuthenticationWithClientCredentials(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                      config.StringVariable("foo"),
			"enabled":                      config.BoolVariable(true),
			"name":                         config.StringVariable(id.Name()),
			"oauth_access_token_validity":  config.IntegerVariable(42),
			"oauth_client_auth_method":     config.StringVariable(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
			"oauth_client_id":              config.StringVariable("foo"),
			"oauth_client_secret":          config.StringVariable("foo"),
			"oauth_refresh_token_validity": config.IntegerVariable(12345),
			"oauth_token_endpoint":         config.StringVariable("https://example.com"),
			"oauth_allowed_scopes":         config.SetVariable(config.StringVariable("foo")),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ApiAuthenticationIntegrationWithClientCredentials),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/api_authentication_with_client_credentials/optionals_set"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "API_AUTHENTICATION"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_access_token_validity.0.value", "42"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "12345"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_id.0.value", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_auth_method.0.value", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allowed_scopes.0.value", "[foo]"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.parent_integration.0.value", ""),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.auth_type.0.value", "OAUTH2"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", "foo")),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/api_authentication_with_client_credentials/optionals_unset"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "API_AUTHENTICATION"),
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

func TestAcc_SecurityIntegrations_ExternalOauth(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	issuer := random.String()
	t.Cleanup(roleCleanup)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_allowed_roles_list":               config.SetVariable(config.StringVariable(role.ID().Name())),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(issuer),
			"external_oauth_jws_keys_url":                     config.SetVariable(config.StringVariable("https://example.com")),
			"external_oauth_scope_delimiter":                  config.StringVariable("."),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable(id.Name()),
			"external_oauth_type": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/external_oauth/optionals_set"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.external_oauth_allowed_roles_list.0.value", role.ID().Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", "foo"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/external_oauth/optionals_unset"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
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

func TestAcc_SecurityIntegrations_OauthForCustomClients(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)

	preAuthorizedRole, preauthorizedRoleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(preauthorizedRoleCleanup)

	blockedRole, blockedRoleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(blockedRoleCleanup)

	validUrl := "https://example.com"
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	key, _ := random.GenerateRSAPublicKey(t)
	comment := random.Comment()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                                  config.StringVariable(id.Name()),
			"oauth_client_type":                     config.StringVariable(string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
			"oauth_redirect_uri":                    config.StringVariable(validUrl),
			"blocked_roles_list":                    config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN"), config.StringVariable(blockedRole.ID().Name())),
			"comment":                               config.StringVariable(comment),
			"enabled":                               config.BoolVariable(true),
			"network_policy":                        config.StringVariable(networkPolicy.ID().Name()),
			"oauth_allow_non_tls_redirect_uri":      config.BoolVariable(true),
			"oauth_allowed_authorization_endpoints": config.SetVariable(config.StringVariable("http://allowed.com")),
			"oauth_allowed_token_endpoints":         config.SetVariable(config.StringVariable("http://allowed.com")),
			"oauth_authorization_endpoint":          config.StringVariable("http://auth.com"),
			"oauth_client_rsa_public_key":           config.StringVariable(key),
			"oauth_client_rsa_public_key_2":         config.StringVariable(key),
			"oauth_enforce_pkce":                    config.BoolVariable(true),
			"oauth_issue_refresh_tokens":            config.BoolVariable(true),
			"oauth_refresh_token_validity":          config.IntegerVariable(86400),
			"oauth_token_endpoint":                  config.StringVariable("http://auth.com"),
			"oauth_use_secondary_roles":             config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
			"pre_authorized_roles_list":             config.SetVariable(config.StringVariable(preAuthorizedRole.ID().Name())),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForCustomClients),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/oauth_for_custom_clients/optionals_set"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_enforce_pkce.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name()),
					// Not asserted, because it also contains other default roles
					// resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.network_policy.0.value", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_rsa_public_key_fp.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_rsa_public_key_2_fp.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allowed_token_endpoints.0.value"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/oauth_for_custom_clients/optionals_unset"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_OauthForPartnerApplications(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                         config.StringVariable(id.Name()),
			"oauth_client":                 config.StringVariable(string(sdk.OauthSecurityIntegrationClientTableauServer)),
			"blocked_roles_list":           config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
			"enabled":                      config.BoolVariable(true),
			"oauth_issue_refresh_tokens":   config.BoolVariable(false),
			"oauth_refresh_token_validity": config.IntegerVariable(86400),
			"oauth_use_secondary_roles":    config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
			"comment":                      config.StringVariable(comment),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/oauth_for_partner_applications/optionals_set"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", comment),
					resource.TestCheckNoResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allowed_token_endpoints.0.value"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "OAUTH - TABLEAU_SERVER"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/oauth_for_partner_applications/optionals_unset"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "OAUTH - TABLEAU_SERVER"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

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
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
		"network_policy": config.StringVariable(networkPolicy.ID().Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.network_policy.0.value", networkPolicy.ID().Name()),
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
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
