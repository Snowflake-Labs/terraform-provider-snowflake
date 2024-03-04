package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OAuthIntegration(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	oauthClient := "CUSTOM"
	clientType := "PUBLIC"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: oauthIntegrationConfig(name, oauthClient, clientType, "SYSADMIN"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "oauth_client", oauthClient),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "oauth_client_type", clientType),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "oauth_issue_refresh_tokens", "true"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "oauth_refresh_token_validity", "3600"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "blocked_roles_list.0", "SYSADMIN"),
				),
			},
			{
				// role change proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2358 issue
				Config: oauthIntegrationConfig(name, oauthClient, clientType, "USERADMIN"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "blocked_roles_list.#", "1"),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "blocked_roles_list.0", "USERADMIN"),
				),
			},
			{
				ResourceName:      "snowflake_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func oauthIntegrationConfig(name, oauthClient, clientType string, blockedRole string) string {
	return fmt.Sprintf(`
	resource "snowflake_oauth_integration" "test" {
		name                         = "%s"
		oauth_client                 = "%s"
		oauth_client_type            = "%s"
		oauth_redirect_uri           = "https://www.example.com/oauth2/callback"
		enabled                      = true
  		oauth_issue_refresh_tokens   = true
  		oauth_refresh_token_validity = 3600
  		blocked_roles_list           = ["%s"]
	}
	`, name, oauthClient, clientType, blockedRole)
}

func TestAcc_OAuthIntegrationTableau(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	oauthClient := "TABLEAU_DESKTOP"
	clientType := "PUBLIC" // not used, but left to fail the test

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: oauthIntegrationConfigTableau(name, oauthClient, clientType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "oauth_client", oauthClient),
					// resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "oauth_client_type", clientType),
				),
			},
			{
				ResourceName:      "snowflake_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func oauthIntegrationConfigTableau(name, oauthClient, clientType string) string {
	return fmt.Sprintf(`
	resource "snowflake_oauth_integration" "test" {
		name                         = "%s"
		oauth_client                 = "%s"
	#	oauth_client_type            = "%s" # this cannot be set for TABLEAU
		enabled                      = true
        oauth_refresh_token_validity = 36000
        oauth_issue_refresh_tokens   = true
        blocked_roles_list           = ["SYSADMIN"]
	}
	`, name, oauthClient, clientType)
}
