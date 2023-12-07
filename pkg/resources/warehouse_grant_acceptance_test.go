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

func TestAcc_WarehouseGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_WAREHOUSE_GRANT_TESTS"); ok {
		t.Skip("Skipping TestAcc_WarehouseGrant")
	}
	wName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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
