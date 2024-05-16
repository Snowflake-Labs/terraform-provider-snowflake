package resources_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Schema(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	comment := "Terraform acceptance test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
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
	oldSchemaName := acc.TestClient().Ids.Alpha()
	newSchemaName := acc.TestClient().Ids.Alpha()
	comment := "Terraform acceptance test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
					},
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
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
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
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
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
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "5"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 15),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "15"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 15),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "0"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 0),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "3"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 3),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Schema_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
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
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 5, 5),
				),
			},
			{
				PreConfig:       acc.TestClient().Schema.UpdateDataRetentionTime(t, id, 20),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "-1"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_days", "3"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 3),
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
	schemaName := acc.TestClient().Ids.Alpha()
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
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_RemoveOutsideOfTerraform"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					acc.TestClient().Schema.DropSchemaFunc(t, sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, schemaName))()
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				RefreshPlanChecks: resource.RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAcc_Schema_RemoveSchemaOutsideOfTerraform(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	configVariables := map[string]config.Variable{
		"schema_name":   config.StringVariable(schemaName),
		"database_name": config.StringVariable(databaseName),
	}

	var cleanupDatabase func()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanupDatabase = acc.TestClient().Database.CreateDatabaseWithName(t, databaseName)
					t.Cleanup(cleanupDatabase)
				},
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

func checkDatabaseAndSchemaDataRetentionTime(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, expectedDatabaseRetentionsDays int, expectedSchemaRetentionDays int) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		schema, err := acc.TestClient().Schema.Show(t, schemaId)
		if err != nil {
			return err
		}

		database, err := acc.TestClient().Database.Show(t, schemaId.DatabaseId())
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
