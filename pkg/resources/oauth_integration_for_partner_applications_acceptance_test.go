package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OauthIntegrationForPartnerApplications_basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role1, role1Cleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(role1Cleanup)

	m := func(oauthClient string, complete bool, redirectUri *string) map[string]config.Variable {
		c := map[string]config.Variable{
			"name":         config.StringVariable(id.Name()),
			"oauth_client": config.StringVariable(oauthClient),
		}
		if complete {
			c["blocked_roles_list"] = config.SetVariable(config.StringVariable(role1.ID().Name()))
			c["comment"] = config.StringVariable("foo")
			c["enabled"] = config.BoolVariable(true)
			c["oauth_issue_refresh_tokens"] = config.StringVariable("false")
			c["oauth_refresh_token_validity"] = config.IntegerVariable(12345)
			c["oauth_use_secondary_roles"] = config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))
		}
		if redirectUri != nil {
			c["oauth_redirect_uri"] = config.StringVariable(*redirectUri)
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/basic"),
				ConfigVariables: m(string(sdk.OauthSecurityIntegrationClientTableauServer), false, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauServer)),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_allowed_authorization_endpoints.#"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_allowed_token_endpoints.#"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_authorization_endpoint"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "enabled"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete_tableau_server"),
				ConfigVariables: m(string(sdk.OauthSecurityIntegrationClientTableauServer), true, nil),
				Check: resource.ComposeTestCheckFunc(
					// TODO: proper check
					// resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "1"),
					// resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.0", role1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauServer)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "12345"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete_tableau_server"),
				ConfigVariables:   m(string(sdk.OauthSecurityIntegrationClientTableauServer), true, nil),
				ResourceName:      "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/basic"),
				ConfigVariables: m(string(sdk.OauthSecurityIntegrationClientTableauServer), false, nil),
				Check: resource.ComposeTestCheckFunc(
					// TODO: proper check
					// resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "0"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientTableauServer)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "7776000"),

					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "created_on"),
				),
			},
			// change client_type
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete_looker"),
				ConfigVariables: m(string(sdk.OauthSecurityIntegrationClientLooker), true, sdk.Pointer("https://example.com")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri", "https://example.com"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "created_on"),
				),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	validURL := "https://example.com"
	role1, role1Cleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(role1Cleanup)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"blocked_roles_list":           config.SetVariable(config.StringVariable(role1.ID().Name())),
			"comment":                      config.StringVariable("foo"),
			"enabled":                      config.BoolVariable(true),
			"name":                         config.StringVariable(id.Name()),
			"oauth_client":                 config.StringVariable(string(sdk.OauthSecurityIntegrationClientLooker)),
			"oauth_issue_refresh_tokens":   config.BoolVariable(true),
			"oauth_redirect_uri":           config.StringVariable(validURL),
			"oauth_refresh_token_validity": config.IntegerVariable(12345),
			"oauth_use_secondary_roles":    config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					// TODO: proper assert, also assert OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST
					// resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.#", "3"),
					// resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "blocked_roles_list.0", role1.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_client", string(sdk.OauthSecurityIntegrationClientLooker)),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_redirect_uri", validURL),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_refresh_token_validity", "12345"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration_for_partner_applications.test", "oauth_use_secondary_roles", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_allowed_authorization_endpoints.#"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_allowed_token_endpoints.#"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_authorization_endpoint"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_token_endpoint"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "oauth_client_id"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration_for_partner_applications.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete"),
				ConfigVariables:   m(),
				ResourceName:      "snowflake_oauth_integration_for_partner_applications.test",
				ImportState:       true,
				ImportStateVerify: true,
				// ignore because this field is not returned from snowflake
				ImportStateVerifyIgnore: []string{"oauth_redirect_uri"},
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_invalid(t *testing.T) {
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"blocked_roles_list":           config.SetVariable(config.StringVariable("foo")),
			"comment":                      config.StringVariable("foo"),
			"enabled":                      config.BoolVariable(true),
			"name":                         config.StringVariable("foo"),
			"oauth_client":                 config.StringVariable("invalid"),
			"oauth_issue_refresh_tokens":   config.BoolVariable(true),
			"oauth_redirect_uri":           config.StringVariable("foo"),
			"oauth_refresh_token_validity": config.IntegerVariable(1),
			"oauth_use_secondary_roles":    config.StringVariable("invalid"),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/complete"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_OauthIntegrationForPartnerApplications_InvalidIncomplete(t *testing.T) {
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
			`The argument "oauth_client" is required, but no definition was found.`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_OauthIntegrationForPartnerApplications/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}
