package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ScimIntegration(t *testing.T) {
	scimIntName := acc.TestClient().Ids.Alpha()
	scimProvisionerRole := "AAD_PROVISIONER"
	scimNetworkPolicy := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: scimIntegrationConfigAzure(scimIntName, scimProvisionerRole, scimNetworkPolicy),
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

func scimIntegrationConfigAzure(name string, role string, policy string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "azure" {
		name = "%s"
		comment = "test comment"
	}

	resource "snowflake_grant_privileges_to_account_role" "azure_grants" {
	  	account_role_name = snowflake_role.azure.name
  		privileges        = ["CREATE USER", "CREATE ROLE"]
		on_account        = true
	}

	resource "snowflake_grant_account_role" "azure" {
		role_name        = snowflake_role.azure.name
		parent_role_name = "ACCOUNTADMIN"
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
			snowflake_grant_privileges_to_account_role.azure_grants,
			snowflake_grant_account_role.azure
		]
	}
	`, role, policy, name)
}
