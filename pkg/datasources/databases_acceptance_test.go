package datasources_test

import (
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Databases_Complete(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()

	configVariables := config.Variables{
		"name":               config.StringVariable(databaseName),
		"comment":            config.StringVariable(comment),
		"account_identifier": config.StringVariable(strconv.Quote(acc.SecondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName())),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StandardDatabase),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.name", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.kind", "STANDARD"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.is_transient", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.is_current", "true"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.origin", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.comment", comment),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.options", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.retention_time"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.resource_group", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.owner_role_type"),

					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.description.#", "2"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.description.0.created_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.description.0.name"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.description.0.kind", "SCHEMA"),

					acc.TestCheckResourceAttrNumberAtLeast("data.snowflake_databases.test", "databases.0.parameters.#", 10),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.key"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.value", ""),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.default", ""),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.level", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.description"),
				),
			},
		},
	})
}
