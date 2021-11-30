package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ResourceMonitor(t *testing.T) {
	// TODO test more attributes
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: resourceMonitorConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "credit_quota", "100"),
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
