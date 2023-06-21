package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Warehouse(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_WAREHOUSE_TESTS"); ok {
		t.Skip("Skipping TestAccWarehouse")
	}

	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: wConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", "test comment"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "60"),
					resource.TestCheckResourceAttrSet("snowflake_warehouse.w", "warehouse_size"),
				),
			},
			// RENAME
			{
				Config: wConfig(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", "test comment"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "60"),
					resource.TestCheckResourceAttrSet("snowflake_warehouse.w", "warehouse_size"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: wConfig2(prefix2, "X-LARGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", "test comment 2"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "60"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", "X-LARGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_warehouse.w",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"initially_suspended",
					"wait_for_provisioning",
					"query_acceleration_max_scale_factor",
					"max_concurrency_level",
					"statement_queued_timeout_in_seconds",
					"statement_timeout_in_seconds",
				},
			},
		},
	})
}

func TestAcc_WarehousePattern(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_WAREHOUSE_TESTS"); ok {
		t.Skip("Skipping TestAccWarehouse")
	}

	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: wConfigPattern(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w1", "name", fmt.Sprintf("%s_", prefix)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w2", "name", fmt.Sprintf("%s1", prefix)),
				),
			},
		},
	})
}

func wConfig(prefix string) string {
	s := `
resource "snowflake_warehouse" "w" {
	name    = "%s"
	comment = "test comment"

	auto_suspend          = 60
	max_cluster_count     = 1
	min_cluster_count     = 1
	scaling_policy        = "STANDARD"
	auto_resume           = true
	initially_suspended   = true
	wait_for_provisioning = false
}
`
	return fmt.Sprintf(s, prefix)
}

func wConfig2(prefix string, size string) string {
	s := `
resource "snowflake_warehouse" "w" {
	name           = "%s"
	comment        = "test comment 2"
	warehouse_size = "%s"

	auto_suspend          = 60
	max_cluster_count     = 1
	min_cluster_count     = 1
	scaling_policy        = "STANDARD"
	auto_resume           = true
	initially_suspended   = true
	wait_for_provisioning = false
}
`
	return fmt.Sprintf(s, prefix, size)
}

func wConfigPattern(prefix string) string {
	s := `
resource "snowflake_warehouse" "w1" {
	name           = "%s_"
}
resource "snowflake_warehouse" "w2" {
	name           = "%s1"
}
`
	return fmt.Sprintf(s, prefix, prefix)
}
