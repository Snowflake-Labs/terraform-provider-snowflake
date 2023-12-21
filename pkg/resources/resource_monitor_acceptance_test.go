package resources_test

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAcc_ResourceMonitor(t *testing.T) {
	// TODO test more attributes
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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

func TestAcc_ResourceMonitorChangeStartEndTimestamp(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
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
  name           = "test"
  comment        = "foo"
  warehouse_size = "XSMALL"
}

resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	frequency 		= "WEEKLY"
	start_timestamp = "2055-01-01 12:00"
	end_timestamp = "2056-01-01 12:00"

}
`, accName)
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
	userEnv := os.Getenv("RESOURCE_MONITOR_NOTIFY_USERS_TEST")
	if userEnv == "" {
		t.Skip("Skipping ResourceMonitorUpdateNotifyUsers")
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
	empty := []string{}
	emptyUsersConfig, err := resourceMonitorNotifyUsersConfig(name, empty)
	if err != nil {
		t.Error(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
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

// TestAcc_ResourceMonitor_issue2167 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2167 issue.
// Second step is purposely error, because tests TestAcc_ResourceMonitorUpdateNotifyUsers and TestAcc_ResourceMonitorNotifyUsers are still skipped.
// It can be fixed with them.
func TestAcc_ResourceMonitor_issue2167(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configNoUsers, err := resourceMonitorNotifyUsersConfig(name, []string{})
	require.NoError(t, err)
	config, err := resourceMonitorNotifyUsersConfig(name, []string{"non_existing_user"})
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
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
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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
