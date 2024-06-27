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

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func(complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"enabled":             config.BoolVariable(true),
			"name":                config.StringVariable(id.Name()),
			"oauth_client_id":     config.StringVariable("foo"),
			"oauth_client_secret": config.StringVariable("foo"),
		}
		if complete {
			c["comment"] = config.StringVariable("foo")
			c["oauth_access_token_validity"] = config.IntegerVariable(42)
			c["oauth_authorization_endpoint"] = config.StringVariable("https://example.com")
			c["oauth_client_auth_method"] = config.StringVariable(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost))
			c["oauth_refresh_token_validity"] = config.IntegerVariable(12345)
			c["oauth_token_endpoint"] = config.StringVariable("https://example.com")
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_authorization_code_grant.test", "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/complete"),
				ConfigVariables: m(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_authorization_endpoint", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_auth_method", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_refresh_token_validity", "12345"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_token_endpoint", "https://example.com"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_authorization_code_grant.test", "created_on"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/basic"),
				ConfigVariables:         m(true),
				ResourceName:            "snowflake_api_authentication_integration_with_authorization_code_grant.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_access_token_validity", "0"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_authorization_endpoint", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_auth_method", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_refresh_token_validity", "0"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_token_endpoint", ""),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_authorization_code_grant.test", "created_on"),
				),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_complete(t *testing.T) {
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_authorization_endpoint", "https://example.com"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_auth_method", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_refresh_token_validity", "12345"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_authorization_code_grant.test", "oauth_token_endpoint", "https://example.com"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_authorization_code_grant.test", "created_on"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/complete"),
				ConfigVariables:         m(),
				ResourceName:            "snowflake_api_authentication_integration_with_authorization_code_grant.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant_invalidIncomplete(t *testing.T) {
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
			`The argument "oauth_client_secret" is required, but no definition was found.`,
			`The argument "oauth_client_id" is required, but no definition was found.`,
			`The argument "enabled" is required, but no definition was found.`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithAuthorizationCodeGrant/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}
