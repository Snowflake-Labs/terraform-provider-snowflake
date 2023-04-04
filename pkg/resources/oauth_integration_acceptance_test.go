package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_OAuthIntegration(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	oauthClient := "CUSTOM"
	clientType := "PUBLIC"

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: oauthIntegrationConfig(name, oauthClient, clientType),
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
				ResourceName:      "snowflake_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func oauthIntegrationConfig(name, oauthClient, clientType string) string {
	return fmt.Sprintf(`
	resource "snowflake_oauth_integration" "test" {
		name                         = "%s"
		oauth_client                 = "%s"
		oauth_client_type            = "%s"
		oauth_redirect_uri           = "https://www.example.com/oauth2/callback"
		enabled                      = true
  		oauth_issue_refresh_tokens   = true
  		oauth_refresh_token_validity = 3600
  		blocked_roles_list           = ["SYSADMIN"]
	}
	`, name, oauthClient, clientType)
}
