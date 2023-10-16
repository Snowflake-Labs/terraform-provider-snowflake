package datasources_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DynamicTables_complete(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dataSourceName := "data.snowflake_dynamic_tables.dts"
	tableName := name + "_table"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(name),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"warehouse":  config.StringVariable(acc.TestWarehouseName),
			"query":      config.StringVariable(fmt.Sprintf("select \"id\" from \"%v\".\"%v\".\"%v\"", acc.TestDatabaseName, acc.TestSchemaName, tableName)),
			"comment":    config.StringVariable("Terraform acceptance test"),
			"table_name": config.StringVariable(tableName),
		}
	}
	variableSet1 := m()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDynamicTableDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "like.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "like.0.pattern", name),
					resource.TestCheckResourceAttr(dataSourceName, "in.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "in.0.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(dataSourceName, "starts_with", name),
					resource.TestCheckResourceAttr(dataSourceName, "limit.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "limit.0.rows", "1"),

					// computed attributes
					resource.TestCheckResourceAttr(dataSourceName, "records.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.created_on"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.schema_name"),
					// unused by Snowflake API at this time (always empty)
					// resource.TestCheckResourceAttrSet(dataSourceName, "records.0.cluster_by"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.rows"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.bytes"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.owner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.target_lag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.refresh_mode"),
					// unused by Snowflake API at this time (always empty)
					// resource.TestCheckResourceAttrSet(dataSourceName, "records.0.refresh_mode_reason"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.warehouse"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.comment"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.text"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.automatic_clustering"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.scheduling_state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.last_suspended_on"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.is_clone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.is_replica"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.data_timestamp"),
				),
			},
		},
	})
}

func testAccCheckDynamicTableDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_dynamic_table" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewSchemaObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["schema"], rs.Primary.Attributes["name"])
		dynamicTable, err := client.DynamicTables.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("dynamic table %v still exists", dynamicTable.Name)
		}
	}
	return nil
}
