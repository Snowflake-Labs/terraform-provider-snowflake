package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	resourcenames "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OauthIntegrationForCustomClients_Basic(t *testing.T) {
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

	m := func(complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"name":               config.StringVariable(id.Name()),
			"oauth_client_type":  config.StringVariable(string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
			"oauth_redirect_uri": config.StringVariable(validUrl),
			"blocked_roles_list": config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
		}
		if complete {
			c["blocked_roles_list"] = config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN"), config.StringVariable(blockedRole.ID().Name()))
			c["comment"] = config.StringVariable(comment)
			c["enabled"] = config.BoolVariable(true)
			c["network_policy"] = config.StringVariable(networkPolicy.Name)
			c["oauth_allow_non_tls_redirect_uri"] = config.BoolVariable(true)
			c["oauth_allowed_authorization_endpoints"] = config.SetVariable(config.StringVariable("http://allowed.com"))
			c["oauth_allowed_token_endpoints"] = config.SetVariable(config.StringVariable("http://allowed.com"))
			c["oauth_authorization_endpoint"] = config.StringVariable("http://auth.com")
			c["oauth_client_rsa_public_key"] = config.StringVariable(key)
			c["oauth_client_rsa_public_key_2"] = config.StringVariable(key)
			c["oauth_enforce_pkce"] = config.BoolVariable(true)
			c["oauth_issue_refresh_tokens"] = config.BoolVariable(true)
			c["oauth_refresh_token_validity"] = config.IntegerVariable(86400)
			c["oauth_token_endpoint"] = config.StringVariable("http://auth.com")
			c["oauth_use_secondary_roles"] = config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))
			c["pre_authorized_roles_list"] = config.SetVariable(config.StringVariable(preAuthorizedRole.ID().Name()))
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
			// create with empty optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "enabled", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_allow_non_tls_redirect_uri", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_enforce_pkce", resources.BooleanDefault),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_issue_refresh_tokens", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "network_policy", ""),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key"),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_enforce_pkce.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.pre_authorized_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.network_policy.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.comment.0.value", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/basic"),
				ConfigVariables: m(false),
				ResourceName:    "snowflake_oauth_integration_for_custom_clients.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_redirect_uri", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_allow_non_tls_redirect_uri", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_enforce_pkce", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.#", "0"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "7776000"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", ""),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key_2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/complete"),
				ConfigVariables: m(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_custom_clients.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_allow_non_tls_redirect_uri", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_enforce_pkce", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "network_policy", networkPolicy.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key", key),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2", key),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_enforce_pkce.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name()),
					// Not asserted, because it also contains other default roles
					// resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.network_policy.0.value", networkPolicy.Name),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_fp.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.comment.0.value", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/complete"),
				ConfigVariables: m(true),
				ResourceName:    "snowflake_oauth_integration_for_custom_clients.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_redirect_uri", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_allow_non_tls_redirect_uri", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_enforce_pkce", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "3"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", networkPolicy.ID().FullyQualifiedName()),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key_2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// change externally
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/complete"),
				ConfigVariables: m(true),
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateOauthForClients(t, sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).WithUnset(
						*sdk.NewOauthForCustomClientsIntegrationUnsetRequest().
							WithEnabled(true).
							WithNetworkPolicy(true).
							WithOauthUseSecondaryRoles(true).
							WithOauthClientRsaPublicKey(true).
							WithOauthClientRsaPublicKey2(true),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_custom_clients.test", plancheck.ResourceActionUpdate),

						planchecks.ExpectDrift("snowflake_oauth_integration_for_custom_clients.test", "enabled", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_custom_clients.test", "network_policy", sdk.String(networkPolicy.ID().FullyQualifiedName()), sdk.String("")),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles", sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone))),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key", sdk.String(key), sdk.String(key)),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2", sdk.String(key), sdk.String(key)),

						planchecks.ExpectChange("snowflake_oauth_integration_for_custom_clients.test", "enabled", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_oauth_integration_for_custom_clients.test", "network_policy", tfjson.ActionUpdate, sdk.String(""), sdk.String(networkPolicy.ID().FullyQualifiedName())),
						planchecks.ExpectChange("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles", tfjson.ActionUpdate, sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))),
						planchecks.ExpectChange("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key", tfjson.ActionUpdate, sdk.String(key), sdk.String(key)),
						planchecks.ExpectChange("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2", tfjson.ActionUpdate, sdk.String(key), sdk.String(key)),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_allow_non_tls_redirect_uri", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_enforce_pkce", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "network_policy", networkPolicy.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key", key),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2", key),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_enforce_pkce.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name()),
					// Not asserted, because it also contains other default roles
					// resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.network_policy.0.value", networkPolicy.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.comment.0.value", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/basic"),
				ConfigVariables: m(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_custom_clients.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "enabled", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_allow_non_tls_redirect_uri", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_enforce_pkce", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_issue_refresh_tokens", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "network_policy", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_enforce_pkce.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.pre_authorized_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.network_policy.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.comment.0.value", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_Complete(t *testing.T) {
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

	configVariables := config.Variables{
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/complete"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_allow_non_tls_redirect_uri", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_enforce_pkce", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "network_policy", networkPolicy.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key", key),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2", key),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_enforce_pkce.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name()),
					// Not asserted, because it also contains other default roles
					// resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.network_policy.0.value", networkPolicy.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_fp.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.comment.0.value", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/complete"),
				ConfigVariables: configVariables,
				ResourceName:    "snowflake_oauth_integration_for_custom_clients.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_redirect_uri", validUrl),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_allow_non_tls_redirect_uri", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_enforce_pkce", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "pre_authorized_roles_list.0", preAuthorizedRole.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "3"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", networkPolicy.ID().FullyQualifiedName()),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client_rsa_public_key_2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_DefaultValues(t *testing.T) {
	validUrl := "https://example.com"
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(defaultsSet bool) config.Variables {
		variables := config.Variables{
			"name":               config.StringVariable(id.Name()),
			"oauth_client_type":  config.StringVariable(string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
			"oauth_redirect_uri": config.StringVariable(validUrl),
			"blocked_roles_list": config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
		}

		variables["enabled"] = config.BoolVariable(false)
		variables["oauth_allow_non_tls_redirect_uri"] = config.BoolVariable(false)
		variables["oauth_enforce_pkce"] = config.BoolVariable(false)
		variables["oauth_use_secondary_roles"] = config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone))
		variables["pre_authorized_roles_list"] = config.SetVariable()
		variables["oauth_issue_refresh_tokens"] = config.BoolVariable(false)
		variables["oauth_refresh_token_validity"] = config.IntegerVariable(7776000)
		variables["network_policy"] = config.StringVariable("")
		variables["oauth_client_rsa_public_key"] = config.StringVariable("")
		variables["oauth_client_rsa_public_key_2"] = config.StringVariable("")
		variables["comment"] = config.StringVariable("")

		return variables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/default_values"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_allow_non_tls_redirect_uri", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_enforce_pkce", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_refresh_token_validity", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "network_policy", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_enforce_pkce.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.pre_authorized_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.network_policy.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.comment.0.value", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/basic"),
				ConfigVariables: configVariables(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_type", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "enabled", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_allow_non_tls_redirect_uri", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_enforce_pkce", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_use_secondary_roles", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "pre_authorized_roles_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_issue_refresh_tokens", resources.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "network_policy", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "oauth_client_rsa_public_key_2", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_enforce_pkce.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.pre_authorized_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.network_policy.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_rsa_public_key_2_fp.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.comment.0.value", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_custom_clients.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_Invalid(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	invalidUseSecondaryRoles := config.Variables{
		"name":                      config.StringVariable(id.Name()),
		"oauth_client_type":         config.StringVariable(string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
		"oauth_redirect_uri":        config.StringVariable("https://example.com"),
		"oauth_use_secondary_roles": config.StringVariable("invalid"),
		"blocked_roles_list":        config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
	}

	invalidClientType := config.Variables{
		"name":                      config.StringVariable(id.Name()),
		"oauth_client_type":         config.StringVariable("invalid"),
		"oauth_redirect_uri":        config.StringVariable("https://example.com"),
		"oauth_use_secondary_roles": config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
		"blocked_roles_list":        config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/invalid"),
				ConfigVariables: invalidUseSecondaryRoles,
				ExpectError:     regexp.MustCompile(`Error: invalid OauthSecurityIntegrationUseSecondaryRolesOption: INVALID`),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForCustomClients/invalid"),
				ConfigVariables: invalidClientType,
				ExpectError:     regexp.MustCompile(`Error: invalid OauthSecurityIntegrationClientTypeOption: INVALID`),
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resourcenames.OauthIntegrationForCustomClients),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: oauthIntegrationForCustomClientsBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "id", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   oauthIntegrationForCustomClientsBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForCustomClients_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resourcenames.OauthIntegrationForCustomClients),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             oauthIntegrationForCustomClientsBasicConfig(quotedId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "id", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   oauthIntegrationForCustomClientsBasicConfig(quotedId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_custom_clients.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_custom_clients.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_custom_clients.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func oauthIntegrationForCustomClientsBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_oauth_integration_for_custom_clients" "test" {
  name               = "%s"
  oauth_client_type  = "CONFIDENTIAL"
  oauth_redirect_uri = "https://example.com"
  blocked_roles_list = [ "ACCOUNTADMIN", "SECURITYADMIN" ]
}
`, name)
}
