package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Warehouses(t *testing.T) {
	warehouseName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: warehouses(warehouseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.s", "warehouses.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.s", "warehouses.0.name"),
				),
			},
		},
	})
}

func warehouses(warehouseName string) string {
	return fmt.Sprintf(`
	resource snowflake_warehouse "s"{
		name                         = "%v"
		warehouse_size               = "XSMALL"
		initially_suspended          = true
		auto_suspend                 = 60
		max_concurrency_level        = 8
		statement_timeout_in_seconds = 172800
		warehouse_type               = "STANDARD"
	}

	data snowflake_warehouses "s" {
		depends_on = [snowflake_warehouse.s]
	}
	`, warehouseName)
}
