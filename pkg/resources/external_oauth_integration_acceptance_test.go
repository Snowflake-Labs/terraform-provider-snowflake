package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ExternalOauthIntegration(t *testing.T) {
	oauthIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	integrationType := "AZURE"

	issuer := fmt.Sprintf("https://sts.windows.net/%s", uuid.NewString())

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalOauthIntegrationConfig(oauthIntName, integrationType, issuer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "name", oauthIntName),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "type", integrationType),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "issuer", issuer),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "snowflake_user_mapping_attribute", "LOGIN_NAME"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.#", "1"),
					resource.TestCheckResourceAttr("snowflake_external_oauth_integration.test", "token_user_mapping_claims.0", "upn"),
				),
			},
			{
				ResourceName:      "snowflake_external_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func externalOauthIntegrationConfig(name, integrationType, issuer string) string {
	return fmt.Sprintf(`
	resource "snowflake_external_oauth_integration" "test" {
		name = "%s"
		type = "%s"
		enabled = true
  		issuer = "%s"
  		snowflake_user_mapping_attribute = "LOGIN_NAME"
		jws_keys_urls = ["https://login.windows.net/common/discovery/keys"]
		audience_urls = ["https://analysis.windows.net/powerbi/connector/Snowflake"]
  		token_user_mapping_claims = ["upn"]
	}
	`, name, integrationType, issuer)
}
