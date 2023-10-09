package resources_test

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

func TestAcc_DynamicTable_basic(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_dynamic_table.dt"
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
	variableSet2 := m()
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")

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
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "warehouse", acc.TestWarehouseName),
					resource.TestCheckResourceAttr(resourceName, "target_lag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target_lag.0.maximum_duration", "2 minutes"),
					resource.TestCheckResourceAttr(resourceName, "query", fmt.Sprintf("select \"id\" from \"%v\".\"%v\".\"%v\"", acc.TestDatabaseName, acc.TestSchemaName, tableName)),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),

					// computed attributes

					// - not used at this time
					//  resource.TestCheckResourceAttrSet(resourceName, "cluster_by"),
					resource.TestCheckResourceAttrSet(resourceName, "rows"),
					resource.TestCheckResourceAttrSet(resourceName, "bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "owner"),
					resource.TestCheckResourceAttrSet(resourceName, "refresh_mode"),
					// - not used at this time
					// resource.TestCheckResourceAttrSet(resourceName, "automatic_clustering"),
					resource.TestCheckResourceAttrSet(resourceName, "scheduling_state"),
					resource.TestCheckResourceAttrSet(resourceName, "last_suspended_on"),
					resource.TestCheckResourceAttrSet(resourceName, "is_clone"),
					resource.TestCheckResourceAttrSet(resourceName, "is_replica"),
					resource.TestCheckResourceAttrSet(resourceName, "data_timestamp"),
				),
			},
			// test target lag to downstream and change comment

			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "target_lag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target_lag.0.downstream", "true"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},
			// test import
			{
				ConfigDirectory:   config.TestStepDirectory(),
				ConfigVariables:   variableSet2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
