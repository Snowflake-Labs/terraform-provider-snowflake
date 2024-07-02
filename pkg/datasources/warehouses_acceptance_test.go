package datasources_test

import (
	"maps"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Warehouses_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.ConfigureClientOnce)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	configVariables := config.Variables{
		"name":    config.StringVariable(id.Name()),
		"comment": config.StringVariable(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Warehouses/optionals_set"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.state"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.type", string(sdk.WarehouseTypeStandard)),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.size", string(sdk.WarehouseSizeXSmall)),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.min_cluster_count", "1"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.max_cluster_count", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.started_clusters"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.running"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.queued"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.is_current", "true"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.auto_suspend", "600"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.auto_resume", "true"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.available"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.provisioning"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.quiescing"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.other"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.created_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.resumed_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.updated_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.enable_query_acceleration", "false"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.query_acceleration_max_scale_factor", "8"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.resource_monitor", ""),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.scaling_policy", string(sdk.ScalingPolicyStandard)),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.owner_role_type"),

					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.describe_output.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.describe_output.0.created_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.describe_output.0.name"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.describe_output.0.kind", "WAREHOUSE"),

					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.parameters.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.parameters.0.max_concurrency_level.0.value", "8"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.parameters.0.statement_queued_timeout_in_seconds.0.value", "0"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.parameters.0.statement_timeout_in_seconds.0.value", "172800"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Warehouses/optionals_unset"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.state"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.type", string(sdk.WarehouseTypeStandard)),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.size", string(sdk.WarehouseSizeXSmall)),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.min_cluster_count", "1"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.max_cluster_count", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.started_clusters"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.running"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.queued"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.is_current", "true"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.auto_suspend", "600"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.auto_resume", "true"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.available"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.provisioning"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.quiescing"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.other"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.created_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.resumed_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.updated_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.enable_query_acceleration", "false"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.query_acceleration_max_scale_factor", "8"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.resource_monitor", ""),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.show_output.0.scaling_policy", string(sdk.ScalingPolicyStandard)),
					resource.TestCheckResourceAttrSet("data.snowflake_warehouses.test", "warehouses.0.show_output.0.owner_role_type"),

					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.describe_output.#", "0"),
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.0.parameters.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Warehouses_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	commonVariables := config.Variables{
		"name_1": config.StringVariable(idOne.Name()),
		"name_2": config.StringVariable(idTwo.Name()),
		"name_3": config.StringVariable(idThree.Name()),
	}

	likeConfig := config.Variables{
		"like": config.StringVariable(idOne.Name()),
	}
	maps.Copy(likeConfig, commonVariables)

	likeConfig2 := config.Variables{
		"like": config.StringVariable(prefix + "%"),
	}
	maps.Copy(likeConfig2, commonVariables)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Warehouses/like"),
				ConfigVariables: likeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.#", "1"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Warehouses/like"),
				ConfigVariables: likeConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_warehouses.test", "warehouses.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Warehouses_WarehouseNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Warehouses/without_warehouse"),
				ExpectError:     regexp.MustCompile("there should be at least one warehouse"),
			},
		},
	})
}
