package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_WarehouseGrant(t *testing.T) {
	wName := acc.TestClient().Ids.Alpha()
	roleName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: warehouseGrantConfig(wName, roleName, "USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse_grant.test", "warehouse_name", wName),
					resource.TestCheckResourceAttr("snowflake_warehouse_grant.test", "privilege", "USAGE"),
				),
			},
			// UPDATE
			{
				Config: warehouseGrantConfig(wName, roleName, "ALL PRIVILEGES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse_grant.test", "warehouse_name", wName),
					resource.TestCheckResourceAttr("snowflake_warehouse_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_warehouse_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func warehouseGrantConfig(n, role, privilege string) string {
	return fmt.Sprintf(`

resource "snowflake_warehouse" "test" {
  name      = "%v"
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_warehouse_grant" "test" {
  warehouse_name = snowflake_warehouse.test.name
  roles          = [snowflake_role.test.name]
  privilege      = "%s"
}
`, n, role, privilege)
}
