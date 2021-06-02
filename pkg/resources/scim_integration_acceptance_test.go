package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccScimIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SCIM_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAccScimIntegration")
	}

	scimIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	scimIntName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: scimIntegrationConfig_azure(scimIntName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "name", scimIntName),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "scim_client", "AZURE"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "provisioner_role", "AAD_PROVISIONER"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test_azure_int", "created_on"),
				),
			},
			{
				Config: scimIntegrationConfig_okta_np(scimIntName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_okta_int_np", "name", scimIntName2),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_okta_int_np", "scim_client", "OKTA"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_okta_int_np", "provisioner_role", "OKTA_PROVISIONER"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_okta_int_np", "network_policy", "OKTA_NETWORK_POLICY"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test_okta_int_np", "created_on"),
				),
			},
		},
	})
}

func scimIntegrationConfig_azure(name string) string {
	return fmt.Sprintf(`
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

	resource "snowflake_scim_integration" "test_azure_int" {
		name = "%s"
		scim_client = "AZURE"
		provisioner_role = snowflake_role.azure.name
		depends_on = [
			snowflake_account_grant.azurecu,
			snowflake_account_grant.azurecr,
			snowflake_role_grants.azure
		]
	}
	`, name)
}

func scimIntegrationConfig_okta_np(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "okta" {
		name = "OKTA_PROVISIONER"
		comment = "test comment"
	}

	resource "snowflake_account_grant" "oktacu" {
		roles     = [snowflake_role.okta.name]
		privilege = "CREATE USER"
	}

	resource "snowflake_account_grant" "oktacr" {
		roles     = [snowflake_role.okta.name]
		privilege = "CREATE ROLE"
	}

	resource "snowflake_role_grants" "okta" {
		role_name = snowflake_role.okta.name
		roles = ["ACCOUNTADMIN"]
	}

	resource "snowflake_network_policy" "okta" {
		name            = "OKTA_NETWORK_POLICY"
		allowed_ip_list = ["192.168.0.100/24", "29.254.123.20"]
	}

	resource "snowflake_scim_integration" "test_okta_int_np" {
		name = "%s"
		scim_client = "OKTA"
		provisioner_role = snowflake_role.okta.name
		network_policy = snowflake_network_policy.okta.name
		depends_on = [
			snowflake_account_grant.oktacu,
			snowflake_account_grant.oktacr,
			snowflake_role_grants.okta
		]
	}
	`, name)
}
