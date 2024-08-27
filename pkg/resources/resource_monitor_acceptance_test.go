package resources_test

import (
	"encoding/json"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_ResourceMonitor(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("-1").
						HasNotifyUsersLen(0).
						HasFrequencyString(string(sdk.FrequencyMonthly)).
						HasNoStartTimestamp(). // TODO: Should be end of the last month
						HasNoEndTimestamp().
						HasTriggerLen(0),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test"),
				),
				//Check: resource.ComposeTestCheckFunc(
				//	resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", id.Name()),
				//	resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "fully_qualified_name", id.Name()),
				//	resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "notify_users.#", "0"),
				//	resource.TestCheckNoResourceAttr("snowflake_resource_monitor.test", "credit_quota"),
				//	resource.TestCheckNoResourceAttr("snowflake_resource_monitor.test", "frequency"),
				//	resource.TestCheckNoResourceAttr("snowflake_resource_monitor.test", "start_timestamp"),
				//	resource.TestCheckNoResourceAttr("snowflake_resource_monitor.test", "end_timestamp"),
				//	resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "trigger.#", "0"),
				//),
			},
			// CHANGE PROPERTIES
			//{
			//	Config: resourceMonitorConfig2(id.Name(), 75),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", id.Name()),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "fully_qualified_name", id.FullyQualifiedName()),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "credit_quota", "150"),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "true"),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "notify_triggers.0", "50"),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "suspend_trigger", "75"),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "suspend_immediate_trigger", "95"),
			//	),
			//},
			//// CHANGE JUST suspend_trigger; proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2316
			//{
			//	Config: resourceMonitorConfig2(id.Name(), 60),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", id.Name()),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "fully_qualified_name", id.FullyQualifiedName()),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "credit_quota", "150"),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "true"),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "notify_triggers.0", "50"),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "suspend_trigger", "60"),
			//		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "suspend_immediate_trigger", "95"),
			//	),
			//},
			//// IMPORT
			//{
			//	ResourceName:      "snowflake_resource_monitor.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}

func TestAcc_ResourceMonitorChangeStartEndTimestamp(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: resourceMonitorConfigInitialTimestamp(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "frequency", "WEEKLY"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "start_timestamp", "2050-01-01 12:00"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "end_timestamp", "2055-01-01 12:00"),
				),
			},
			{
				Config: resourceMonitorConfigUpdatedTimestamp(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "frequency", "WEEKLY"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "start_timestamp", "2055-01-01 12:00"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "end_timestamp", "2056-01-01 12:00"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_resource_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceMonitorConfigUpdatedTimestamp(accName string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "warehouse" {
  name           = "test%v"
  comment        = "foo"
  warehouse_size = "XSMALL"
}

resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	frequency 		= "WEEKLY"
	start_timestamp = "2055-01-01 12:00"
	end_timestamp = "2056-01-01 12:00"

}
`, accName, accName)
}

// fix 2 added empy notifiy user
// Config for changed timestamp frequency validation test
func resourceMonitorConfigInitialTimestamp(accName string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "warehouse" {
  name           = "test"
  comment        = "foo"
  warehouse_size = "XSMALL"
}

resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	frequency 		= "WEEKLY"
	start_timestamp = "2050-01-01 12:00"
	end_timestamp = "2055-01-01 12:00"

}
`, accName)
}

func TestAcc_ResourceMonitorUpdateNotifyUsers(t *testing.T) {
	userEnv := testenvs.GetOrSkipTest(t, testenvs.ResourceMonitorNotifyUsers)
	users := strings.Split(userEnv, ",")
	name := acc.TestClient().Ids.Alpha()
	config, err := resourceMonitorNotifyUsersConfig(name, users)
	if err != nil {
		t.Error(err)
	}
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "false"),
	}
	for _, s := range users {
		checks = append(checks, resource.TestCheckTypeSetElemAttr("snowflake_resource_monitor.test", "notify_users.*", s))
	}
	empty := []string{}
	emptyUsersConfig, err := resourceMonitorNotifyUsersConfig(name, empty)
	if err != nil {
		t.Error(err)
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: emptyUsersConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "false"),
				),
			},
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				ResourceName:      "snowflake_resource_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceMonitorConfig(accName string, warehouse string) string {
	return fmt.Sprintf(`
resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	credit_quota    = 100
	set_for_account = false
 	notify_triggers = [40]
	suspend_trigger = 80
	suspend_immediate_trigger = 90
	warehouses      = ["%s"]
}
`, accName, warehouse)
}

func resourceMonitorConfig2(accName string, suspendTrigger int) string {
	return fmt.Sprintf(`
resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	credit_quota    = 150
	set_for_account = true
	notify_triggers = [50]
	warehouses      = []
	suspend_trigger = %d
	suspend_immediate_trigger = 95
}
`, accName, suspendTrigger)
}

// TestAcc_ResourceMonitor_issue2167 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2167 issue.
// Second step is purposely error, because tests TestAcc_ResourceMonitorUpdateNotifyUsers and TestAcc_ResourceMonitorNotifyUsers are still skipped.
// It can be fixed with them.
func TestAcc_ResourceMonitor_issue2167(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	configNoUsers, err := resourceMonitorNotifyUsersConfig(name, []string{})
	require.NoError(t, err)
	config, err := resourceMonitorNotifyUsersConfig(name, []string{"non_existing_user"})
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: configNoUsers,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
				),
			},
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`.*090268 \(22023\): User non_existing_user does not exist.*`),
			},
		},
	})
}

func TestAcc_ResourceMonitorNotifyUsers(t *testing.T) {
	userEnv := testenvs.GetOrSkipTest(t, testenvs.ResourceMonitorNotifyUsers)
	users := strings.Split(userEnv, ",")
	name := acc.TestClient().Ids.Alpha()
	config, err := resourceMonitorNotifyUsersConfig(name, users)
	if err != nil {
		t.Error(err)
	}
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
		resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "false"),
	}
	for _, s := range users {
		checks = append(checks, resource.TestCheckTypeSetElemAttr("snowflake_resource_monitor.test", "notify_users.*", s))
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				ResourceName:      "snowflake_resource_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceMonitorNotifyUsersConfig(accName string, accNotifyUsers []string) (string, error) {
	notifyUsers, err := json.Marshal(accNotifyUsers)
	if err != nil {
		return "", err
	}
	config := fmt.Sprintf(`
resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	set_for_account = false
	notify_users    = %v
}
`, accName, string(notifyUsers))
	return config, nil
}
