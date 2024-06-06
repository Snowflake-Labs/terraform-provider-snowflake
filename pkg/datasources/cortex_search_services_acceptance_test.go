package datasources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CortexSearchServices_complete(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	dataSourceName := "data.snowflake_cortex_search_services.csss"
	tableName := name + "_table"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(name),
			"on":         config.StringVariable("id"),
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
		CheckDestroy: nil,
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
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.comment"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "records.#", "1"),
				),
			},
		},
	})
}

func TestAcc_CortexSearchServices_badCombination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      cortexSearchServicesDatasourceConfigDbAndSchema(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_CortexSearchServices_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      cortexSearchServicesDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func cortexSearchServicesDatasourceConfigDbAndSchema() string {
	return fmt.Sprintf(`
data "snowflake_cortex_search_services" "csss" {
  in {
    database = "%s"
    schema   = "%s"
  }
}
`, acc.TestDatabaseName, acc.TestSchemaName)
}

func cortexSearchServicesDatasourceEmptyIn() string {
	return `
data "snowflake_cortex_search_services" "csss" {
  in {
  }
}
`
}
