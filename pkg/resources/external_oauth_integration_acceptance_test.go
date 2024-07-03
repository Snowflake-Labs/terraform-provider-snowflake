package resources_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalOauthIntegration_basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role1, role1Cleanup := acc.TestClient().Role.CreateRole(t)
	issuer := random.String()
	t.Cleanup(role1Cleanup)
	m := func(complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"enabled":             config.BoolVariable(true),
			"name":                config.StringVariable(id.Name()),
			"external_oauth_type": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(issuer),
			"external_oauth_jws_keys_url":                     config.SetVariable(config.StringVariable("https://example.com")),
		}
		if complete {
			c["external_oauth_add_privileged_roles_to_blocked_list"] = config.BoolVariable(true)
			c["external_oauth_allowed_roles_list"] = config.SetVariable(config.StringVariable(role1.ID().Name()))
			c["external_oauth_any_role_mode"] = config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable))
			c["external_oauth_audience_list"] = config.SetVariable(config.StringVariable("foo"))
			c["external_oauth_scope_delimiter"] = config.StringVariable(".")
			c["external_oauth_scope_mapping_attribute"] = config.StringVariable("foo")
			c["comment"] = config.StringVariable("foo")
		}
		return c
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.0", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", ","),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "")),
			},
			// import - without optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/basic"),
				ConfigVariables: m(false),
				ResourceName:    "snowflake_external_oauth_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_token_user_mapping_claim.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_token_user_mapping_claim.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_jws_keys_url.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_jws_keys_url.0", "https://example.com"),
				),
			},
			// set optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.0", role1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.0", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", role1.Name),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo")),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(true),
				ResourceName:    "snowflake_external_oauth_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_token_user_mapping_claim.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_token_user_mapping_claim.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_jws_keys_url.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_jws_keys_url.0", "https://example.com"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_allowed_roles_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_allowed_roles_list.0", role1.Name),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_audience_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_audience_list.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "external_oauth_scope_delimiter", "."),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", "foo"),
				),
			},
			// change values externally
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(true),
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateExternalOauth(t, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).
						WithSet(*sdk.NewExternalOauthIntegrationSetRequest().
							WithExternalOauthSnowflakeUserMappingAttribute(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName),
						))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", sdk.String(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)), sdk.String(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName))),

						planchecks.ExpectChange("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", tfjson.ActionUpdate, sdk.String(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)), sdk.String(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress))),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.0", role1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.0", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", role1.Name),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo")),
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckNoResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.#"),
					resource.TestCheckNoResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode"),
					resource.TestCheckNoResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.#"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.0", "https://example.com"),
					resource.TestCheckNoResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter"),
					resource.TestCheckNoResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_mapping_attribute"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", ","),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "")),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_completeWithJwsKeysUrlAndAllowedRolesList(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role1, role1Cleanup := acc.TestClient().Role.CreateRole(t)
	issuer := random.String()
	t.Cleanup(role1Cleanup)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment": config.StringVariable("foo"),
			"enabled": config.BoolVariable(true),
			"external_oauth_add_privileged_roles_to_blocked_list": config.BoolVariable(true),
			"external_oauth_allowed_roles_list":                   config.SetVariable(config.StringVariable(role1.ID().Name())),
			"external_oauth_any_role_mode":                        config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                        config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                               config.StringVariable(issuer),
			"external_oauth_jws_keys_url":                         config.SetVariable(config.StringVariable("https://example.com")),
			"external_oauth_scope_delimiter":                      config.StringVariable("."),
			"external_oauth_scope_mapping_attribute":              config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute":     config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":             config.SetVariable(config.StringVariable("foo")),
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.0", role1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.0", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", role1.Name),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables:         m(),
				ResourceName:            "snowflake_external_oauth_integration.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2", "external_oauth_scope_mapping_attribute"},
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_completeWithRsaPublicKeysAndBlockedRolesList_paramSet(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role1, role1Cleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(role1Cleanup)
	expectedRoles := []string{"ACCOUNTADMIN", "SECURITYADMIN", role1.Name}
	sort.Strings(expectedRoles)
	issuer := random.String()
	rsaKey := random.GenerateRSAPublicKey(t)
	paramCleanup := acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList, "true")
	t.Cleanup(paramCleanup)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment": config.StringVariable("foo"),
			"enabled": config.BoolVariable(true),
			"external_oauth_add_privileged_roles_to_blocked_list": config.BoolVariable(true),
			"external_oauth_blocked_roles_list":                   config.SetVariable(config.StringVariable(role1.ID().Name())),
			"external_oauth_any_role_mode":                        config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                        config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                               config.StringVariable(issuer),
			"external_oauth_rsa_public_key":                       config.StringVariable(rsaKey),
			"external_oauth_rsa_public_key_2":                     config.StringVariable(rsaKey),
			"external_oauth_scope_delimiter":                      config.StringVariable("."),
			"external_oauth_scope_mapping_attribute":              config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute":     config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":             config.SetVariable(config.StringVariable("foo")),
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.0", role1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_rsa_public_key", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_rsa_public_key.0.value", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_rsa_public_key_2.0.value", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_blocked_roles_list.0.value", strings.Join(expectedRoles, ",")),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables:         m(),
				ResourceName:            "snowflake_external_oauth_integration.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2", "external_oauth_scope_mapping_attribute"},
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_completeWithRsaPublicKeysAndBlockedRolesList_paramUnset(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role1, role1Cleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(role1Cleanup)
	issuer := random.String()
	rsaKey := random.GenerateRSAPublicKey(t)
	paramCleanup := acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList, "false")
	t.Cleanup(paramCleanup)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment": config.StringVariable("foo"),
			"enabled": config.BoolVariable(true),
			"external_oauth_add_privileged_roles_to_blocked_list": config.BoolVariable(true),
			"external_oauth_blocked_roles_list":                   config.SetVariable(config.StringVariable(role1.ID().Name())),
			"external_oauth_any_role_mode":                        config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                        config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                               config.StringVariable(issuer),
			"external_oauth_rsa_public_key":                       config.StringVariable(rsaKey),
			"external_oauth_rsa_public_key_2":                     config.StringVariable(rsaKey),
			"external_oauth_scope_delimiter":                      config.StringVariable("."),
			"external_oauth_scope_mapping_attribute":              config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute":     config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":             config.SetVariable(config.StringVariable("foo")),
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.0", role1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_rsa_public_key", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_jws_keys_url.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_rsa_public_key.0.value", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_rsa_public_key_2.0.value", rsaKey),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_blocked_roles_list.0.value", role1.Name),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables:         m(),
				ResourceName:            "snowflake_external_oauth_integration.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2", "external_oauth_scope_mapping_attribute"},
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_invalid(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment": config.StringVariable("foo"),
			"enabled": config.BoolVariable(true),
			"external_oauth_add_privileged_roles_to_blocked_list": config.BoolVariable(true),
			"external_oauth_allowed_roles_list":                   config.SetVariable(config.StringVariable("foo")),
			"external_oauth_any_role_mode":                        config.StringVariable("invalid"),
			"external_oauth_audience_list":                        config.SetVariable(config.StringVariable("foo")),
			"external_oauth_blocked_roles_list":                   config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                               config.StringVariable("foo"),
			"external_oauth_jws_keys_url":                         config.SetVariable(config.StringVariable("foo")),
			"external_oauth_rsa_public_key":                       config.StringVariable("foo"),
			"external_oauth_rsa_public_key_2":                     config.StringVariable("foo"),
			"external_oauth_scope_delimiter":                      config.StringVariable("foo"),
			"external_oauth_scope_mapping_attribute":              config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute":     config.StringVariable("invalid"),
			"external_oauth_token_user_mapping_claim":             config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable("foo"),
			"external_oauth_type": config.StringVariable("invalid"),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(),
				ExpectError: helpers.MatchAllStringsInOrderNonOverlapping([]string{
					fmt.Sprintf("expected external_oauth_any_role_mode to be one of %q, got invalid", sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationAnyRoleModes)),
					fmt.Sprintf("expected external_oauth_snowflake_user_mapping_attribute to be one of %q, got invalid", sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationSnowflakeUserMappingAttributes)),
					fmt.Sprintf("expected external_oauth_type to be one of %q, got invalid", sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationTypes)),
				}),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_InvalidIncomplete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name": config.StringVariable(id.Name()),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			// Some strings are trimmed because of inconsistent '\n' placement from tf error messages.
			`The argument "external_oauth_type" is required, but no definition was found.`,
			`The argument "external_oauth_snowflake_user_mapping_attribute" is required,`,
			`The argument "enabled" is required, but no definition was found.`,
			`The argument "external_oauth_issuer" is required,`,
			`The argument "external_oauth_token_user_mapping_claim" is required,`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}
