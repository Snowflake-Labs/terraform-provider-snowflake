package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SystemGenerateSCIMAccessToken(t *testing.T) {
	scimIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
	resource "snowflake_role" "azured" {
		name = "AAD_PROVISIONER"
		comment = "test comment"
	}

	resource "snowflake_account_grant" "azurecud" {
		roles     = [snowflake_role.azured.name]
		privilege = "CREATE USER"
	}
	resource "snowflake_account_grant" "azurecrd" {
		roles     = [snowflake_role.azured.name]
		privilege = "CREATE ROLE"
	}
	resource "snowflake_role_grants" "azured" {
		role_name = snowflake_role.azured.name
		roles = ["ACCOUNTADMIN"]
	}

	resource "snowflake_scim_integration" "azured" {
		name = "%s"
		scim_client = "AZURE"
		provisioner_role = snowflake_role.azured.name
		depends_on = [
			snowflake_account_grant.azurecud,
			snowflake_account_grant.azurecrd,
			snowflake_role_grants.azured
		]
	}

	data snowflake_system_generate_scim_access_token p {
		integration_name = snowflake_scim_integration.azured.name
		depends_on = [snowflake_scim_integration.azured]
	}
	`, name)
}
