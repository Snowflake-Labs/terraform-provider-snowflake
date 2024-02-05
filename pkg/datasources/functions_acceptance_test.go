package datasources_test

import (
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Functions(t *testing.T) {
	functionNameOne := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	functionNameTwo := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dataSourceName := "data.snowflake_functions.functions"

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":          config.StringVariable(acc.TestDatabaseName),
			"schema":            config.StringVariable(acc.TestSchemaName),
			"function_name_one": config.StringVariable(functionNameOne),
			"function_name_two": config.StringVariable(functionNameTwo),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Functions/complete"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(dataSourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttrSet(dataSourceName, "functions.#"),
				),
			},
		},
	})
}
