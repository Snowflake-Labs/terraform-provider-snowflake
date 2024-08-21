package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SystemGenerateSCIMAccessToken(t *testing.T) {
	scimIntName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: generateAccessTokenConfig(scimIntName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_system_generate_scim_access_token.p", "integration_name", scimIntName),
					resource.TestCheckResourceAttrSet("data.snowflake_system_generate_scim_access_token.p", "access_token"),
				),
			},
		},
	})
}

func generateAccessTokenConfig(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_account_role" "azured" {
		name = "AAD_PROVISIONER"
		comment = "test comment"
	}

	resource "snowflake_grant_privileges_to_account_role" "azure_grants" {
	  	account_role_name = snowflake_account_role.azured.name
  		privileges        = ["CREATE USER", "CREATE ROLE"]
		on_account        = true
	}

	resource "snowflake_grant_account_role" "azured" {
		role_name        = snowflake_account_role.azured.name
		parent_role_name = "ACCOUNTADMIN"
	}

	resource "snowflake_scim_integration" "azured" {
		name = "%s"
		enabled = true
		scim_client = "AZURE"
		run_as_role = snowflake_account_role.azured.name
		depends_on = [
			snowflake_grant_privileges_to_account_role.azure_grants,
			snowflake_grant_account_role.azured
		]
	}

	data snowflake_system_generate_scim_access_token p {
		integration_name = snowflake_scim_integration.azured.name
		depends_on = [snowflake_scim_integration.azured]
	}
	`, name)
}
