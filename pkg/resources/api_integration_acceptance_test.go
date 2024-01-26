package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	dummyAwsPrefix            = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	dummyAwsOtherPrefix       = "https://123456.execute-api.us-west-2.amazonaws.com/prod/"
	dummyAzurePrefix          = "https://apim-hello-world.azure-api.net/dev"
	dummyAzureOtherPrefix     = "https://apim-hello-world.azure-api.net/prod"
	dummyGooglePrefix         = "https://gateway-id-123456.uc.gateway.dev/prod"
	dummyGoogleOtherPrefix    = "https://gateway-id-123456.uc.gateway.dev/dev"
	dummyApiAwsRoleArn        = "arn:aws:iam::000000000001:/role/test"
	dummyAzureTenantId        = "00000000-0000-0000-0000-000000000000"
	dummyAzureAdApplicationId = "11111111-1111-1111-1111-111111111111"
	dummyGoogleAudience       = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
)

func TestAcc_ApiIntegration(t *testing.T) {
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

func TestAcc_ApiIntegration_aws(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := "acceptance test"
	key := "12345"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":             config.StringVariable(name),
			"api_provider":     config.StringVariable("aws_api_gateway"),
			"api_aws_role_arn": config.StringVariable(dummyApiAwsRoleArn),
			"api_allowed_prefixes": config.ListVariable(
				config.StringVariable(dummyAwsPrefix),
			),
			"api_blocked_prefixes": config.ListVariable(
				config.StringVariable(dummyAwsOtherPrefix),
			),
			"api_key": config.StringVariable("12345"),
			"comment": config.StringVariable(comment),
			"enabled": config.BoolVariable(true),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckApiIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "name", name),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_provider", "aws_api_gateway"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_aws_role_arn", dummyApiAwsRoleArn),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_allowed_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_allowed_prefixes.0", dummyAwsPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_blocked_prefixes.#", "1"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_blocked_prefixes.0", dummyAwsOtherPrefix),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_iam_user_arn"),
					resource.TestCheckResourceAttrSet("snowflake_api_integration.test_aws_int", "api_aws_external_id"),
					resource.TestCheckResourceAttr("snowflake_api_integration.test_aws_int", "api_key", key),
				),
			},
			// IMPORT
			{
				ConfigVariables:         m(),
				ResourceName:            "snowflake_api_integration.test_aws_int",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func testAccCheckApiIntegrationDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_api_integration" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewAccountObjectIdentifier(rs.Primary.Attributes["name"])
		existingApiIntegration, err := client.ApiIntegrations.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("api integration %v still exists", existingApiIntegration.ID().FullyQualifiedName())
		}
	}
	return nil
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
