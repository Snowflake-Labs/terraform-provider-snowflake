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

func TestAcc_MaterializedViews(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: materializedViews(acc.TestWarehouseName, acc.TestDatabaseName, acc.TestSchemaName, tableName, viewName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_materialized_views.v", "materialized_views.#"),
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "materialized_views.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "materialized_views.0.name", viewName),
				),
			},
		},
	})
}

func materializedViews(warehouseName string, databaseName string, schemaName string, tableName string, viewName string) string {
	return fmt.Sprintf(`
	resource snowflake_table "t"{
		name 	 = "%[4]v"
		database = "%[2]s"
		schema 	 = "%[3]s"
		column {
			name = "column2"
			type = "VARCHAR(16)"
		}
	}

	resource snowflake_materialized_view "v"{
		name 	   = "%[5]v"
		comment    = "Terraform test resource"
		database   = "%[2]s"
		schema 	   = "%[3]s"
		is_secure  = true
		or_replace = false
		statement  = "SELECT * FROM ${snowflake_table.t.name}"
		warehouse  = "%[1]s"
	}

	data snowflake_materialized_views "v" {
		database = "%[2]s"
		schema = "%[3]s"
		depends_on = [snowflake_materialized_view.v]
	}
	`, warehouseName, databaseName, schemaName, tableName, viewName)
}
