package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Warehouse(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	newComment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				Config: wConfig(prefix, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "60"),
					resource.TestCheckResourceAttrSet("snowflake_warehouse.w", "warehouse_size"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "min_cluster_count", "1"),
				),
			},
			// RENAME
			{
				Config: wConfig(prefix2, newComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_warehouse.w", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "60"),
					resource.TestCheckResourceAttrSet("snowflake_warehouse.w", "warehouse_size"),
				),
			},
			// CHANGE PROPERTIES (proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2652)
			{
				Config: wConfig2(prefix2, "X-LARGE", 20, 2, newComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "60"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", "XLARGE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_concurrency_level", "20"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "min_cluster_count", "2"),
				),
			},
			// CHANGE JUST max_concurrency_level
			{
				Config: wConfig2(prefix2, "XLARGE", 16, 2, newComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "60"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", "XLARGE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_concurrency_level", "16"),
				),
			},
			// CHANGE max_concurrency_level EXTERNALLY (proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2318)
			{
				PreConfig: func() { alterWarehouseMaxConcurrencyLevelExternally(t, prefix2, 10) },
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Config: wConfig2(prefix2, "XLARGE", 16, 2, newComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "60"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", "XLARGE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_concurrency_level", "16"),
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
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
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

func wConfig(prefix string, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name    = "%s"
	comment = "%s"

	auto_suspend          = 60
	max_cluster_count     = 4
	min_cluster_count     = 1
	scaling_policy        = "STANDARD"
	auto_resume           = true
	initially_suspended   = true
	wait_for_provisioning = false
}
`, prefix, comment)
}

func wConfig2(prefix string, size string, maxConcurrencyLevel int, minClusterCount int, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name           = "%[1]s"
	comment        = "%[5]s"
	warehouse_size = "%[2]s"

	auto_suspend          = 60
	max_cluster_count     = 4
	min_cluster_count     = %[4]d
	scaling_policy        = "STANDARD"
	auto_resume           = true
	initially_suspended   = true
	wait_for_provisioning = false
	max_concurrency_level = %[3]d
}
`, prefix, size, maxConcurrencyLevel, minClusterCount, comment)
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

func alterWarehouseMaxConcurrencyLevelExternally(t *testing.T, warehouseId string, level int) {
	t.Helper()

	client := acc.Client(t)
	ctx := context.Background()

	err := client.Warehouses.Alter(ctx, sdk.NewAccountObjectIdentifier(warehouseId), &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{MaxConcurrencyLevel: sdk.Int(level)}})
	require.NoError(t, err)
}
