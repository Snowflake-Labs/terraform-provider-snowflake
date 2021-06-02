package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSystemGenerateSCIMAccessToken(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: generateAccessTokenConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_system_generate_scim_access_token.p", "integration_name", "AAD_PROVISIONING"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_generate_scim_access_token.p", "access_token"),
				),
			},
		},
	})
}

func generateAccessTokenConfig() string {
	s := `
	resource "snowflake_role" "azure" {
		name = "AAD_PROVISIONER"
		comment = "test comment"
	}

	resource "snowflake_account_grant" "azurecu" {
		roles     = [snowflake_role.azure.name]
		privilege = "CREATE USER"
	}
	resource "snowflake_account_grant" "azurecr" {
		roles     = [snowflake_role.azure.name]
		privilege = "CREATE ROLE"
	}
	resource "snowflake_role_grants" "azure" {
		role_name = snowflake_role.azure.name
		roles = ["ACCOUNTADMIN"]
	}

	resource "snowflake_scim_integration" "azure" {
		name = "AAD_PROVISIONING"
		scim_client = "AZURE"
		provisioner_role = snowflake_role.azure.name
		depends_on = [
			snowflake_account_grant.azurecu,
			snowflake_account_grant.azurecr,
			snowflake_role_grants.azure
		]
	}

	data snowflake_system_generate_scim_access_token p {
		integration_name = snowflake_scim_integration.azure.name
	}
	`
	return s
}
