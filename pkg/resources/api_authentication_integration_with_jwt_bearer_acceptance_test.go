package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_basic(t *testing.T) {
	// TODO [SNOW-1452191]: unskip
	t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")
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
			c["oauth_authorization_endpoint"] = config.StringVariable("foo")
			c["oauth_client_auth_method"] = config.StringVariable("foo")
			c["oauth_refresh_token_validity"] = config.IntegerVariable(42)
			c["oauth_token_endpoint"] = config.StringVariable("foo")
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
				ConfigVariables: m(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_authorization_endpoint", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_auth_method", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_refresh_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_token_endpoint", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
				ConfigVariables:   m(true),
				ResourceName:      "snowflake_api_authentication_integration_with_jwt_bearer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_authorization_endpoint", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_auth_method", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_refresh_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_token_endpoint", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on"),
				),
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_complete(t *testing.T) {
	// TODO [SNOW-1452191]: unskip
	t.Skip("Skip because of the error: Invalid value specified for property 'OAUTH_CLIENT_SECRET'")
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"comment":                      config.StringVariable("foo"),
			"enabled":                      config.BoolVariable(true),
			"name":                         config.StringVariable(id.Name()),
			"oauth_access_token_validity":  config.IntegerVariable(42),
			"oauth_authorization_endpoint": config.StringVariable("foo"),
			"oauth_client_auth_method":     config.StringVariable("foo"),
			"oauth_client_id":              config.StringVariable("foo"),
			"oauth_client_secret":          config.StringVariable("foo"),
			"oauth_refresh_token_validity": config.IntegerVariable(42),
			"oauth_token_endpoint":         config.StringVariable("foo"),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_access_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_authorization_endpoint", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_auth_method", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_id", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_client_secret", "foo"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_refresh_token_validity", "42"),
					resource.TestCheckResourceAttr("snowflake_api_authentication_integration_with_jwt_bearer.test", "oauth_token_endpoint", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_api_authentication_integration_with_jwt_bearer.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/complete"),
				ConfigVariables:   m(),
				ResourceName:      "snowflake_api_authentication_integration_with_jwt_bearer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ApiAuthenticationIntegrationWithJwtBearer_invalidIncomplete(t *testing.T) {
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ApiAuthenticationIntegrationWithJwtBearer/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}
