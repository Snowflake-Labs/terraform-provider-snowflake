package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Procedures(t *testing.T) {
	procNameOne := acc.TestClient().Ids.Alpha()
	procNameTwo := acc.TestClient().Ids.Alpha()
	dataSourceName := "data.snowflake_procedures.procedures"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":      config.StringVariable(acc.TestDatabaseName),
			"schema":        config.StringVariable(acc.TestSchemaName),
			"proc_name_one": config.StringVariable(procNameOne),
			"proc_name_two": config.StringVariable(procNameTwo),
		}
	}
	configVariables := m()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Procedures/complete"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(dataSourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttrSet(dataSourceName, "procedures.#"),
					// resource.TestCheckResourceAttr(dataSourceName, "procedures.#", "3"),
					// Extra 1 in procedure count above due to ASSOCIATE_SEMANTIC_CATEGORY_TAGS appearing in all "SHOW PROCEDURES IN ..." commands
				),
			},
		},
	})
}
