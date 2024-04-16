package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Schema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := "Terraform acceptance test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(name),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
					checkBool("snowflake_schema.test", "is_transient", false),
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
			// UPDATE COMMENT (proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2606)
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(name),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(""),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", ""),
					checkBool("snowflake_schema.test", "is_transient", false),
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

func TestAcc_Schema_Rename(t *testing.T) {
	oldSchemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newSchemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := "Terraform acceptance test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(oldSchemaName),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", oldSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
					checkBool("snowflake_schema.test", "is_transient", false),
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(newSchemaName),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", newSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
					checkBool("snowflake_schema.test", "is_transient", false),
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

// TestAcc_Schema_TwoSchemasWithTheSameNameOnDifferentDatabases proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2209 issue.
func TestAcc_Schema_TwoSchemasWithTheSameNameOnDifferentDatabases(t *testing.T) {
	name := "test_schema"
	// It seems like Snowflake orders the output of SHOW command based on names, so they do matter
	newDatabaseName := "SELDQBXEKC"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(name),
					"database": config.StringVariable(acc.TestDatabaseName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":         config.StringVariable(name),
					"database":     config.StringVariable(acc.TestDatabaseName),
					"new_database": config.StringVariable(newDatabaseName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test_2", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test_2", "database", newDatabaseName),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Schema_DefaultDataRetentionTime(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	id := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)

	configVariablesWithoutSchemaDataRetentionTime := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseName),
			"schema":                       config.StringVariable(schemaName),
			"database_data_retention_time": config.IntegerVariable(databaseDataRetentionTime),
		}
	}

	configVariablesWithSchemaDataRetentionTime := func(databaseDataRetentionTime int, schemaDataRetentionTime int) config.Variables {
		vars := configVariablesWithoutSchemaDataRetentionTime(databaseDataRetentionTime)
		vars["schema_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(id, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(id, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "5"),
					checkDatabaseAndSchemaDataRetentionTime(id, 10, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 15),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "15"),
					checkDatabaseAndSchemaDataRetentionTime(id, 10, 15),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(id, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "0"),
					checkDatabaseAndSchemaDataRetentionTime(id, 10, 0),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "3"),
					checkDatabaseAndSchemaDataRetentionTime(id, 10, 3),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Schema_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	id := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)

	configVariablesWithoutSchemaDataRetentionTime := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseName),
			"schema":                       config.StringVariable(schemaName),
			"database_data_retention_time": config.IntegerVariable(databaseDataRetentionTime),
		}
	}

	configVariablesWithSchemaDataRetentionTime := func(databaseDataRetentionTime int, schemaDataRetentionTime int) config.Variables {
		vars := configVariablesWithoutSchemaDataRetentionTime(databaseDataRetentionTime)
		vars["schema_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(id, 5, 5),
				),
			},
			{
				PreConfig:       setSchemaDataRetentionTime(t, id, 20),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(id, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "3"),
					checkDatabaseAndSchemaDataRetentionTime(id, 10, 3),
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

func TestAcc_Schema_RemoveDatabaseOutsideOfTerraform(t *testing.T) {
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := map[string]config.Variable{
		"schema_name":   config.StringVariable(schemaName),
		"database_name": config.StringVariable(acc.TestDatabaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_RemoveOutsideOfTerraform"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					removeSchemaOutsideOfTerraform(t, acc.TestDatabaseName, schemaName)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				RefreshPlanChecks: resource.RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
						expectsCreatePlan("snowflake_schema.test"),
					},
				},
			},
		},
	})
}

func TestAcc_Schema_RemoveSchemaOutsideOfTerraform(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := map[string]config.Variable{
		"schema_name":   config.StringVariable(schemaName),
		"database_name": config.StringVariable(databaseName),
	}

	cleanupDatabase := createDatabaseOutsideTerraform(t, databaseName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_RemoveOutsideOfTerraform"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					cleanupDatabase()
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_RemoveOutsideOfTerraform"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("error creating schema"),
			},
		},
	})
}

func checkDatabaseAndSchemaDataRetentionTime(id sdk.DatabaseObjectIdentifier, expectedDatabaseRetentionsDays int, expectedSchemaRetentionDays int) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		schema, err := client.Schemas.ShowByID(ctx, id)
		if err != nil {
			return err
		}

		database, err := client.Databases.ShowByID(ctx, sdk.NewAccountObjectIdentifier(id.DatabaseName()))
		if err != nil {
			return err
		}

		// "retention_time" may sometimes be an empty string instead of an integer
		var schemaRetentionTime int64
		{
			rt := schema.RetentionTime
			if rt == "" {
				rt = "0"
			}

			schemaRetentionTime, err = strconv.ParseInt(rt, 10, 64)
			if err != nil {
				return err
			}
		}

		if database.RetentionTime != expectedDatabaseRetentionsDays {
			return fmt.Errorf("invalid database retention time, expected: %d, got: %d", expectedDatabaseRetentionsDays, database.RetentionTime)
		}

		if schemaRetentionTime != int64(expectedSchemaRetentionDays) {
			return fmt.Errorf("invalid schema retention time, expected: %d, got: %d", expectedSchemaRetentionDays, schemaRetentionTime)
		}

		return nil
	}
}

func setSchemaDataRetentionTime(t *testing.T, id sdk.DatabaseObjectIdentifier, days int) func() {
	t.Helper()

	return func() {
		client := acc.Client(t)
		ctx := context.Background()

		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Set: &sdk.SchemaSet{
				DataRetentionTimeInDays: sdk.Int(days),
			},
		})
		require.NoError(t, err)
	}
}

func testAccCheckSchemaDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_schema" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewDatabaseObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["name"])
		schema, err := client.Schemas.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("schema %v still exists", schema.Name)
		}
	}
	return nil
}

func removeSchemaOutsideOfTerraform(t *testing.T, databaseName string, schemaName string) {
	t.Helper()

	client := acc.Client(t)
	ctx := context.Background()

	err := client.Schemas.Drop(ctx, sdk.NewDatabaseObjectIdentifier(databaseName, schemaName), new(sdk.DropSchemaOptions))
	require.NoError(t, err)
}
