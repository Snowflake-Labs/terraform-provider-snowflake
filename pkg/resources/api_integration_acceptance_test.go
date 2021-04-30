package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ApiIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_API_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAccApiIntegration")
	}

	apiIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	apiIntName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: apiIntegrationConfig_aws(apiIntName, []string{"https://123456.execute-api.us-west-2.amazonaws.com/prod/"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "name", apiIntName),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_provider", "aws_api_gateway"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_external_id"),
				),
			},
			{
				Config: apiIntegrationConfig_azure(apiIntName2, []string{"https://apim-hello-world.azure-api.net/"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "name", apiIntName2),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_provider", "azure_api_management"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "azure_multi_tenant_app_name"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "azure_consent_url"),
				),
			},
		},
	})
}

func apiIntegrationConfig_aws(name string, prefixes []string) string {
	return fmt.Sprintf(`
	resource "snowflake_api_integration" "test_aws_int" {
		name = "%s"
		api_provider = "aws_api_gateway"
		api_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
		api_allowed_prefixes = %q
		enabled = true
	}
	`, name, prefixes)
}

func apiIntegrationConfig_azure(name string, prefixes []string) string {
	return fmt.Sprintf(`
	resource "snowflake_api_integration" "test_azure_int" {
		name = "%s"
		api_provider = "azure_api_management"
		azure_tenant_id = "123456"
		azure_ad_application_id = "7890"
		api_allowed_prefixes = %q
		enabled = true
	}
	`, name, prefixes)
}
