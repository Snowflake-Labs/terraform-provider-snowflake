package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

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
			"external_oauth_token_user_mapping_claims":            config.SetVariable(config.StringVariable("foo")),
			"name": config.StringVariable(id.Name()),
			"type": config.StringVariable(string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrl"),
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
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claims.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_token_user_mapping_claims.0", "foo"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "external_oauth_type", string(sdk.ExternalOauthSecurityIntegrationTypeCustom)),
					resource.TestCheckResourceAttrSet("snowflake_external_oauth_integration.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/completeWithJwsKeysUrl"),
				ConfigVariables:   m(),
				ResourceName:      "snowflake_external_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
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
			"external_oauth_any_role_mode":                        config.StringVariable("foo"),
			"external_oauth_audience_list":                        config.SetVariable(config.StringVariable("foo")),
			"external_oauth_blocked_roles_list":                   config.SetVariable(config.StringVariable("foo")),
			"external_oauth_issuer":                               config.StringVariable("foo"),
			"external_oauth_jws_keys_url":                         config.SetVariable(config.StringVariable("foo")),
			"external_oauth_rsa_public_key":                       config.StringVariable("foo"),
			"external_oauth_rsa_public_key_2":                     config.StringVariable("foo"),
			"external_oauth_scope_delimiter":                      config.StringVariable("foo"),
			"external_oauth_scope_mapping_attribute":              config.StringVariable("foo"),
			"external_oauth_snowflake_user_mapping_attribute":     config.StringVariable("foo"),
			"external_oauth_token_user_mapping_claims":            config.SetVariable(config.StringVariable("foo")),
			"name": config.StringVariable("foo"),
			"type": config.StringVariable("foo"),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ExternalOauthIntegration/complete"),
				ConfigVariables: m(),
				ExpectError:     helpers.MatchAllStringsInOrderNonOverlapping([]string{
					// TODO: Implement
				}),
			},
		},
	})
}
