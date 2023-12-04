package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO: test empty database/schema combinations
// TODO: test with pattern
func TestAcc_Alerts(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: alertsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
		},
	})
}

func alertsConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_alert" "test_resource_alert" {
	name     	      = "%s"
	database  	      = "%s"
	schema   	      = "%s"
	warehouse 	      = "%s"
	condition         = "select 0 as c"
	action            = "select 0 as c"
	enabled  	      = false
	comment           = "some comment"
	alert_schedule 	  {
		interval = "60"
	}
}

data "snowflake_alerts" "test_datasource_alert" {
}
`, name, acc.TestDatabaseName, acc.TestSchemaName, acc.TestWarehouseName)
}
