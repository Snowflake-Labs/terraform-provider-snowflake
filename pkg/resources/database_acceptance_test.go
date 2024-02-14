package resources_test

import (
	"context"
	"fmt"
	"os"
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
	if _, ok := os.LookupEnv("SKIP_DATABASE_TESTS"); ok {
		t.Skip("Skipping TestAcc_DatabaseWithUnderscore")
	}

	prefix := "_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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
	if _, ok := os.LookupEnv("SKIP_DATABASE_TESTS"); ok {
		t.Skip("Skipping TestAcc_Database")
	}

	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	secondaryAccountName := getSecondaryAccount(t)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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
