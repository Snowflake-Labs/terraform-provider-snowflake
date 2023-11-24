package resources_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAcc_UnsafeMigration_basic(t *testing.T) {
	id := "unsafe_migration_test_database"
	up := fmt.Sprintf("create database %s", id)
	down := fmt.Sprintf("drop database %s", id)

	resourceName := "snowflake_unsafe_migration.migration"
	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"up":   config.StringVariable(up),
			"down": config.StringVariable(down),
		}
	}

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
				ConfigVariables: createConfigVariables(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "up", up),
					resource.TestCheckResourceAttr(resourceName, "down", down),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						return nil
					}),
					// TODO: check if exists after apply
					// testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeMigration_downChanged(t *testing.T) {
	id := "unsafe_migration_test_database"
	up := fmt.Sprintf("create database %s", id)
	down := fmt.Sprintf("drop database %s", id)
	invalidDown := "select 1"
	var savedId string

	resourceName := "snowflake_unsafe_migration.migration"
	createConfigVariables := func(up string, down string) map[string]config.Variable {
		return map[string]config.Variable{
			"up":   config.StringVariable(up),
			"down": config.StringVariable(down),
		}
	}

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
				ConfigVariables: createConfigVariables(up, invalidDown),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "up", up),
					resource.TestCheckResourceAttr(resourceName, "down", invalidDown),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						savedId = value
						return nil
					}),
					// TODO: check if exists after apply
					// testAccCheckDatabaseExistence(id, true),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: createConfigVariables(up, down),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "up", up),
					resource.TestCheckResourceAttr(resourceName, "down", down),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						if savedId != value {
							return errors.New("different id after down update")
						}
						return nil
					}),
					// TODO: check if exists after down update
					// testAccCheckDatabaseExistence(id, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeMigration_upChanged(t *testing.T) {
	id := "unsafe_migration_test_database"
	up := fmt.Sprintf("create database %s", id)
	down := fmt.Sprintf("drop database %s", id)

	newId := "unsafe_migration_test_database_2"
	newUp := fmt.Sprintf("create database %s", newId)
	newDown := fmt.Sprintf("drop database %s", newId)

	var savedId string

	resourceName := "snowflake_unsafe_migration.migration"
	createConfigVariables := func(up string, down string) map[string]config.Variable {
		return map[string]config.Variable{
			"up":   config.StringVariable(up),
			"down": config.StringVariable(down),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: func(state *terraform.State) error {
			err := testAccCheckDatabaseExistence(t, id, false)(state)
			if err != nil {
				return err
			}
			err = testAccCheckDatabaseExistence(t, newId, false)(state)
			if err != nil {
				return err
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: createConfigVariables(up, down),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "up", up),
					resource.TestCheckResourceAttr(resourceName, "down", down),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						savedId = value
						return nil
					}),
					// TODO: check if exists after apply
					// testAccCheckDatabaseExistence(id, true),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: createConfigVariables(up, down),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "up", newUp),
					resource.TestCheckResourceAttr(resourceName, "down", newDown),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						if savedId == value {
							return errors.New("same id after up update")
						}
						return nil
					}),
					// TODO: check if exists after up update and check that old database doesn't exist (may be duplicate with check destroy)
					// testAccCheckDatabaseExistence(id, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeMigration_grants(t *testing.T) {
	id := "unsafe_migration_test_database"
	up := fmt.Sprintf("create database %s", id)
	down := fmt.Sprintf("drop database %s", id)
	// TODO: before test
	// create role
	// create database

	// create migration
	// - up: grant ... to role xyz
	// - down: revoke ... from role xyz

	resourceName := "snowflake_unsafe_migration.migration"
	createConfigVariables := func(up string, down string) map[string]config.Variable {
		return map[string]config.Variable{
			"up":   config.StringVariable(up),
			"down": config.StringVariable(down),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: func(state *terraform.State) error {
			return dropResourcesForMigrationTestCaseForGrants(t)
		},
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createResourcesForMigrationTestCaseForGrants(t) },
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: createConfigVariables(up, down),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "up", up),
					resource.TestCheckResourceAttr(resourceName, "down", down),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						return nil
					}),
					// TODO: check if exists after apply
					// testAccCheckDatabaseExistence(id, true),
				),
			},
		},
	})
}

func testAccCheckDatabaseExistence(t *testing.T, id string, shouldExist bool) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		db := acc.TestAccProvider.Meta().(*sql.DB)
		client := sdk.NewClientFromDB(db)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "snowflake_unsafe_migration" {
				continue
			}
			ctx := context.Background()
			_, err := client.Databases.ShowByID(ctx, sdk.NewAccountObjectIdentifier(id))
			if shouldExist {
				if err != nil {
					return fmt.Errorf("error while retrieving database %s, err = %w", id, err)
				}
			} else {
				if err == nil {
					return fmt.Errorf("database %v still exists", id)
				}
			}
		}
		return nil
	}
}

func createResourcesForMigrationTestCaseForGrants(t *testing.T, databaseName string, roleName string) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()

	err = client.Databases.Create(ctx, sdk.NewAccountObjectIdentifier("unsafe_migration_test_database"), &sdk.CreateDatabaseOptions{})
	require.NoError(t, err)

	err = client.Roles.Create(ctx, sdk.NewCreateRoleRequest(sdk.NewAccountObjectIdentifier("unsafe_migration_test_role")))
	require.NoError(t, err)
}

func dropResourcesForMigrationTestCaseForGrants(t *testing.T, databaseName string, roleName string) error {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()

	err = client.Databases.Drop(ctx, sdk.NewAccountObjectIdentifier(databaseName), &sdk.DropDatabaseOptions{})
	require.NoError(t, err)

	err = client.Roles.Drop(ctx, sdk.NewDropRoleRequest(sdk.NewAccountObjectIdentifier(roleName)))
	require.NoError(t, err)
}
