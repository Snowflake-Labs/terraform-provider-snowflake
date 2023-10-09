package resources_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ResourceMonitor(t *testing.T) {
	// TODO test more attributes
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: resourceMonitorConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "credit_quota", "100"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "false"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "notify_triggers.0", "40"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "suspend_trigger", "80"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "suspend_immediate_trigger", "90"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: resourceMonitorConfig2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "credit_quota", "150"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "true"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "notify_triggers.0", "50"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "suspend_trigger", "75"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "suspend_immediate_trigger", "95"),
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

func resourceMonitorConfig(accName string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "warehouse" {
  name           = "test"
  comment        = "foo"
  warehouse_size = "SMALL"
}

resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	credit_quota    = 100
	set_for_account = false
 	notify_triggers = [40]
	suspend_trigger = 80
	suspend_immediate_trigger = 90
	warehouses      = [snowflake_warehouse.warehouse.id]
}
`, accName)
}

func resourceMonitorConfig2(accName string) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse" "warehouse" {
  name           = "test"
  comment        = "foo"
  warehouse_size = "SMALL"
}

resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	credit_quota    = 150
	set_for_account = true
	notify_triggers = [50]
	warehouses      = []
	suspend_trigger = 75
	suspend_immediate_trigger = 95
}
`, accName)
}

func TestAcc_ResourceMonitorNotifyUsers(t *testing.T) {
	userEnv := os.Getenv("RESOURCE_MONITOR_NOTIFY_USERS_TEST")
	if userEnv == "" {
		t.Skip("Skipping TestAcc_ResourceMonitorNotifyUsers")
	}
	users := strings.Split(userEnv, ",")
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
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
