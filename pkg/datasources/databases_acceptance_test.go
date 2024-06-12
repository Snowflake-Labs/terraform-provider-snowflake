package datasources_test

import (
	"maps"
	"regexp"
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
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/optionals-set"),
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

					resource.TestCheckResourceAttrWith("data.snowflake_databases.test", "databases.0.parameters.#", acc.IsGreaterOrEqualTo(10)),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.key"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.value", ""),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.default", ""),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.level", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.description"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/optionals-unset"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.name", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.kind", "STANDARD"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.is_transient", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.is_current", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.origin", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.comment", comment),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.options", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.retention_time"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.resource_group", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.owner_role_type"),

					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.description.#", "0"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Databases_DifferentFiltering(t *testing.T) {
	prefix := acc.TestClient().Ids.Alpha()
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	commonVariables := config.Variables{
		"name_1": config.StringVariable(idOne.Name()),
		"name_2": config.StringVariable(idTwo.Name()),
		"name_3": config.StringVariable(idThree.Name()),
	}

	likeConfig := config.Variables{
		"like": config.StringVariable(idOne.Name()),
	}
	maps.Copy(likeConfig, commonVariables)

	startsWithConfig := config.Variables{
		"starts_with": config.StringVariable(prefix),
	}
	maps.Copy(startsWithConfig, commonVariables)

	limitConfig := config.Variables{
		"rows": config.IntegerVariable(1),
		"from": config.StringVariable(prefix),
	}
	maps.Copy(limitConfig, commonVariables)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/like"),
				ConfigVariables: likeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.#", "1"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/starts_with"),
				ConfigVariables: startsWithConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.#", "2"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/limit"),
				ConfigVariables: limitConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.#", "1"),
				),
			},
		},
	})
}

func TestAcc_Databases_DatabaseNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/without-database"),
				ExpectError:     regexp.MustCompile("there should be at least one database"),
			},
		},
	})
}
