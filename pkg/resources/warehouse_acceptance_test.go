package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/snowflakechecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

// TODO: test import for empty config
func TestAcc_Warehouse(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	warehouseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	warehouseId2 := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	name := warehouseId.Name()
	name2 := warehouseId2.Name()
	comment := random.Comment()
	newComment := random.Comment()

	resourceMonitor, resourceMonitorCleanup := acc.TestClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)
	resourceMonitorId := resourceMonitor.ID()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				Config: warehouseBasicConfigWithComment(name, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", name),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "warehouse_type"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "warehouse_size"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "max_cluster_count"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "min_cluster_count"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "scaling_policy"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "-1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_resume", "unknown"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "initially_suspended"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "resource_monitor"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "enable_query_acceleration", "unknown"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "query_acceleration_max_scale_factor", "-1"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_concurrency_level", "-1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", "-1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "-1"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.type", string(sdk.WarehouseTypeStandard)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.size", string(sdk.WarehouseSizeXSmall)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.max_cluster_count", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.min_cluster_count", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.scaling_policy", string(sdk.ScalingPolicyStandard)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_suspend", "600"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_resume", "true"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.resource_monitor", "null"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.comment", comment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.enable_query_acceleration", "false"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.query_acceleration_max_scale_factor", "8"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.max_concurrency_level.0.value", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
				),
			},
			// RENAME
			{
				Config: warehouseBasicConfigWithComment(name2, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_warehouse.w", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", name2),
				),
			},
			// Change config but use defaults for every attribute (expecting changes in the current approach) - ideally it should result in no changes
			{
				Config: warehouseFullDefaultConfig(name2, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_type", tfjson.ActionUpdate, nil, sdk.String(string(sdk.WarehouseTypeStandard))),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_size", tfjson.ActionUpdate, nil, sdk.String(string(sdk.WarehouseSizeXSmall))),
						planchecks.ExpectChange("snowflake_warehouse.w", "max_cluster_count", tfjson.ActionUpdate, nil, sdk.String("1")),
						planchecks.ExpectChange("snowflake_warehouse.w", "min_cluster_count", tfjson.ActionUpdate, nil, sdk.String("1")),
						planchecks.ExpectChange("snowflake_warehouse.w", "scaling_policy", tfjson.ActionUpdate, nil, sdk.String(string(sdk.ScalingPolicyStandard))),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_suspend", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("600")),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_resume", tfjson.ActionUpdate, sdk.String("unknown"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_warehouse.w", "enable_query_acceleration", tfjson.ActionUpdate, sdk.String("unknown"), sdk.String("false")),
						planchecks.ExpectChange("snowflake_warehouse.w", "query_acceleration_max_scale_factor", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("8")),

						planchecks.ExpectChange("snowflake_warehouse.w", "max_concurrency_level", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("8")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("0")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("172800")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", string(sdk.WarehouseTypeStandard)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeXSmall)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_cluster_count", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "min_cluster_count", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "scaling_policy", string(sdk.ScalingPolicyStandard)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "600"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_resume", "true"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "initially_suspended"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "resource_monitor"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "enable_query_acceleration", "false"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "query_acceleration_max_scale_factor", "8"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "172800"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: warehouseFullConfigNoDefaults(name2, newComment, resourceMonitorId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", string(sdk.WarehouseTypeSnowparkOptimized)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeMedium)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_cluster_count", "4"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "min_cluster_count", "2"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "scaling_policy", string(sdk.ScalingPolicyEconomy)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "1200"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_resume", "false"),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "initially_suspended"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "resource_monitor", resourceMonitorId.Name()),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "enable_query_acceleration", "true"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "query_acceleration_max_scale_factor", "4"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_concurrency_level", "4"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", "5"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "86400"),
				),
			},
			// CHANGE max_concurrency_level EXTERNALLY (proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2318)
			{
				PreConfig: func() { acc.TestClient().Warehouse.UpdateMaxConcurrencyLevel(t, warehouseId2, 10) },
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift("snowflake_warehouse.w", "max_concurrency_level", sdk.String("4"), sdk.String("10")),
						planchecks.ExpectChange("snowflake_warehouse.w", "max_concurrency_level", tfjson.ActionUpdate, sdk.String("10"), sdk.String("4")),
					},
				},
				Config: warehouseFullConfigNoDefaults(name2, newComment, resourceMonitorId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", name2),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "max_concurrency_level", "4"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.max_concurrency_level.0.value", "4"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_warehouse.w",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resource_monitor",
				},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(warehouseId2.Name(), "resource_monitor", resourceMonitorId.Name()),
				),
			},
		},
	})
}

func TestAcc_Warehouse_WarehouseType(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// set up with concrete type
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_type", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_type", tfjson.ActionCreate, nil, sdk.String(string(sdk.WarehouseTypeStandard))),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Config: warehouseWithTypeConfig(id.Name(), sdk.WarehouseTypeStandard, sdk.WarehouseSizeMedium),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", string(sdk.WarehouseTypeStandard)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.type", string(sdk.WarehouseTypeStandard)),
					snowflakechecks.CheckWarehouseType(t, id, sdk.WarehouseTypeStandard),
				),
			},
			// import when type in config
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "warehouse_type", string(sdk.WarehouseTypeStandard)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.0.type", string(sdk.WarehouseTypeStandard)),
				),
			},
			// change type in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_type", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Config: warehouseWithTypeConfig(id.Name(), sdk.WarehouseTypeSnowparkOptimized, sdk.WarehouseSizeMedium),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", string(sdk.WarehouseTypeSnowparkOptimized)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.type", string(sdk.WarehouseTypeSnowparkOptimized)),
					snowflakechecks.CheckWarehouseType(t, id, sdk.WarehouseTypeSnowparkOptimized),
				),
			},
			// remove type from config
			{
				Config: warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeMedium)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_warehouse.w", plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_type", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized)), nil),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", ""),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.type", string(sdk.WarehouseTypeStandard)),
					snowflakechecks.CheckWarehouseType(t, id, sdk.WarehouseTypeStandard),
				),
			},
			// add config (lower case)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_type", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_type", tfjson.ActionUpdate, nil, sdk.String(strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized)))),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Config: warehouseWithTypeConfig(id.Name(), sdk.WarehouseType(strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized))), sdk.WarehouseSizeMedium),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized))),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.type", string(sdk.WarehouseTypeSnowparkOptimized)),
					snowflakechecks.CheckWarehouseType(t, id, sdk.WarehouseTypeSnowparkOptimized),
				),
			},
			// remove type from config but update warehouse externally to default (still expecting non-empty plan because we do not know the default)
			{
				PreConfig: func() {
					acc.TestClient().Warehouse.UpdateWarehouseType(t, id, sdk.WarehouseTypeStandard)
				},
				Config: warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeMedium)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_type", "show_output"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "warehouse_type", sdk.String(strings.ToLower(string(sdk.WarehouseTypeSnowparkOptimized))), sdk.String(string(sdk.WarehouseTypeStandard))),
						planchecks.ExpectDrift("snowflake_warehouse.w", "show_output.0.type", sdk.String(string(sdk.WarehouseTypeSnowparkOptimized)), sdk.String(string(sdk.WarehouseTypeStandard))),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeStandard)), nil),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", ""),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.type", string(sdk.WarehouseTypeStandard)),
					snowflakechecks.CheckWarehouseType(t, id, sdk.WarehouseTypeStandard),
				),
			},
			// change the size externally
			{
				PreConfig: func() {
					// we change the type to the type different from default, expecting action
					acc.TestClient().Warehouse.UpdateWarehouseType(t, id, sdk.WarehouseTypeSnowparkOptimized)
				},
				Config: warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeMedium)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_type", "show_output"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "warehouse_type", nil, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectDrift("snowflake_warehouse.w", "show_output.0.type", sdk.String(string(sdk.WarehouseTypeStandard)), sdk.String(string(sdk.WarehouseTypeSnowparkOptimized))),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_type", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseTypeSnowparkOptimized)), nil),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", ""),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.type", string(sdk.WarehouseTypeStandard)),
					snowflakechecks.CheckWarehouseType(t, id, sdk.WarehouseTypeStandard),
				),
			},
			// import when no type in config
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "warehouse_type", string(sdk.WarehouseTypeStandard)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.0.type", string(sdk.WarehouseTypeStandard)),
				),
			},
		},
	})
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
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_size", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_size", tfjson.ActionCreate, nil, sdk.String(string(sdk.WarehouseSizeSmall))),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Config: warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeSmall)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeSmall)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.size", string(sdk.WarehouseSizeSmall)),
					snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeSmall),
				),
			},
			// import when size in config
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "warehouse_size", string(sdk.WarehouseSizeSmall)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.0.size", string(sdk.WarehouseSizeSmall)),
				),
			},
			// change size in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_size", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_size", tfjson.ActionUpdate, sdk.String(string(sdk.WarehouseSizeSmall)), sdk.String(string(sdk.WarehouseSizeMedium))),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Config: warehouseWithSizeConfig(id.Name(), string(sdk.WarehouseSizeMedium)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeMedium)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.size", string(sdk.WarehouseSizeMedium)),
					snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeMedium),
				),
			},
			// remove size from config
			{
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_warehouse.w", plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_size", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeMedium)), nil),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "warehouse_size"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.size", string(sdk.WarehouseSizeXSmall)),
					snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeXSmall),
				),
			},
			// add config (lower case)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_size", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_size", tfjson.ActionUpdate, nil, sdk.String(strings.ToLower(string(sdk.WarehouseSizeSmall)))),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Config: warehouseWithSizeConfig(id.Name(), strings.ToLower(string(sdk.WarehouseSizeSmall))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", strings.ToLower(string(sdk.WarehouseSizeSmall))),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.size", string(sdk.WarehouseSizeSmall)),
					snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeSmall),
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
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_size", "show_output"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "warehouse_size", sdk.String(strings.ToLower(string(sdk.WarehouseSizeSmall))), sdk.String(string(sdk.WarehouseSizeXSmall))),
						planchecks.ExpectDrift("snowflake_warehouse.w", "show_output.0.size", sdk.String(string(sdk.WarehouseSizeSmall)), sdk.String(string(sdk.WarehouseSizeXSmall))),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeXSmall)), nil),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "warehouse_size"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.size", string(sdk.WarehouseSizeXSmall)),
					snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeXSmall),
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
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "warehouse_size", "show_output"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "warehouse_size", nil, sdk.String(string(sdk.WarehouseSizeSmall))),
						planchecks.ExpectDrift("snowflake_warehouse.w", "show_output.0.size", sdk.String(string(sdk.WarehouseSizeXSmall)), sdk.String(string(sdk.WarehouseSizeSmall))),
						planchecks.ExpectChange("snowflake_warehouse.w", "warehouse_size", tfjson.ActionCreate, sdk.String(string(sdk.WarehouseSizeSmall)), nil),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "warehouse_size"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.size", string(sdk.WarehouseSizeXSmall)),
					snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeXSmall),
				),
			},
			// import when no size in config
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "warehouse_size", string(sdk.WarehouseSizeXSmall)),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.0.size", string(sdk.WarehouseSizeXSmall)),
				),
			},
		},
	})
}

// [SNOW-1348102 - next PR]: add more validations
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
				ExpectError: regexp.MustCompile("invalid warehouse size: SMALLa"),
			},
		},
	})
}

// TestAcc_Warehouse_AutoResume validates behavior for falling back to Snowflake default for boolean attribute
func TestAcc_Warehouse_AutoResume(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// set up with auto resume set in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "auto_resume", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_resume", tfjson.ActionCreate, nil, sdk.String("true")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Config: warehouseWithAutoResumeConfig(id.Name(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_resume", "true"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_resume", "true"),
					snowflakechecks.CheckAutoResume(t, id, true),
				),
			},
			// import when type in config
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "auto_resume", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.0.auto_resume", "true"),
				),
			},
			// change value in config
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "auto_resume", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_resume", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Config: warehouseWithAutoResumeConfig(id.Name(), false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_resume", "false"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_resume", "false"),
					snowflakechecks.CheckAutoResume(t, id, false),
				),
			},
			// remove type from config (expecting non-empty plan because we do not know the default)
			{
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_warehouse.w", plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "auto_resume", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_resume", tfjson.ActionUpdate, sdk.String("false"), sdk.String("unknown")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_resume", "unknown"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_resume", "false"),
					snowflakechecks.CheckAutoResume(t, id, false),
				),
			},
			// change auto resume externally
			{
				PreConfig: func() {
					// we change the auto resume to the type different from default, expecting action
					acc.TestClient().Warehouse.UpdateAutoResume(t, id, true)
				},
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "auto_resume", "show_output"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "auto_resume", sdk.String("unknown"), sdk.String("true")),
						planchecks.ExpectDrift("snowflake_warehouse.w", "show_output.0.auto_resume", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_resume", tfjson.ActionUpdate, sdk.String("true"), sdk.String("unknown")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_resume", "unknown"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_resume", "false"),
					snowflakechecks.CheckWarehouseType(t, id, sdk.WarehouseTypeStandard),
				),
			},
			// import when no type in config
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "auto_resume", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "show_output.0.auto_resume", "false"),
				),
			},
		},
	})
}

func TestAcc_Warehouse_ZeroValues(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// create with valid "zero" values
			{
				Config: warehouseWithAllValidZeroValuesConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_suspend", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange("snowflake_warehouse.w", "query_acceleration_max_scale_factor", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("0")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "query_acceleration_max_scale_factor", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "0"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_suspend", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.query_acceleration_max_scale_factor", "0"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_queued_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),

					// TODO [SNOW-1348102 - next PR]: snowflake checks?
					// snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeSmall),
				),
			},
			// remove all from config (to validate that unset is run correctly)
			{
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_suspend", tfjson.ActionUpdate, sdk.String("0"), sdk.String("-1")),
						planchecks.ExpectChange("snowflake_warehouse.w", "query_acceleration_max_scale_factor", tfjson.ActionUpdate, sdk.String("0"), sdk.String("-1")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("0"), sdk.String("-1")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("0"), sdk.String("-1")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "-1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "query_acceleration_max_scale_factor", "-1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", "-1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "-1"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_suspend", "600"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.query_acceleration_max_scale_factor", "8"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_queued_timeout_in_seconds.0.level", ""),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
			// add valid "zero" values again (to validate if set is run correctly)
			{
				Config: warehouseWithAllValidZeroValuesConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "auto_suspend", "query_acceleration_max_scale_factor", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds", "show_output"),
						planchecks.ExpectChange("snowflake_warehouse.w", "auto_suspend", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("0")),
						planchecks.ExpectChange("snowflake_warehouse.w", "query_acceleration_max_scale_factor", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("0")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("0")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("0")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "show_output", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "auto_suspend", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "query_acceleration_max_scale_factor", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_queued_timeout_in_seconds", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "0"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.auto_suspend", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.query_acceleration_max_scale_factor", "0"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_queued_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
		},
	})
}

func TestAcc_Warehouse_Parameter(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			// create with setting one param
			{
				Config: warehouseWithParameterConfig(id.Name(), 86400),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionCreate, nil, sdk.String("86400")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "parameters", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "86400"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),

					// TODO [SNOW-1348102 - next PR]: snowflake checks?
					// snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeSmall),
				),
			},
			// do not make any change (to check if there is no drift)
			{
				Config: warehouseWithParameterConfig(id.Name(), 86400),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// import when param in config
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "statement_timeout_in_seconds", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change the param value in config
			{
				Config: warehouseWithParameterConfig(id.Name(), 43200),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), sdk.String("43200")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "parameters", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "43200"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "43200"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change param value on account - expect no changes
			{
				PreConfig: func() {
					param := acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					require.Equal(t, "", string(param.Level))
					revert := acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "86400")
					t.Cleanup(revert)
				},
				Config: warehouseWithParameterConfig(id.Name(), 43200),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectChange("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", tfjson.ActionNoop, sdk.String("43200"), sdk.String("43200")),
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "43200"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "43200"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change the param value externally
			{
				PreConfig: func() {
					// clean after previous step
					acc.TestClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					// update externally
					acc.TestClient().Warehouse.UpdateStatementTimeoutInSeconds(t, id, 86400)
				},
				Config: warehouseWithParameterConfig(id.Name(), 43200),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "statement_timeout_in_seconds", sdk.String("43200"), sdk.String("86400")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), sdk.String("43200")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "parameters", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "43200"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "43200"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// change the param value on account to the value from config (but on different level)
			{
				PreConfig: func() {
					acc.TestClient().Warehouse.UnsetStatementTimeoutInSeconds(t, id)
					acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "43200")
				},
				Config: warehouseWithParameterConfig(id.Name(), 43200),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "statement_timeout_in_seconds", sdk.String("43200"), sdk.String("-1")),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("43200")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "parameters", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "43200"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "43200"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// remove the param from config
			{
				PreConfig: func() {
					// clean after previous step
					acc.TestClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					param := acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					require.Equal(t, "", string(param.Level))
				},
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("43200"), sdk.String("-1")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "parameters", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "-1"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
			// import when param not in config (snowflake default)
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "statement_timeout_in_seconds", "-1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
			// change the param value in config to snowflake default
			{
				Config: warehouseWithParameterConfig(id.Name(), 172800),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("-1"), sdk.String("172800")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "parameters", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "172800"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeWarehouse)),
				),
			},
			// remove the param from config
			{
				PreConfig: func() {
					param := acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					require.Equal(t, "", string(param.Level))
				},
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("172800"), sdk.String("-1")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "parameters", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "-1"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
			// change param value on account - change expected to be noop
			{
				PreConfig: func() {
					param := acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
					require.Equal(t, "", string(param.Level))
					revert := acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "86400")
					t.Cleanup(revert)
				},
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", sdk.String("172800"), sdk.String("86400")),
						planchecks.ExpectChange("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", tfjson.ActionNoop, sdk.String("86400"), sdk.String("86400")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "-1"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeAccount)),
				),
			},
			// import when param not in config (set on account)
			{
				ResourceName: "snowflake_warehouse.w",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "statement_timeout_in_seconds", "-1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeAccount)),
				),
			},
			// change param value on warehouse
			{
				PreConfig: func() {
					acc.TestClient().Warehouse.UpdateStatementTimeoutInSeconds(t, id, 86400)
				},
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectChange("snowflake_warehouse.w", "statement_timeout_in_seconds", tfjson.ActionUpdate, sdk.String("86400"), sdk.String("-1")),
						planchecks.ExpectComputed("snowflake_warehouse.w", "parameters", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "-1"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "86400"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", string(sdk.ParameterTypeAccount)),
				),
			},
			// unset param on account
			{
				PreConfig: func() {
					acc.TestClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterStatementTimeoutInSeconds)
				},
				Config: warehouseBasicConfig(id.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "statement_timeout_in_seconds", "parameters"),
						planchecks.ExpectDrift("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", sdk.String("86400"), sdk.String("172800")),
						planchecks.ExpectDrift("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", sdk.String(string(sdk.ParameterTypeAccount)), sdk.String("")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "statement_timeout_in_seconds", "-1"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.value", "172800"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "parameters.0.statement_timeout_in_seconds.0.level", ""),
				),
			},
		},
	})
}

func TestAcc_Warehouse_InitiallySuspendedChangesPostCreation(t *testing.T) {
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
				Config: warehouseWithInitiallySuspendedConfig(id.Name(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "initially_suspended", "true"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.state", string(sdk.WarehouseStateSuspended)),

					// TODO [SNOW-1348102 - next PR]: snowflake checks?
					// snowflakechecks.CheckWarehouseSize(t, id, sdk.WarehouseSizeSmall),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: warehouseWithInitiallySuspendedConfig(id.Name(), false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "initially_suspended", "true"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.state", string(sdk.WarehouseStateSuspended)),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: warehouseBasicConfig(id.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "initially_suspended", "true"),

					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "show_output.0.state", string(sdk.WarehouseStateSuspended)),
				),
			},
		},
	})
}

// TODO: test with all other values set in config - even wait_for_provisioning deprecated field (to validate no change)
// TODO: test with no entries in config in the previous version and full config after
func TestAcc_Warehouse_migrateFromVersion092_withWarehouseSize(t *testing.T) {
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
						VersionConstraint: "=0.92.0",
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
					PreApply: []plancheck.PlanCheck{
						// add checks to all
						plancheck.ExpectResourceAction("snowflake_warehouse.w", plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails("snowflake_warehouse.w", "comment", "enable_query_acceleration", "query_acceleration_max_scale_factor", "warehouse_type", "max_concurrency_level", "statement_queued_timeout_in_seconds", "statement_timeout_in_seconds"),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", string(sdk.WarehouseSizeX4Large)),
					resource.TestCheckNoResourceAttr("snowflake_warehouse.w", "wait_for_provisioning"),
				),
			},
		},
	})
}

// TODO [SNOW-1348102 - next PR]: test defaults removal
// TODO [SNOW-1348102 - next PR]: test basic creation (check previous defaults)
// TODO [SNOW-1348102 - next PR]: test auto_suspend set to 0 before migration
// TODO [SNOW-1348102 - next PR]: do we care about drift in warehouse for is_current warehouse? (test)
// TODO [SNOW-1348102 - next PR]: test boolean type change (with leaving boolean/int in config) and add migration
// TODO [SNOW-1348102 - next PR]: test int, string, identifier changed externally
// TODO [SNOW-1348102 - next PR]: unskip - it fails currently because of other state upograders missing
func TestAcc_Warehouse_migrateFromVersion091_withoutWarehouseSize(t *testing.T) {
	t.Skip("Skipped due to the missing state migrators for other props")
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
						VersionConstraint: "=0.92.0",
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

func warehouseBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name           = "%s"
}
`, name)
}

func warehouseBasicConfigWithComment(name string, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name           = "%s"
	comment        = "%s"
}
`, name, comment)
}

func warehouseFullDefaultConfig(name string, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name                                = "%[1]s"
	warehouse_type                      = "STANDARD"
	warehouse_size                      = "XSMALL"
	max_cluster_count                   = 1
	min_cluster_count                   = 1
	scaling_policy                      = "STANDARD"
	auto_suspend                        = 600
	auto_resume                         = true
	initially_suspended                 = false
	comment                             = "%[2]s"
    enable_query_acceleration           = false
    query_acceleration_max_scale_factor = 8

    max_concurrency_level               = 8
    statement_queued_timeout_in_seconds = 0
    statement_timeout_in_seconds        = 172800
}
`, name, comment)
}

func warehouseFullConfigNoDefaults(name string, comment string, id sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name                                = "%[1]s"
	warehouse_type                      = "SNOWPARK-OPTIMIZED"
	warehouse_size                      = "MEDIUM"
	max_cluster_count                   = 4
	min_cluster_count                   = 2
	scaling_policy                      = "ECONOMY"
	auto_suspend                        = 1200
	auto_resume                         = false
	initially_suspended                 = false
	resource_monitor                    = "%[3]s"
	comment                             = "%[2]s"
    enable_query_acceleration           = true
    query_acceleration_max_scale_factor = 4

    max_concurrency_level               = 4
    statement_queued_timeout_in_seconds = 5
    statement_timeout_in_seconds        = 86400
}
`, name, comment, id.Name())
}

func warehouseWithSizeConfig(name string, size string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name           = "%s"
	warehouse_size = "%s"
}
`, name, size)
}

func warehouseWithTypeConfig(name string, warehouseType sdk.WarehouseType, size sdk.WarehouseSize) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name           = "%s"
	warehouse_type = "%s"
	warehouse_size = "%s"
}
`, name, warehouseType, size)
}

func warehouseWithAutoResumeConfig(name string, autoResume bool) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name        = "%s"
	auto_resume = "%t"
}
`, name, autoResume)
}

func warehouseWithAllValidZeroValuesConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name                                = "%s"
	auto_suspend                        = 0
    query_acceleration_max_scale_factor = 0
    statement_queued_timeout_in_seconds = 0
    statement_timeout_in_seconds        = 0
}
`, name)
}

func warehouseWithParameterConfig(name string, value int) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name                                = "%s"
    statement_timeout_in_seconds        = %d
}
`, name, value)
}

func warehouseWithInitiallySuspendedConfig(name string, initiallySuspened bool) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "w" {
	name                = "%s"
	initially_suspended = %t
}
`, name, initiallySuspened)
}
