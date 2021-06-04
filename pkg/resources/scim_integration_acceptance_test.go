package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ScimIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SCIM_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAccScimIntegration")
	}

	scimIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	scimProvisionerRole := "AAD_PROVISIONER"
	scimNetworkPolicy := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: scimIntegrationConfig_azure(scimIntName, scimProvisionerRole, scimNetworkPolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", scimIntName),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "AZURE"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "provisioner_role", scimProvisionerRole),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", scimNetworkPolicy),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test", "created_on"),
				),
			},
			{
				ResourceName:      "snowflake_scim_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func scimIntegrationConfig_azure(name string, role string, policy string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "azure" {
		name = "%s"
		comment = "test comment"
	}

	resource "snowflake_account_grant" "azurecua" {
		roles     = [snowflake_role.azure.name]
		privilege = "CREATE USER"
	}

	resource "snowflake_account_grant" "azurecra" {
		roles     = [snowflake_role.azure.name]
		privilege = "CREATE ROLE"
	}

	resource "snowflake_role_grants" "azure" {
		role_name = snowflake_role.azure.name
		roles = ["ACCOUNTADMIN"]
	}

	resource "snowflake_network_policy" "azure" {
		name            = "%s"
		allowed_ip_list = ["192.168.0.100/24", "29.254.123.20"]
	}

	resource "snowflake_scim_integration" "test" {
		name = "%s"
		scim_client = "AZURE"
		provisioner_role = snowflake_role.azure.name
		network_policy = snowflake_network_policy.azure.name
		depends_on = [
			snowflake_account_grant.azurecua,
			snowflake_account_grant.azurecra,
			snowflake_role_grants.azure
		]
	}
	`, role, policy, name)
}
