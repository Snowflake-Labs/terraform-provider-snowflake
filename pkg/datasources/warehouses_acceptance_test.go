package datasources_test

import (
// "fmt"
// "strings"
// "testing"
//
// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//func TestAcc_Warehouses(t *testing.T) {
//	warehouseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
//	resource.ParallelTest(t, resource.TestCase{
//		Providers:    providers(),
//		CheckDestroy: nil,
//		Steps: []resource.TestStep{
//			{
//				Config: warehouses(warehouseName),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.s", "warehouses.#"),
//					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.s", "warehouses.0.name"),
//				),
//			},
//		},
//	})
//}

// func warehouses(warehouseName string) string {
// 	return fmt.Sprintf(`
// 	resource snowflake_warehouse "s"{
// 		name                         = "%v"
// 		warehouse_size               = "XSMALL"
// 		initially_suspended          = true
// 		auto_suspend                 = 60
// 		max_concurrency_level        = 8
// 		statement_timeout_in_seconds = 172800
// 	}

// 	data snowflake_warehouses "s" {
// 		depends_on = [snowflake_warehouse.s]
// 	}
// 	`, warehouseName)
// }
