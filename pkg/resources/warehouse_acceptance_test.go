package resources_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Warehouse(t *testing.T) {
	warehouseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	warehouseId2 := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	prefix := warehouseId.Name()
	prefix2 := warehouseId2.Name()
	comment := random.Comment()
	newComment := random.Comment()

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
				PreConfig: func() { acc.TestClient().Warehouse.UpdateMaxConcurrencyLevel(t, warehouseId2, 10) },
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
	prefix := acc.TestClient().Ids.Alpha()

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

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2763
// TODO [SNOW-1348102]: probably to remove with warehouse rework (we will remove default and also logic with enable_query_acceleration seems superficial - nothing in the docs)
func TestAcc_Warehouse_Issue2763(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouseWithOptions(t, id, &sdk.CreateWarehouseOptions{
						EnableQueryAcceleration: sdk.Bool(false),
					})
					t.Cleanup(warehouseCleanup)
				},
				Config:             wConfigWithQueryAcceleration(id.Name()),
				ResourceName:       "snowflake_warehouse.w",
				ImportState:        true,
				ImportStateId:      id.Name(),
				ImportStatePersist: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					var warehouse *terraform.InstanceState
					if len(s) != 1 {
						return fmt.Errorf("expected 1 state: %#v", s)
					}
					warehouse = s[0]
					// verify that query_acceleration_max_scale_factor is not set in state after import
					_, ok := warehouse.Attributes["query_acceleration_max_scale_factor"]
					if ok {
						return fmt.Errorf("query_acceleration_max_scale_factor is present in state but shouldn't")
					}
					warehouseInSnowflake, err := acc.TestClient().Warehouse.Show(t, id)
					if err != nil {
						return fmt.Errorf("error getting warehouse from SF: %w", err)
					}
					// verify that by default QueryAccelerationMaxScaleFactor is 8 in SF
					if warehouseInSnowflake.QueryAccelerationMaxScaleFactor != 8 {
						return fmt.Errorf("expected QueryAccelerationMaxScaleFactor to be equal to 8 but got %d", warehouseInSnowflake.QueryAccelerationMaxScaleFactor)
					}
					return nil
				},
			},
			{
				Config: wConfigWithQueryAcceleration(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "enable_query_acceleration", "false"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "query_acceleration_max_scale_factor"),
				),
			},
		},
	})
}

func wConfigWithQueryAcceleration(name string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name                      = "%s"
    enable_query_acceleration = false
    query_acceleration_max_scale_factor = 8
}
`, name)
}

func TestAcc_Warehouse_WarehouseSizes(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// set up with concrete size
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						printPlanDetails("snowflake_warehouse.w", "warehouse_size", "warehouse_size_sf"),
					},
				},
				Config: warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeSmall)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeSmall)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size_sf", string(sdk.WarehouseSizeSmall)),
					testAccCheckWarehouseSize(t, id, sdk.WarehouseSizeSmall),
				),
			},
			// import when size in config
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
				ImportStateCheck: ComposeImportStateCheck(
					testCheckResourceAttrInstanceState(id.Name(), "warehouse_size", string(sdk.WarehouseSizeSmall)),
					testCheckResourceAttrInstanceState(id.Name(), "warehouse_size_sf", string(sdk.WarehouseSizeSmall)),
					//testCheckResourceAttrInstanceState(id.Name(), "warehouse_size_sf_changed", "false"),
				),
			},
			// change size in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						printPlanDetails("snowflake_warehouse.w", "warehouse_size", "warehouse_size_sf"),
					},
				},
				Config: warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeMedium)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeMedium)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size_sf", string(sdk.WarehouseSizeMedium)),
					testAccCheckWarehouseSize(t, id, sdk.WarehouseSizeMedium),
				),
			},
			// remove size from config
			{
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						printPlanDetails("snowflake_warehouse.w", "warehouse_size", "warehouse_size_sf"),
						plancheck.ExpectResourceAction("snowflake_warehouse.w", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", ""),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size_sf", string(sdk.WarehouseSizeXSmall)),
					testAccCheckWarehouseSize(t, id, sdk.WarehouseSizeXSmall),
				),
			},
			// add config (lower case)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						printPlanDetails("snowflake_warehouse.w", "warehouse_size", "warehouse_size_sf"),
					},
				},
				Config: warehouseWithSizeConfig(id.Name(), strings.ToLower(string(sdk.WarehouseSizeSmall))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", strings.ToLower(string(sdk.WarehouseSizeSmall))),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size_sf", string(sdk.WarehouseSizeSmall)),
					testAccCheckWarehouseSize(t, id, sdk.WarehouseSizeSmall),
				),
			},
			// remove size from config but update warehouse externally to default (still expecting non-empty plan because we do not know the default)
			{
				PreConfig: func() {
					acc.TestClient().Warehouse.UpdateWarehouseSize(t, id, sdk.WarehouseSizeXSmall)
				},
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						printPlanDetails("snowflake_warehouse.w", "warehouse_size", "warehouse_size_sf"),
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", ""),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size_sf", string(sdk.WarehouseSizeXSmall)),
					testAccCheckWarehouseSize(t, id, sdk.WarehouseSizeXSmall),
				),
			},
			// change the size externally
			{
				PreConfig: func() {
					// we change the size to the size different from default, expecting action
					acc.TestClient().Warehouse.UpdateWarehouseSize(t, id, sdk.WarehouseSizeSmall)
				},
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						printPlanDetails("snowflake_warehouse.w", "warehouse_size", "warehouse_size_sf"),
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", ""),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size_sf", string(sdk.WarehouseSizeXSmall)),
					testAccCheckWarehouseSize(t, id, sdk.WarehouseSizeXSmall),
				),
			},
			// import when no size in config
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				//ImportStateVerify: true,
				//ImportStateVerifyIgnore: []string{
				//	"initially_suspended",
				//	"wait_for_provisioning",
				//	"query_acceleration_max_scale_factor",
				//	"max_concurrency_level",
				//	"statement_queued_timeout_in_seconds",
				//	"statement_timeout_in_seconds",
				//},
				ImportStateCheck: ComposeImportStateCheck(
					testCheckResourceAttrInstanceState(id.Name(), "warehouse_size", string(sdk.WarehouseSizeXSmall)),
					testCheckResourceAttrInstanceState(id.Name(), "warehouse_size_sf", string(sdk.WarehouseSizeXSmall)),
					//testCheckResourceAttrInstanceState(id.Name(), "warehouse_size_sf_changed", "false"),
				),
			},
		},
	})
}

func TestAcc_Warehouse_SizeValidation(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				Config:      warehouseWithSizeConfig(id.Name(), "SMALLa"),
				ExpectError: regexp.MustCompile(`expected a valid warehouse size, got "SMALLa"`),
			},
		},
	})
}

func TestAcc_Warehouse_migrateFromVersion091_withWarehouseSize(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.91.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeX4Large)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", "4XLARGE"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeX4Large)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeX4Large)),
				),
			},
		},
	})
}

func TestAcc_Warehouse_migrateFromVersion091_withoutWarehouseSize(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.91.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: warehouseBasicConfig(id.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeXSmall)),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeXSmall)),
				),
			},
		},
	})
}

func warehouseWithSizeConfig(name string, size string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name           = "%s"
	warehouse_size = "%s"
}
`, name, size)
}

func warehouseBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name           = "%s"
}
`, name)
}

func testAccCheckWarehouseSize(t *testing.T, id sdk.AccountObjectIdentifier, expectedSize sdk.WarehouseSize) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		warehouse, err := acc.TestClient().Warehouse.Show(t, id)
		if err != nil {
			return err
		}
		if warehouse.Size != expectedSize {
			return fmt.Errorf("expected size: %s; got: %s", expectedSize, warehouse.Size)
		}
		return nil
	}
}

func ComposeImportStateCheck(fs ...resource.ImportStateCheckFunc) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		for i, f := range fs {
			if err := f(s); err != nil {
				return fmt.Errorf("check %d/%d error: %s", i+1, len(fs), err)
			}
		}
		return nil
	}
}

func testCheckResourceAttrInstanceState(id string, attributeName, attributeValue string) resource.ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		for _, v := range is {
			if v.ID != id {
				continue
			}

			if attrVal, ok := v.Attributes[attributeName]; ok {
				if attrVal != attributeValue {
					return fmt.Errorf("expected: %s got: %s", attributeValue, attrVal)
				}

				return nil
			}
		}

		return fmt.Errorf("attribute %s not found in instance state", attributeName)
	}
}

var _ plancheck.PlanCheck = printingPlanCheck{}

type printingPlanCheck struct {
	resourceAddress string
	attributes      []string
}

func (e printingPlanCheck) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result []error

	for _, change := range req.Plan.ResourceDrift {
		if e.resourceAddress != change.Address {
			continue
		}
		actions := change.Change.Actions
		var before, after, computed map[string]any
		if change.Change.Before != nil {
			before = change.Change.Before.(map[string]any)
		}
		if change.Change.After != nil {
			after = change.Change.After.(map[string]any)
		}
		if change.Change.AfterUnknown != nil {
			computed = change.Change.AfterUnknown.(map[string]any)
		}
		fmt.Printf("resource drift for [%s]: actions: %v\n", change.Address, actions)
		for _, attr := range e.attributes {
			valueBefore, _ := before[attr]
			valueAfter, _ := after[attr]
			_, isComputed := computed[attr]
			fmt.Printf("\t[%s]: before: %v, after: %v, computed: %t\n", attr, valueBefore, valueAfter, isComputed)
		}
	}

	for _, change := range req.Plan.ResourceChanges {
		if e.resourceAddress != change.Address {
			continue
		}
		actions := change.Change.Actions
		var before, after, computed map[string]any
		if change.Change.Before != nil {
			before = change.Change.Before.(map[string]any)
		}
		if change.Change.After != nil {
			after = change.Change.After.(map[string]any)
		}
		if change.Change.AfterUnknown != nil {
			computed = change.Change.AfterUnknown.(map[string]any)
		}
		fmt.Printf("resource change for [%s]: actions: %v\n", change.Address, actions)
		for _, attr := range e.attributes {
			valueBefore, _ := before[attr]
			valueAfter, _ := after[attr]
			_, isComputed := computed[attr]
			fmt.Printf("\t[%s]: before: %v, after: %v, computed: %t\n", attr, valueBefore, valueAfter, isComputed)
		}
	}

	resp.Error = errors.Join(result...)
}

func printPlanDetails(resourceAddress string, attributes ...string) plancheck.PlanCheck {
	return printingPlanCheck{
		resourceAddress,
		attributes,
	}
}
