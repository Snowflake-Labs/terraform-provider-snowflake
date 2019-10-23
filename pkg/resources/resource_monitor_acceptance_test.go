package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceMonitor(t *testing.T) {
	// TODO test more attributes
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: resourceMonitorConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", name),
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
	name           = "%v"
}
`, accName)
}
