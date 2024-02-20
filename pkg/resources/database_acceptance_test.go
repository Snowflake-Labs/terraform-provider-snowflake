package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_DatabaseWithUnderscore(t *testing.T) {
	prefix := "_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: dbConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_database.db", "data_retention_time_in_days"),
				),
			},
		},
	})
}

func TestAcc_Database(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	secondaryAccountName := getSecondaryAccount(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: dbConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_database.db", "data_retention_time_in_days"),
				),
			},
			// RENAME
			{
				Config: dbConfig(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_database.db", "data_retention_time_in_days"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: dbConfig2(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment 2"),
					resource.TestCheckResourceAttr("snowflake_database.db", "data_retention_time_in_days", "3"),
				),
			},
			// ADD REPLICATION
			// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2369 error
			{
				Config: dbConfigWithReplication(prefix2, secondaryAccountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment 2"),
					resource.TestCheckResourceAttr("snowflake_database.db", "data_retention_time_in_days", "3"),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_configuration.0.accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_configuration.0.accounts.0", secondaryAccountName),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_database.db",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_configuration"},
			},
		},
	})
}

func TestAcc_DatabaseRemovedOutsideOfTerraform(t *testing.T) {
	id := generateUnsafeExecuteTestDatabaseName(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"db": config.StringVariable(id),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", id),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				PreConfig:       func() { dropDatabaseOutsideTerraform(t, id) },
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"db": config.StringVariable(id),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", id),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2021
func TestAcc_Database_issue2021(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	secondaryAccountName := getSecondaryAccount(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: dbConfigWithReplication(name, secondaryAccountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", name),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_configuration.0.accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_configuration.0.accounts.0", secondaryAccountName),
					testAccCheckIfDatabaseIsReplicated(t, name),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Database_DefaultDataRetentionTime(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	id := sdk.NewAccountObjectIdentifier(databaseName)

	configVariablesWithoutDatabaseDataRetentionTime := func(accountDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                    config.StringVariable(databaseName),
			"account_data_retention_time": config.IntegerVariable(accountDataRetentionTime),
		}
	}

	configVariablesWithDatabaseDataRetentionTime := func(databaseDataRetentionTime int, schemaDataRetentionTime int) config.Variables {
		vars := configVariablesWithoutDatabaseDataRetentionTime(databaseDataRetentionTime)
		vars["database_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithoutDatabase"),
				ConfigVariables: configVariablesWithoutDatabaseDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "-1"),
					checkAccountAndDatabaseDataRetentionTime(id, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithoutDatabase"),
				ConfigVariables: configVariablesWithoutDatabaseDataRetentionTime(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "-1"),
					checkAccountAndDatabaseDataRetentionTime(id, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithDatabase"),
				ConfigVariables: configVariablesWithDatabaseDataRetentionTime(10, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "5"),
					checkAccountAndDatabaseDataRetentionTime(id, 10, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithDatabase"),
				ConfigVariables: configVariablesWithDatabaseDataRetentionTime(10, 15),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "15"),
					checkAccountAndDatabaseDataRetentionTime(id, 10, 15),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithoutDatabase"),
				ConfigVariables: configVariablesWithoutDatabaseDataRetentionTime(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "-1"),
					checkAccountAndDatabaseDataRetentionTime(id, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithDatabase"),
				ConfigVariables: configVariablesWithDatabaseDataRetentionTime(10, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "0"),
					checkAccountAndDatabaseDataRetentionTime(id, 10, 0),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithDatabase"),
				ConfigVariables: configVariablesWithDatabaseDataRetentionTime(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "3"),
					checkAccountAndDatabaseDataRetentionTime(id, 10, 3),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Database_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	id := sdk.NewAccountObjectIdentifier(databaseName)

	configVariablesWithoutDatabaseDataRetentionTime := func(accountDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                    config.StringVariable(databaseName),
			"account_data_retention_time": config.IntegerVariable(accountDataRetentionTime),
		}
	}

	configVariablesWithDatabaseDataRetentionTime := func(databaseDataRetentionTime int, schemaDataRetentionTime int) config.Variables {
		vars := configVariablesWithoutDatabaseDataRetentionTime(databaseDataRetentionTime)
		vars["database_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithoutDatabase"),
				ConfigVariables: configVariablesWithoutDatabaseDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "-1"),
					checkAccountAndDatabaseDataRetentionTime(id, 5, 5),
				),
			},
			{
				PreConfig:       setDatabaseDataRetentionTime(t, id, 20),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithoutDatabase"),
				ConfigVariables: configVariablesWithoutDatabaseDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "-1"),
					checkAccountAndDatabaseDataRetentionTime(id, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database_DefaultDataRetentionTime/WithDatabase"),
				ConfigVariables: configVariablesWithDatabaseDataRetentionTime(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "3"),
					checkAccountAndDatabaseDataRetentionTime(id, 10, 3),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func dbConfig(prefix string) string {
	s := `
resource "snowflake_database" "db" {
	name = "%s"
	comment = "test comment"
}
`
	return fmt.Sprintf(s, prefix)
}

func dbConfig2(prefix string) string {
	s := `
resource "snowflake_database" "db" {
	name = "%s"
	comment = "test comment 2"
	data_retention_time_in_days = 3
}
`
	return fmt.Sprintf(s, prefix)
}

func dbConfigWithReplication(prefix string, secondaryAccountName string) string {
	s := `
resource "snowflake_database" "db" {
	name = "%s"
	comment = "test comment 2"
	data_retention_time_in_days = 3
	replication_configuration {
		accounts = [
			"%s"
		]
	}
}
`
	return fmt.Sprintf(s, prefix, secondaryAccountName)
}

func dropDatabaseOutsideTerraform(t *testing.T, id string) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()

	err = client.Databases.Drop(ctx, sdk.NewAccountObjectIdentifier(id), &sdk.DropDatabaseOptions{})
	require.NoError(t, err)
}

func getSecondaryAccount(t *testing.T) string {
	t.Helper()

	secondaryConfig, err := sdk.ProfileConfig("secondary_test_account")
	require.NoError(t, err)

	secondaryClient, err := sdk.NewClient(secondaryConfig)
	require.NoError(t, err)

	ctx := context.Background()

	account, err := secondaryClient.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)

	return account
}

func testAccCheckDatabaseExistence(t *testing.T, id string, shouldExist bool) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		client, err := sdk.NewDefaultClient()
		require.NoError(t, err)
		ctx := context.Background()

		_, err = client.Databases.ShowByID(ctx, sdk.NewAccountObjectIdentifier(id))
		if shouldExist {
			if err != nil {
				return fmt.Errorf("error while retrieving database %s, err = %w", id, err)
			}
		} else {
			if err == nil {
				return fmt.Errorf("database %v still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckIfDatabaseIsReplicated(t *testing.T, id string) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		client, err := sdk.NewDefaultClient()
		if err != nil {
			return err
		}

		ctx := context.Background()
		replicationDatabases, err := client.ReplicationFunctions.ShowReplicationDatabases(ctx, nil)
		if err != nil {
			return err
		}

		var exists bool
		for _, o := range replicationDatabases {
			if o.Name == id {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("database %s should be replicated", id)
		}

		return nil
	}
}

func checkAccountAndDatabaseDataRetentionTime(id sdk.AccountObjectIdentifier, expectedAccountRetentionDays int, expectedDatabaseRetentionsDays int) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		db := acc.TestAccProvider.Meta().(*sql.DB)
		client := sdk.NewClientFromDB(db)
		ctx := context.Background()

		database, err := client.Databases.ShowByID(ctx, id)
		if err != nil {
			return err
		}

		if database.RetentionTime != expectedDatabaseRetentionsDays {
			return fmt.Errorf("invalid database retention time, expected: %d, got: %d", expectedDatabaseRetentionsDays, database.RetentionTime)
		}

		param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterDataRetentionTimeInDays)
		if err != nil {
			return err
		}
		accountRetentionDays, err := strconv.Atoi(param.Value)
		if err != nil {
			return err
		}

		if accountRetentionDays != expectedAccountRetentionDays {
			return fmt.Errorf("invalid account retention time, expected: %d, got: %d", expectedAccountRetentionDays, accountRetentionDays)
		}

		return nil
	}
}

func setDatabaseDataRetentionTime(t *testing.T, id sdk.AccountObjectIdentifier, days int) func() {
	t.Helper()

	return func() {
		client, err := sdk.NewDefaultClient()
		require.NoError(t, err)
		ctx := context.Background()

		err = client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			Set: &sdk.DatabaseSet{
				DataRetentionTimeInDays: sdk.Int(days),
			},
		})
		require.NoError(t, err)
	}
}
