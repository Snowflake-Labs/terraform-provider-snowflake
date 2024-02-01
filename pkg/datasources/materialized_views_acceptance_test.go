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
	warehouseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
				Config: materializedViews(warehouseName, databaseName, schemaName, tableName, viewName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_materialized_views.v", "schema", schemaName),
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
	resource "snowflake_warehouse" "w" {
		name                = "%v"
		initially_suspended = false
	}

	resource snowflake_database "d" {
		name = "%v"
	}

	resource snowflake_schema "s"{
		name 	 = "%v"
		database = snowflake_database.d.name
	}

	resource snowflake_table "t"{
		name 	 = "%v"
		database = snowflake_schema.s.database
		schema 	 = snowflake_schema.s.name
		column {
			name = "column2"
			type = "VARCHAR(16)"
		}
	}

	resource snowflake_materialized_view "v"{
		name 	   = "%v"
		comment    = "Terraform test resource"
		database   = snowflake_schema.s.database
		schema 	   = snowflake_schema.s.name
		is_secure  = true
		or_replace = false
		statement  = "SELECT * FROM ${snowflake_table.t.name}"
		warehouse  = snowflake_warehouse.w.name
	}

	data snowflake_materialized_views "v" {
		database = snowflake_materialized_view.v.database
		schema = snowflake_materialized_view.v.schema
		depends_on = [snowflake_materialized_view.v]
	}
	`, warehouseName, databaseName, schemaName, tableName, viewName)
}
