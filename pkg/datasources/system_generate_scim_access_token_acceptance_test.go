package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SystemGenerateSCIMAccessToken(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	scimId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleId := snowflakeroles.AadProvisioner

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config: generateAccessTokenConfig(scimId, roleId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_system_generate_scim_access_token.p", "integration_name", scimId.Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_system_generate_scim_access_token.p", "access_token"),
				),
			},
		},
	})
}

func generateAccessTokenConfig(scimId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
	resource "snowflake_scim_integration" "azured" {
		name = "%[1]s"
		enabled = true
		scim_client = "AZURE"
		run_as_role = "%[2]s"
	}

	data snowflake_system_generate_scim_access_token p {
		integration_name = snowflake_scim_integration.azured.name
	}
	`, scimId.Name(), roleId.Name())
}
