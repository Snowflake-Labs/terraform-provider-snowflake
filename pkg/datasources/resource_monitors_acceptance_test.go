package datasources_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ResourceMonitors(t *testing.T) {
	prefix := "data_source_resource_monitor_"
	resourceMonitorName := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	resourceMonitorName2 := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: resourceMonitors(resourceMonitorName.Name(), resourceMonitorName2.Name(), prefix+"%"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_resource_monitors.test", "resource_monitors.#", "2"),
				),
			},
			{
				Config: resourceMonitors(resourceMonitorName.Name(), resourceMonitorName2.Name(), resourceMonitorName.Name()),
				Check: assert.AssertThat(t,
					resourceshowoutputassert.ResourceMonitorDatasourceShowOutput(t, "snowflake_resource_monitors.test").
						HasName(resourceMonitorName.Name()).
						HasCreditQuota(5).
						HasUsedCredits(0).
						HasRemainingCredits(5).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTime("").
						HasSuspendAt(0).
						HasSuspendImmediateAt(0).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
		},
	})
}

func resourceMonitors(resourceMonitorName, resourceMonitorName2, searchPrefix string) string {
	return fmt.Sprintf(`
	resource "snowflake_resource_monitor" "rm1" {
		name 		 = "%s"
		credit_quota = 5
	}

	resource "snowflake_resource_monitor" "rm2" {
		name 		 = "%s"
		credit_quota = 15
	}

	data "snowflake_resource_monitors" "test" {
		depends_on = [ snowflake_resource_monitor.rm1, snowflake_resource_monitor.rm2 ]
		like = "%s"
	}
	`, resourceMonitorName, resourceMonitorName2, searchPrefix)
}
