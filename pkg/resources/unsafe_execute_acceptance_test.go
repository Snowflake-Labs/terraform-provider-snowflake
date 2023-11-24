package resources_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcc_UnsafeExecute_basic(t *testing.T) {
	// TODO: capitalized
	id := fmt.Sprintf("UNSAFE_EXECUTE_TEST_DATABASE_%d", rand.Intn(10000))
	execute := fmt.Sprintf("create database %s", id)
	revert := fmt.Sprintf("drop database %s", id)

	resourceName := "snowflake_unsafe_execute.test"
	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(execute),
			"revert":  config.StringVariable(revert),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_revertUpdated(t *testing.T) {
	// TODO: capitalized
	id := fmt.Sprintf("UNSAFE_EXECUTE_TEST_DATABASE_%d", rand.Intn(10000))
	execute := fmt.Sprintf("create database %s", id)
	revert := fmt.Sprintf("drop database %s", id)
	// TODO: this is not invalid but it does not match the execute
	invalidRevert := "select 1"
	var savedId string

	resourceName := "snowflake_unsafe_execute.test"
	createConfigVariables := func(execute string, revert string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(execute),
			"revert":  config.StringVariable(revert),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, invalidRevert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", invalidRevert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						savedId = value
						return nil
					}),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, revert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						if savedId != value {
							return errors.New("different id after revert update")
						}
						return nil
					}),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_executeUpdated(t *testing.T) {
	// TODO: capitalized
	id := fmt.Sprintf("UNSAFE_EXECUTE_TEST_DATABASE_%d", rand.Intn(10000))
	execute := fmt.Sprintf("create database %s", id)
	revert := fmt.Sprintf("drop database %s", id)

	// TODO: capitalized
	newId := fmt.Sprintf("%s_2", id)
	newExecute := fmt.Sprintf("create database %s", newId)
	newRevert := fmt.Sprintf("drop database %s", newId)

	var savedId string

	resourceName := "snowflake_unsafe_execute.test"
	createConfigVariables := func(execute string, revert string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(execute),
			"revert":  config.StringVariable(revert),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, revert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						savedId = value
						return nil
					}),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(newExecute, newRevert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", newExecute),
					resource.TestCheckResourceAttr(resourceName, "revert", newRevert),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == "" {
							return errors.New("empty id")
						}
						if savedId == value {
							return errors.New("same id after execute update")
						}
						return nil
					}),
					testAccCheckDatabaseExistence(t, id, false),
					testAccCheckDatabaseExistence(t, newId, true),
				),
			},
		},
	})
}

// TODO: make this test pass
func TestAcc_UnsafeExecute_grants(t *testing.T) {
	id := "UNSAFE_EXECUTE_test_database"
	execute := fmt.Sprintf("create database %s", id)
	revert := fmt.Sprintf("drop database %s", id)
	// TODO: before test
	// create role
	// create database

	// create migration
	// - execute: grant ... to role xyz
	// - revert: revoke ... from role xyz

	resourceName := "snowflake_unsafe_execute.test"
	createConfigVariables := func(execute string, revert string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(execute),
			"revert":  config.StringVariable(revert),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: func(state *terraform.State) error {
			return dropResourcesForUnsafeExecuteTestCaseForGrants(t)
		},
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createResourcesForExecuteUnsafeTestCaseForGrants(t) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, revert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
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
		client, err := sdk.NewDefaultClient()
		if err != nil {
			t.Fatal(err)
		}
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

// TODO: tweak this method
func createResourcesForExecuteUnsafeTestCaseForGrants(t *testing.T) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()

	err = client.Databases.Create(ctx, sdk.NewAccountObjectIdentifier("UNSAFE_EXECUTE_test_database"), &sdk.CreateDatabaseOptions{})
	require.NoError(t, err)

	err = client.Roles.Create(ctx, sdk.NewCreateRoleRequest(sdk.NewAccountObjectIdentifier("UNSAFE_EXECUTE_test_role")))
	require.NoError(t, err)
}

// TODO: fix this method
func dropResourcesForUnsafeExecuteTestCaseForGrants(t *testing.T) error {
	t.Helper()

	databaseName := "TODO"
	roleName := "TODO"

	client, err := sdk.NewDefaultClient()
	if err != nil {
		return err
	}
	ctx := context.Background()

	err = client.Databases.Drop(ctx, sdk.NewAccountObjectIdentifier(databaseName), &sdk.DropDatabaseOptions{})
	assert.NoError(t, err)
	if err != nil {
		return err
	}

	err = client.Roles.Drop(ctx, sdk.NewDropRoleRequest(sdk.NewAccountObjectIdentifier(roleName)))
	assert.NoError(t, err)
	if err != nil {
		return err
	}

	return nil
}
