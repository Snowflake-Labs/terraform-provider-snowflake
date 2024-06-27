package datasources_test

import (
	"maps"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Databases_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.ConfigureClientOnce)
	databaseName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()

	configVariables := config.Variables{
		"name":               config.StringVariable(databaseName),
		"comment":            config.StringVariable(comment),
		"account_identifier": config.StringVariable(acc.SecondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/optionals_set"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.name", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.kind", "STANDARD"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.transient", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.is_current", "true"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.origin", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.show_output.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.options", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.show_output.0.retention_time"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.resource_group", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.show_output.0.owner_role_type"),

					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.describe_output.#", "2"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.describe_output.0.created_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.describe_output.0.name"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.describe_output.0.kind", "SCHEMA"),

					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.data_retention_time_in_days.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.max_data_extension_time_in_days.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.external_volume.0.value", ""),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.catalog.0.value", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.replace_invalid_characters.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.0.default_ddl_collation.0.value", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.storage_serialization_policy.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.log_level.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.trace_level.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.suspend_task_after_num_failures.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.task_auto_retry_attempts.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.user_task_managed_initial_warehouse_size.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.user_task_minimum_trigger_interval_in_seconds.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.quoted_identifiers_ignore_case.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.parameters.0.enable_console_output.0.value"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/optionals_unset"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.name", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.kind", "STANDARD"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.transient", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.is_current", "true"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.origin", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.show_output.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.options", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.show_output.0.retention_time"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.show_output.0.resource_group", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.test", "databases.0.show_output.0.owner_role_type"),

					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.describe_output.#", "0"),
					resource.TestCheckResourceAttr("data.snowflake_databases.test", "databases.0.parameters.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Databases_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Databases/without_database"),
				ExpectError:     regexp.MustCompile("there should be at least one database"),
			},
		},
	})
}
