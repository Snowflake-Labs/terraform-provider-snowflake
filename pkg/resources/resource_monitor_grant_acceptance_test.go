package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMonitor_defaults(t *testing.T) {
	wName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	roleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: resourceMonitorGrantConfig(wName, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor_grant.test", "monitor_name", wName),
					resource.TestCheckResourceAttr("snowflake_resource_monitor_grant.test", "privilege", "MONITOR"),
				),
			},
		},
	})
}

func resourceMonitorGrantConfig(n, role string) string {
	return fmt.Sprintf(`

resource "snowflake_resource_monitor" "test" {
  name      = "%v"
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_resource_monitor_grant" "test" {
  monitor_name = snowflake_resource_monitor.test.name
  roles          = [snowflake_role.test.name]
}
`, n, role)
}
