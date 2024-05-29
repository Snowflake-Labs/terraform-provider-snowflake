package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Alerts(t *testing.T) {
	alertId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigNoOptionals(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigDbOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigDbAndSchema(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
					resource.TestCheckResourceAttr("data.snowflake_alerts.test_datasource_alert", "alerts.0.name", alertId.Name()),
				),
			},
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigAllOptionals(alertId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
					resource.TestCheckResourceAttr("data.snowflake_alerts.test_datasource_alert", "alerts.0.name", alertId.Name()),
				),
			},
			{
				Config: alertsResourceConfig(alertId) + alertsDatasourceConfigSchemaOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_alerts.test_datasource_alert", "alerts.#"),
				),
			},
		},
	})
}

func alertsResourceConfig(alertId sdk.SchemaObjectIdentifier) string {
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
`, alertId.Name(), alertId.DatabaseName(), alertId.SchemaName(), acc.TestWarehouseName)
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

func alertsDatasourceConfigAllOptionals(alertId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	database  	      = "%s"
	schema  	      = "%s"
	pattern  	      = "%s"
}
`, alertId.DatabaseName(), alertId.SchemaName(), alertId.Name())
}

func alertsDatasourceConfigSchemaOnly() string {
	return fmt.Sprintf(`
data "snowflake_alerts" "test_datasource_alert" {
	schema  	      = "%s"
}
`, acc.TestSchemaName)
}
