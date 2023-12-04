package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ApiIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_API_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAcc_ApiIntegration")
	}

	apiIntNameAWS := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	apiIntNameAzure := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	apiIntNameGCP := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: apiIntegrationConfigAWS(apiIntNameAWS, []string{"https://123456.execute-api.us-west-2.amazonaws.com/prod/"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "name", apiIntNameAWS),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_provider", "aws_api_gateway"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "comment", "acceptance test"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_external_id"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_key"),
				),
			},
			{
				Config: apiIntegrationConfigAzure(apiIntNameAzure, []string{"https://apim-hello-world.azure-api.net/"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "name", apiIntNameAzure),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "api_provider", "azure_api_management"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_azure_int", "comment", "acceptance test"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "azure_multi_tenant_app_name"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "azure_consent_url"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_azure_int", "api_key"),
				),
			},
			{
				Config: apiIntegrationConfigGCP(apiIntNameGCP, []string{"https://gateway-id-123456.uc.gateway.dev/"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "name", apiIntNameGCP),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "api_provider", "google_api_gateway"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_gcp_int", "comment", "acceptance test"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "google_audience"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_gcp_int", "api_gcp_service_account"),
				),
			},
		},
	})
}

func apiIntegrationConfigAWS(name string, prefixes []string) string {
	return fmt.Sprintf(`
	resource "snowflake_api_integration" "test_aws_int" {
		name = "%s"
		api_provider = "aws_api_gateway"
		api_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
		api_allowed_prefixes = %q
		api_key = "12345"
		comment = "acceptance test"
		enabled = true
	}
	`, name, prefixes)
}

func apiIntegrationConfigAzure(name string, prefixes []string) string {
	return fmt.Sprintf(`
	resource "snowflake_api_integration" "test_azure_int" {
		name = "%s"
		api_provider = "azure_api_management"
		azure_tenant_id = "123456"
		azure_ad_application_id = "7890"
		api_allowed_prefixes = %q
		api_key = "12345"
		comment = "acceptance test"
		enabled = true
	}
	`, name, prefixes)
}

func apiIntegrationConfigGCP(name string, prefixes []string) string {
	return fmt.Sprintf(`
	resource "snowflake_api_integration" "test_gcp_int" {
		name = "%s"
		api_provider = "google_api_gateway"
		google_audience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
		api_allowed_prefixes = %q
		comment = "acceptance test"
		enabled = true
	}
	`, name, prefixes)
}
