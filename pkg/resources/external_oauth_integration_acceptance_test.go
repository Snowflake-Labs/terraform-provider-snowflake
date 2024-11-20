package resources_test

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

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
	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	issuer := random.String()
	t.Cleanup(roleCleanup)
	m := func(complete, unset bool) map[string]config.Variable {
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
			c["external_oauth_allowed_roles_list"] = config.SetVariable(config.StringVariable(role.ID().Name()))
			c["external_oauth_any_role_mode"] = config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable))
			c["external_oauth_audience_list"] = config.SetVariable(config.StringVariable("foo"))
			c["external_oauth_scope_delimiter"] = config.StringVariable(".")
			c["external_oauth_scope_mapping_attribute"] = config.StringVariable("foo")
			c["comment"] = config.StringVariable("foo")
		}
		if unset {
			c["external_oauth_scope_mapping_attribute"] = config.StringVariable("foo")
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
				ConfigVariables: m(false, false),
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
				ConfigVariables: m(false, false),
				ResourceName:    "snowflake_external_oauth_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_token_user_mapping_claim.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_token_user_mapping_claim.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_jws_keys_url.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_jws_keys_url.0", "https://example.com"),
				),
			},
			// set optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.0", role.ID().Name()),
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
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo")),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(true, true),
				ResourceName:    "snowflake_external_oauth_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_token_user_mapping_claim.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_token_user_mapping_claim.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_jws_keys_url.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_jws_keys_url.0", "https://example.com"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_allowed_roles_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_allowed_roles_list.0", role.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_audience_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_audience_list.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_scope_delimiter", "."),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", "foo"),
				),
			},
			// change values externally
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(true, true),
				PreConfig: func() {
					acc.TestClient().SecurityIntegration.UpdateExternalOauth(t, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).
						WithSet(*sdk.NewExternalOauthIntegrationSetRequest().
							WithExternalOauthSnowflakeUserMappingAttribute(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName),
						))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_oauth_integration.test", plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", sdk.String(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)), sdk.String(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName))),
						planchecks.ExpectChange("snowflake_external_oauth_integration.test", "external_oauth_snowflake_user_mapping_attribute", tfjson.ActionUpdate, sdk.String(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)), sdk.String(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress))),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.0", role.ID().Name()),
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
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", role.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_audience_list.0.value", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_token_user_mapping_claim.0.value", "['foo']"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.comment.0.value", "foo")),
			},
			// unset without force new
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/unset"),
				ConfigVariables: m(false, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_oauth_integration.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckNoResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.#"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_any_role_mode", ""),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_jws_keys_url.0", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_scope_delimiter", ""),
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
			// unset all
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/basic"),
				ConfigVariables: m(false, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_oauth_integration.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrlAndAllowedRolesList"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_allowed_roles_list.0", role.ID().Name()),
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
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
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
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_allowed_roles_list.0.value", role.ID().Name()),
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
	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	expectedRoles := []string{"ACCOUNTADMIN", "SECURITYADMIN", role.ID().Name()}
	sort.Strings(expectedRoles)
	issuer := random.String()
	rsaKey, _ := random.GenerateRSAPublicKey(t)
	paramCleanup := acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList, "true")
	t.Cleanup(paramCleanup)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable(role.ID().Name())),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(issuer),
			"external_oauth_rsa_public_key":                   config.StringVariable(rsaKey),
			"external_oauth_rsa_public_key_2":                 config.StringVariable(rsaKey),
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.0", role.ID().Name()),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				ResourceName:    "snowflake_external_oauth_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_blocked_roles_list.#", "3"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_blocked_roles_list.0", expectedRoles[0]),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_blocked_roles_list.1", expectedRoles[1]),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_blocked_roles_list.2", expectedRoles[2]),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_audience_list.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_audience_list.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_issuer", issuer),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_rsa_public_key", rsaKey),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_scope_delimiter", "."),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_token_user_mapping_claim.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_token_user_mapping_claim.0", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
				),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_completeWithRsaPublicKeysAndBlockedRolesList_paramUnset(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	issuer := random.String()
	rsaKey, _ := random.GenerateRSAPublicKey(t)
	paramCleanup := acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList, "false")
	t.Cleanup(paramCleanup)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable(role.ID().Name())),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable(issuer),
			"external_oauth_rsa_public_key":                   config.StringVariable(rsaKey),
			"external_oauth_rsa_public_key_2":                 config.StringVariable(rsaKey),
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
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithRsaPublicKeysAndBlockedRolesList"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_blocked_roles_list.0", role.ID().Name()),
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
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "describe_output.0.external_oauth_blocked_roles_list.0.value", role.ID().Name()),
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

func TestAcc_ExternalOauthIntegration_invalidAnyRoleMode(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_allowed_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_any_role_mode":                    config.StringVariable("invalid"),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable("foo"),
			"external_oauth_jws_keys_url":                     config.SetVariable(config.StringVariable("foo")),
			"external_oauth_rsa_public_key":                   config.StringVariable("foo"),
			"external_oauth_rsa_public_key_2":                 config.StringVariable("foo"),
			"external_oauth_scope_delimiter":                  config.StringVariable("foo"),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable("foo"),
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
				ExpectError:     regexp.MustCompile("Error: invalid ExternalOauthSecurityIntegrationAnyRoleModeOption: INVALID"),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_invalidSnowflakeUserMappingAttribute(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_allowed_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable("foo"),
			"external_oauth_jws_keys_url":                     config.SetVariable(config.StringVariable("foo")),
			"external_oauth_rsa_public_key":                   config.StringVariable("foo"),
			"external_oauth_rsa_public_key_2":                 config.StringVariable("foo"),
			"external_oauth_scope_delimiter":                  config.StringVariable("foo"),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable("invalid"),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
			"name":                config.StringVariable("foo"),
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
				ExpectError:     regexp.MustCompile("Error: invalid ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption: INVALID"),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_invalidOauthType(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                                         config.StringVariable("foo"),
			"enabled":                                         config.BoolVariable(true),
			"external_oauth_allowed_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_any_role_mode":                    config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
			"external_oauth_audience_list":                    config.SetVariable(config.StringVariable("foo")),
			"external_oauth_blocked_roles_list":               config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                           config.StringVariable("foo"),
			"external_oauth_jws_keys_url":                     config.SetVariable(config.StringVariable("foo")),
			"external_oauth_rsa_public_key":                   config.StringVariable("foo"),
			"external_oauth_rsa_public_key_2":                 config.StringVariable("foo"),
			"external_oauth_scope_delimiter":                  config.StringVariable("foo"),
			"external_oauth_scope_mapping_attribute":          config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
			"external_oauth_token_user_mapping_claim":         config.SetVariable(config.StringVariable("foo")),
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
				ExpectError:     regexp.MustCompile("Error: invalid ExternalOauthSecurityIntegrationTypeOption: INVALID"),
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

func TestAcc_ExternalOauthIntegration_migrateFromVersion092_withRsaPublicKeysAndBlockedRolesList(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	issuer := random.String()
	rsaKey, _ := random.GenerateRSAPublicKey(t)
	resourceName := "snowflake_external_oauth_integration.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: externalOauthIntegrationWithRsaPublicKeysAndBlockedRolesListv092(id.Name(), issuer, rsaKey, role.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttr(resourceName, "issuer", issuer),
					resource.TestCheckResourceAttr(resourceName, "token_user_mapping_claims.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_user_mapping_claims.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)),
					resource.TestCheckResourceAttr(resourceName, "scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr(resourceName, "rsa_public_key", rsaKey),
					resource.TestCheckResourceAttr(resourceName, "rsa_public_key_2", rsaKey),
					resource.TestCheckResourceAttr(resourceName, "blocked_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "blocked_roles.0", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "audience_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "audience_urls.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(resourceName, "scope_delimiter", ":"),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   externalOauthIntegrationWithRsaPublicKeysAndBlockedRolesListv093(id.Name(), issuer, rsaKey, role.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_rsa_public_key", rsaKey),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_rsa_public_key_2", rsaKey),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_blocked_roles_list.0", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_scope_delimiter", ":"),
				),
			},
		},
	})
}

func externalOauthIntegrationWithRsaPublicKeysAndBlockedRolesListv092(name, issuer, rsaKey, roleName string) string {
	s := `
locals {
  key_raw = <<-EOT
%s
  EOT
  key = trimsuffix(local.key_raw, "\n")
}
resource "snowflake_external_oauth_integration" "test" {
	name                             = "%s"
	enabled = true
	type = "CUSTOM"
	issuer = "%s"
	token_user_mapping_claims = ["foo"]
	snowflake_user_mapping_attribute = "LOGIN_NAME"
	scope_mapping_attribute = "foo"
	rsa_public_key = local.key
	rsa_public_key_2 = local.key
	blocked_roles = ["%s"]
	audience_urls = ["foo"]
	any_role_mode = "DISABLE"
	scope_delimiter = ":"
}`
	return fmt.Sprintf(s, rsaKey, name, issuer, roleName)
}

func externalOauthIntegrationWithRsaPublicKeysAndBlockedRolesListv093(name, issuer, rsaKey, roleName string) string {
	s := `
locals {
  key_raw = <<-EOT
%s
  EOT
  key = trimsuffix(local.key_raw, "\n")
}
resource "snowflake_external_oauth_integration" "test" {
	name                             = "%s"
	enabled = true
	external_oauth_type = "CUSTOM"
	external_oauth_issuer = "%s"
	external_oauth_token_user_mapping_claim = ["foo"]
	external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
	external_oauth_scope_mapping_attribute = "foo"
	external_oauth_rsa_public_key = local.key
	external_oauth_rsa_public_key_2 = local.key
	external_oauth_blocked_roles_list = ["%s"]
	external_oauth_audience_list = ["foo"]
	external_oauth_any_role_mode = "DISABLE"
	external_oauth_scope_delimiter = ":"
}`
	return fmt.Sprintf(s, rsaKey, name, issuer, roleName)
}

func TestAcc_ExternalOauthIntegration_migrateFromVersion092_withJwsKeysUrlAndAllowedRolesList(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	issuer := random.String()
	resourceName := "snowflake_external_oauth_integration.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: externalOauthIntegrationWithJwsKeysUrlAndAllowedRolesListv092(id.Name(), issuer, role.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttr(resourceName, "issuer", issuer),
					resource.TestCheckResourceAttr(resourceName, "token_user_mapping_claims.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_user_mapping_claims.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)),
					resource.TestCheckResourceAttr(resourceName, "scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr(resourceName, "jws_keys_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "jws_keys_urls.0", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "allowed_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allowed_roles.0", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "audience_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "audience_urls.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(resourceName, "scope_delimiter", ":"),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   externalOauthIntegrationWithJwsKeysUrlAndAllowedRolesListv093(id.Name(), issuer, role.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_issuer", issuer),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_token_user_mapping_claim.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_token_user_mapping_claim.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_snowflake_user_mapping_attribute", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_scope_mapping_attribute", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_jws_keys_url.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_jws_keys_url.0", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_allowed_roles_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_allowed_roles_list.0", role.ID().Name()),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_audience_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_audience_list.0", "foo"),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_any_role_mode", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(resourceName, "external_oauth_scope_delimiter", ":"),
				),
			},
		},
	})
}

func externalOauthIntegrationWithJwsKeysUrlAndAllowedRolesListv092(name, issuer, roleName string) string {
	s := `
resource "snowflake_external_oauth_integration" "test" {
	name                             = "%s"
	enabled = true
	type = "CUSTOM"
	issuer = "%s"
	token_user_mapping_claims = ["foo"]
	snowflake_user_mapping_attribute = "LOGIN_NAME"
	scope_mapping_attribute = "foo"
	jws_keys_urls = ["https://example.com"]
	allowed_roles = ["%s"]
	audience_urls = ["foo"]
	any_role_mode = "DISABLE"
	scope_delimiter = ":"
}`
	return fmt.Sprintf(s, name, issuer, roleName)
}

func externalOauthIntegrationWithJwsKeysUrlAndAllowedRolesListv093(name, issuer, roleName string) string {
	s := `
resource "snowflake_external_oauth_integration" "test" {
	name                             = "%s"
	enabled = true
	external_oauth_type = "CUSTOM"
	external_oauth_issuer = "%s"
	external_oauth_token_user_mapping_claim = ["foo"]
	external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
	external_oauth_scope_mapping_attribute = "foo"
	external_oauth_jws_keys_url = ["https://example.com"]
	external_oauth_allowed_roles_list = ["%s"]
	external_oauth_audience_list = ["foo"]
	external_oauth_any_role_mode = "DISABLE"
	external_oauth_scope_delimiter = ":"
}`
	return fmt.Sprintf(s, name, issuer, roleName)
}

func TestAcc_ExternalOauthIntegration_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: externalOauthIntegrationBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   externalOauthIntegrationBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_ExternalOauthIntegration_WithQuotedName(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             externalOauthIntegrationBasicConfig(quotedId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   externalOauthIntegrationBasicConfig(quotedId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_oauth_integration.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_external_oauth_integration.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "id", id.Name()),
				),
			},
		},
	})
}

func externalOauthIntegrationBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_external_oauth_integration" "test" {
	name               = "%s"
	external_oauth_type                             = "CUSTOM"
	enabled                                         = true
	external_oauth_issuer                           = "issuer"
	external_oauth_token_user_mapping_claim         = [ "foo" ]
	external_oauth_snowflake_user_mapping_attribute = "EMAIL_ADDRESS"
	external_oauth_jws_keys_url                     = [ "https://example.com" ]
}
`, name)
}
