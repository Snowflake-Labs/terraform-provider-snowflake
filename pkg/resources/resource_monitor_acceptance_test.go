package resources_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ResourceMonitor(t *testing.T) {
	// TODO test more attributes
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: resourceMonitorConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "credit_quota", "100"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "false"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: resourceMonitorConfig2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "credit_quota", "150"),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "set_for_account", "false"),
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
resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	credit_quota    = 100
	set_for_account = false
}
`, accName)
}

func resourceMonitorConfig2(accName string) string {
	return fmt.Sprintf(`
resource "snowflake_resource_monitor" "test" {
	name            = "%v"
	credit_quota    = 150
	set_for_account = false
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
		Providers:    providers(),
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
