package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMonitors(t *testing.T) {
	resourceMonitorName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: resourceMonitors(resourceMonitorName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_resource_monitors.s", "resource_monitors.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_resource_monitors.s", "resource_monitors.0.name"),
				),
			},
		},
	})
}

func resourceMonitors(resourceMonitorName string) string {
	return fmt.Sprintf(`
	resource snowflake_resource_monitor "s"{
		name 		 = "%v"
		credit_quota = 5
	}

	data snowflake_resource_monitors "s" {
		depends_on = [snowflake_resource_monitor.s]
	}
	`, resourceMonitorName)
}
