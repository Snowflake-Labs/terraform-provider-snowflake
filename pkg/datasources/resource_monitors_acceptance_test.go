package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ResourceMonitors(t *testing.T) {
	resourceMonitorName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
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
