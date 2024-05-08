package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Alerts(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: alertsResourceConfig(name) + alertsDatasourceConfigNoOptionals(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
			{
				Config: alertsResourceConfig(name) + alertsDatasourceConfigDbOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
			{
				Config: alertsResourceConfig(name) + alertsDatasourceConfigDbAndSchema(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
					resource.TestCheckResourceAttr("data.snowflake_alerts.test_datasource_alert", "alerts.0.name", name),
				),
			},
			{
				Config: alertsResourceConfig(name) + alertsDatasourceConfigAllOptionals(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
					resource.TestCheckResourceAttr("data.snowflake_alerts.test_datasource_alert", "alerts.0.name", name),
				),
			},
			{
				Config: alertsResourceConfig(name) + alertsDatasourceConfigSchemaOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
		},
	})
}

func alertsResourceConfig(name string) string {
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
`, name, acc.TestDatabaseName, acc.TestSchemaName, acc.TestWarehouseName)
}

func alertsDatasourceConfigNoOptionals() string {
	return `
data "snowflake_alerts" "test_datasource_alert" {}
`
}

func alertsDatasourceConfigDbOnly() string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	database  	      = "%s"
}
`, acc.TestDatabaseName)
}

func alertsDatasourceConfigDbAndSchema() string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	database  	      = "%s"
	schema  	      = "%s"
}
`, acc.TestDatabaseName, acc.TestSchemaName)
}

func alertsDatasourceConfigAllOptionals(name string) string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	database  	      = "%s"
	schema  	      = "%s"
	pattern  	      = "%s"
}
`, acc.TestDatabaseName, acc.TestSchemaName, name)
}

func alertsDatasourceConfigSchemaOnly() string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	schema  	      = "%s"
}
`, acc.TestSchemaName)
}
