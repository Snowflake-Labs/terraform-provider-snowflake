package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OauthIntegrationForPartnerApplications_Basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	validUrl := "https://example.com"
	comment := random.Comment()

	configVariables := func(complete bool) config.Variables {
		values := config.Variables{
			"name":               config.StringVariable(id.Name()),
			"oauth_client":       config.StringVariable(string(sdk.OauthSecurityIntegrationClientLooker)),
			"blocked_roles_list": config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
			"oauth_redirect_uri": config.StringVariable(validUrl),
		}
		if complete {
			values["enabled"] = config.BoolVariable(true)
			values["oauth_issue_refresh_tokens"] = config.BoolVariable(false)
			values["oauth_refresh_token_validity"] = config.IntegerVariable(86400)
			values["oauth_use_secondary_roles"] = config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))
			values["comment"] = config.StringVariable(comment)
		}
		return values
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/basic"),
				ConfigVariables: configVariables(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "default"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "default"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "-1"),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - LOOKER"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/basic"),
				ConfigVariables: configVariables(false),
				ResourceName:    "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "7776000"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_partner_applications.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - LOOKER"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete"),
				ConfigVariables: configVariables(true),
				ResourceName:    "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// change externally
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete"),
				ConfigVariables: configVariables(true),
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().
							WithBlockedRolesList(*sdk.NewBlockedRolesListRequest([]sdk.AccountObjectIdentifier{})).
							WithComment("").
							WithOauthIssueRefreshTokens(true).
							WithOauthRefreshTokenValidity(3600),
					))
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithUnset(
						*sdk.NewOauthForPartnerApplicationsIntegrationUnsetRequest().
							WithEnabled(true).
							WithOauthUseSecondaryRoles(true),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_partner_applications.test", plancheck.ResourceActionUpdate),

						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "enabled", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "comment", sdk.String(comment), sdk.String("")),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone))),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", sdk.String("86400"), sdk.String("3600")),

						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "enabled", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "comment", tfjson.ActionUpdate, sdk.String(""), sdk.String(comment)),
						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", tfjson.ActionUpdate, sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))),
						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", tfjson.ActionUpdate, sdk.String("3600"), sdk.String("86400")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - LOOKER"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/basic"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_partner_applications.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri", validUrl),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "default"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "default"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - LOOKER"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_BasicTableauDesktop(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	configVariables := func(complete bool) config.Variables {
		values := config.Variables{
			"name":               config.StringVariable(id.Name()),
			"oauth_client":       config.StringVariable(string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
			"blocked_roles_list": config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
		}
		if complete {
			values["blocked_roles_list"] = config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN"), config.StringVariable(role.ID().Name()))
			values["enabled"] = config.BoolVariable(true)
			values["oauth_issue_refresh_tokens"] = config.BoolVariable(false)
			values["oauth_refresh_token_validity"] = config.IntegerVariable(86400)
			values["oauth_use_secondary_roles"] = config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))
			values["comment"] = config.StringVariable(comment)
		}
		return values
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/basic_tableau"),
				ConfigVariables: configVariables(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "default"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "default"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "-1"),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - TABLEAU_DESKTOP"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", ""),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/basic_tableau"),
				ConfigVariables: configVariables(false),
				ResourceName:    "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "7776000"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete_tableau"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_partner_applications.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - TABLEAU_DESKTOP"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", comment),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete_tableau"),
				ConfigVariables: configVariables(true),
				ResourceName:    "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "3"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// change externally
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete_tableau"),
				ConfigVariables: configVariables(true),
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(
						*sdk.NewOauthForPartnerApplicationsIntegrationSetRequest().
							WithBlockedRolesList(*sdk.NewBlockedRolesListRequest([]sdk.AccountObjectIdentifier{})).
							WithComment("").
							WithOauthIssueRefreshTokens(true).
							WithOauthRefreshTokenValidity(3600),
					))
					acc.TestClient().SecurityIntegration.UpdateOauthForPartnerApplications(t, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithUnset(
						*sdk.NewOauthForPartnerApplicationsIntegrationUnsetRequest().
							WithEnabled(true).
							WithOauthUseSecondaryRoles(true),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_partner_applications.test", plancheck.ResourceActionUpdate),

						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "enabled", sdk.String("true"), sdk.String("false")),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "comment", sdk.String(comment), sdk.String("")),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone))),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectDrift("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", sdk.String("86400"), sdk.String("3600")),

						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "enabled", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "comment", tfjson.ActionUpdate, sdk.String(""), sdk.String(comment)),
						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", tfjson.ActionUpdate, sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)), sdk.String(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))),
						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectChange("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", tfjson.ActionUpdate, sdk.String("3600"), sdk.String("86400")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "3"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - TABLEAU_DESKTOP"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", "IMPLICIT"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", comment),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/basic_tableau"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_partner_applications.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "default"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "default"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "-1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - TABLEAU_DESKTOP"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", "NONE"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "7776000"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", ""),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_Complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	configVariables := config.Variables{
		"name":                         config.StringVariable(id.Name()),
		"oauth_client":                 config.StringVariable(string(sdk.OauthSecurityIntegrationClientTableauServer)),
		"blocked_roles_list":           config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
		"enabled":                      config.BoolVariable(true),
		"oauth_issue_refresh_tokens":   config.BoolVariable(false),
		"oauth_refresh_token_validity": config.IntegerVariable(86400),
		"oauth_use_secondary_roles":    config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
		"comment":                      config.StringVariable(comment),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete_tableau"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauServer)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "2"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.integration_type", "OAUTH - TABLEAU_SERVER"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.comment.0.value", comment),
					resource.TestCheckNoResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "describe_output.0.oauth_allowed_token_endpoints.0.value"),
				),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete_tableau"),
				ConfigVariables: configVariables,
				ResourceName:    "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauServer)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_issue_refresh_tokens", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_refresh_token_validity", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_roles_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_Invalid(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	invalidUseSecondaryRoles := config.Variables{
		"name":                      config.StringVariable(id.Name()),
		"oauth_client":              config.StringVariable(string(sdk.OauthSecurityIntegrationClientTableauDesktop)),
		"oauth_use_secondary_roles": config.StringVariable("invalid"),
		"blocked_roles_list":        config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
	}

	invalidOauthClient := config.Variables{
		"name":                      config.StringVariable(id.Name()),
		"oauth_client":              config.StringVariable("invalid"),
		"oauth_use_secondary_roles": config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/invalid"),
				ConfigVariables: invalidUseSecondaryRoles,
				ExpectError:     regexp.MustCompile(`Error: invalid OauthSecurityIntegrationUseSecondaryRolesOption: INVALID`),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/invalid"),
				ConfigVariables: invalidOauthClient,
				ExpectError:     regexp.MustCompile(`Error: invalid OauthSecurityIntegrationClientOption: INVALID`),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: oauthIntegrationForPartnerApplicationsBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "id", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   oauthIntegrationForPartnerApplicationsBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_WithQuotedName(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             oauthIntegrationForPartnerApplicationsBasicConfig(quotedId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "id", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   oauthIntegrationForPartnerApplicationsBasicConfig(quotedId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_partner_applications.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_oauth_integration_for_partner_applications.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "id", id.Name()),
				),
			},
		},
	})
}

func oauthIntegrationForPartnerApplicationsBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_oauth_integration_for_partner_applications" "test" {
  name               = "%s"
  oauth_client       = "LOOKER"
  oauth_redirect_uri = "https://example.com"
  blocked_roles_list = [ "ACCOUNTADMIN", "SECURITYADMIN" ]
}
`, name)
}
