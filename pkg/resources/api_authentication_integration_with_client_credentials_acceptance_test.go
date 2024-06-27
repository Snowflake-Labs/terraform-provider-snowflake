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

func TestAcc_ApiAuthenticationIntegrationWithClientCredentials_basic(t *testing.T) {
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
			c["oauth_allowed_scopes"] = config.SetVariable(config.StringVariable("foo"))
			c["oauth_client_auth_method"] = config.StringVariable(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost))
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithClientCredentials/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_client_credentials.test", "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithClientCredentials/complete"),
				ConfigVariables: m(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_access_token_validity", "42"),
					// TODO: proper test
					// resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_allowed_scopes", "config.StringVariable("foo")"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_auth_method", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_token_endpoint", "https://example.com"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_client_credentials.test", "created_on"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithClientCredentials/basic"),
				ConfigVariables:         m(true),
				ResourceName:            "snowflake_api_authentication_integration_with_client_credentials.test",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithClientCredentials/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_access_token_validity", "0"),
					// TODO: proper test
					// resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_allowed_scopes", "config.StringVariable("foo")"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_auth_method", ""),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_token_endpoint", ""),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_client_credentials.test", "created_on"),
				),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithClientCredentials_complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                     config.StringVariable("foo"),
			"enabled":                     config.BoolVariable(true),
			"name":                        config.StringVariable(id.Name()),
			"oauth_access_token_validity": config.IntegerVariable(42),
			"oauth_allowed_scopes":        config.SetVariable(config.StringVariable("foo")),
			"oauth_client_auth_method":    config.StringVariable(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
			"oauth_client_id":             config.StringVariable("foo"),
			"oauth_client_secret":         config.StringVariable("foo"),
			"oauth_token_endpoint":        config.StringVariable("https://example.com"),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithClientCredentials/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_access_token_validity", "42"),
					// TODO: proper test
					// resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_allowed_scopes", "config.StringVariable("foo")"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_auth_method", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_client_credentials.test", "oauth_token_endpoint", "https://example.com"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_client_credentials.test", "created_on"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithClientCredentials/complete"),
				ConfigVariables:         m(),
				ResourceName:            "snowflake_api_authentication_integration_with_client_credentials.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_client_secret"},
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithClientCredentials_invalidIncomplete(t *testing.T) {
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
			`The argument "enabled" is required, but no definition was found.`,
			`The argument "oauth_client_id" is required, but no definition was found.`,
			`The argument "oauth_client_secret" is required, but no definition was found.`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithClientCredentials/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}
