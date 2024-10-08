package resources_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
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
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	name := id.Name()
	secondId := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	nameLowerCase := strings.ToLower(secondId.Name())
	secondIdLowerCased := sdk.NewAccountObjectIdentifier(nameLowerCase)
	nameLowerCaseEscaped := fmt.Sprintf(`"%s"`, nameLowerCase)
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }

	resourceName := "snowflake_unsafe_execute.test"
	createConfigVariables := func(id string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(createDatabaseStatement(id)),
			"revert":  config.StringVariable(dropDatabaseStatement(id)),
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
				ConfigVariables: createConfigVariables(name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(name)),
					resource.TestCheckNoResourceAttr(resourceName, "query"),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, secondIdLowerCased, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(nameLowerCaseEscaped),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(nameLowerCaseEscaped)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(nameLowerCaseEscaped)),
					resource.TestCheckNoResourceAttr(resourceName, "query"),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, secondIdLowerCased, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_withRead(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	name := id.Name()
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }
	showDatabaseStatement := func(id string) string { return fmt.Sprintf("show databases like '%%%s%%'", id) }

	resourceName := "snowflake_unsafe_execute.test"
	createConfigVariables := func(id string) map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(createDatabaseStatement(id)),
			"revert":  config.StringVariable(dropDatabaseStatement(id)),
			"query":   config.StringVariable(showDatabaseStatement(id)),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_withRead"),
				ConfigVariables: createConfigVariables(name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(name)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.name", name),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.0.created_on"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.budget", ""),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.comment", ""),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_readRemoved(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	name := id.Name()
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }
	showDatabaseStatement := func(id string) string { return fmt.Sprintf("show databases like '%%%s%%'", id) }
	resourceName := "snowflake_unsafe_execute.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
					"query":   config.StringVariable(showDatabaseStatement(name)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(name)),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
					"query":   config.StringVariable(""),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", ""),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_badQuery(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	name := id.Name()
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }
	showDatabaseStatement := func(id string) string { return fmt.Sprintf("show databases like '%%%s%%'", id) }
	resourceName := "snowflake_unsafe_execute.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
					"query":   config.StringVariable("bad query"),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "query", "bad query"),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
					"query":   config.StringVariable(showDatabaseStatement(name)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(name)),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.name", name),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_invalidExecuteStatement(t *testing.T) {
	invalidCreateStatement := "create database"
	invalidDropStatement := "drop database"

	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"execute": config.StringVariable(invalidCreateStatement),
			"revert":  config.StringVariable(invalidDropStatement),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				ExpectError: regexp.MustCompile("SQL compilation error"),
			},
		},
	})
}

func TestAcc_UnsafeExecute_invalidRevertStatement(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	name := id.Name()
	updatedId := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	updatedName := updatedId.Name()
	createDatabaseStatement := func(id string) string { return fmt.Sprintf("create database %s", id) }
	dropDatabaseStatement := func(id string) string { return fmt.Sprintf("drop database %s", id) }
	invalidDropStatement := "drop database"

	resourceName := "snowflake_unsafe_execute.test"

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
			err = testAccCheckDatabaseExistence(t, updatedId, false)(state)
			if err != nil {
				return err
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(invalidDropStatement),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", invalidDropStatement),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(updatedName)),
					"revert":  config.StringVariable(invalidDropStatement),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				ExpectError: regexp.MustCompile("SQL compilation error"),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(name)),
					"revert":  config.StringVariable(dropDatabaseStatement(name)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(name)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(name)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
					testAccCheckDatabaseExistence(t, updatedId, false),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(updatedName)),
					"revert":  config.StringVariable(dropDatabaseStatement(updatedName)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(updatedName)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(updatedName)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, false),
					testAccCheckDatabaseExistence(t, updatedId, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_revertUpdated(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	name := id.Name()
	execute := fmt.Sprintf("create database %s", name)
	revert := fmt.Sprintf("drop database %s", name)
	notMatchingRevert := "select 1"
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
				ConfigVariables: createConfigVariables(execute, notMatchingRevert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", notMatchingRevert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
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
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	name := id.Name()
	execute := fmt.Sprintf("create database %s", name)
	revert := fmt.Sprintf("drop database %s", name)

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix("UNSAFE_EXECUTE_TEST_DATABASE_")
	newName := newId.Name()
	newExecute := fmt.Sprintf("create database %s", newName)
	newRevert := fmt.Sprintf("drop database %s", newName)

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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
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
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
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

func TestAcc_UnsafeExecute_grants(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	client := acc.TestClient()

	database, databaseCleanup := client.Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	role, roleCleanup := client.Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	privilege := sdk.AccountObjectPrivilegeCreateSchema
	execute := fmt.Sprintf("GRANT %s ON DATABASE %s TO ROLE %s", privilege, database.ID().FullyQualifiedName(), role.ID().FullyQualifiedName())
	revert := fmt.Sprintf("REVOKE %s ON DATABASE %s FROM ROLE %s", privilege, database.ID().FullyQualifiedName(), role.ID().FullyQualifiedName())

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
			err := verifyGrantExists(t, role.ID(), privilege, false)(state)
			return err
		},
		Steps: []resource.TestStep{
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
					verifyGrantExists(t, role.ID(), privilege, true),
				),
			},
		},
	})
}

// TestAcc_UnsafeExecute_grantsComplex test fails with:
//
//	testing_new_config.go:156: unexpected index type (string) for "snowflake_unsafe_execute.test[\"0\"]", for_each is not supported
//	testing_new.go:68: unexpected index type (string) for "snowflake_unsafe_execute.test[\"0\"]", for_each is not supported
//
// Quick search unveiled this issue: https://github.com/hashicorp/terraform-plugin-sdk/issues/536.
//
// It also seems that it is working correctly underneath; with TF_LOG set to DEBUG we have:
//
//	2023/11/26 17:16:03 [DEBUG] SQL "GRANT CREATE SCHEMA,MODIFY ON DATABASE UNSAFE_EXECUTE_TEST_DATABASE_4397 TO ROLE UNSAFE_EXECUTE_TEST_ROLE_1145" applied successfully
//	2023/11/26 17:16:03 [DEBUG] SQL "GRANT MODIFY,USAGE ON DATABASE UNSAFE_EXECUTE_TEST_DATABASE_3740 TO ROLE UNSAFE_EXECUTE_TEST_ROLE_3008" applied successfully
func TestAcc_UnsafeExecute_grantsComplex(t *testing.T) {
	t.Skip("Skipping TestAcc_UnsafeExecute_grantsComplex because of https://github.com/hashicorp/terraform-plugin-sdk/issues/536 issue")

	client := acc.TestClient()

	database1, database1Cleanup := client.Database.CreateDatabase(t)
	t.Cleanup(database1Cleanup)

	database2, database2Cleanup := client.Database.CreateDatabase(t)
	t.Cleanup(database2Cleanup)

	role1, role1Cleanup := client.Role.CreateRole(t)
	t.Cleanup(role1Cleanup)

	role2, role2Cleanup := client.Role.CreateRole(t)
	t.Cleanup(role2Cleanup)

	dbId1 := database1.ID()
	dbId2 := database2.ID()
	roleId1 := role1.ID()
	roleId2 := role2.ID()
	privilege1 := sdk.AccountObjectPrivilegeCreateSchema
	privilege2 := sdk.AccountObjectPrivilegeModify
	privilege3 := sdk.AccountObjectPrivilegeUsage

	// resourceName1 := "snowflake_unsafe_execute.test.0"
	// resourceName2 := "snowflake_unsafe_execute.test.1"
	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database_grants": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
				"database_name": config.StringVariable(dbId1.Name()),
				"role_id":       config.StringVariable(roleId1.Name()),
				"privileges":    config.ListVariable(config.StringVariable(privilege1.String()), config.StringVariable(privilege2.String())),
			}), config.ObjectVariable(map[string]config.Variable{
				"database_name": config.StringVariable(dbId2.Name()),
				"role_id":       config.StringVariable(roleId2.Name()),
				"privileges":    config.ListVariable(config.StringVariable(privilege2.String()), config.StringVariable(privilege3.String())),
			})),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: func(state *terraform.State) error {
			err := verifyGrantExists(t, roleId1, privilege1, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId1, privilege2, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId1, privilege3, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId2, privilege1, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId2, privilege2, false)(state)
			if err != nil {
				return err
			}
			err = verifyGrantExists(t, roleId2, privilege3, false)(state)
			if err != nil {
				return err
			}
			return err
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: createConfigVariables(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet(resourceName1, "id"),
					// resource.TestCheckResourceAttrSet(resourceName2, "id"),
					verifyGrantExists(t, roleId1, privilege1, true),
					verifyGrantExists(t, roleId1, privilege2, true),
					verifyGrantExists(t, roleId1, privilege3, false),
					verifyGrantExists(t, roleId2, privilege1, false),
					verifyGrantExists(t, roleId2, privilege2, true),
					verifyGrantExists(t, roleId2, privilege3, true),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2491
func TestAcc_UnsafeExecute_queryResultsBug(t *testing.T) {
	resourceName := "snowflake_unsafe_execute.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: unsafeExecuteConfig(108),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", "SELECT 108"),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.0.108"),
				),
			},
			{
				Config: unsafeExecuteConfig(96),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", "SELECT 96"),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.0.96"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func unsafeExecuteConfig(queryNumber int) string {
	return fmt.Sprintf(`
resource "snowflake_unsafe_execute" "test" {
  execute = "SELECT 18"
  revert  = "SELECT 36"
  query  = "SELECT %d"
}

output "unsafe" {
  value = snowflake_unsafe_execute.test.query_results
}
`, queryNumber)
}

func verifyGrantExists(t *testing.T, roleId sdk.AccountObjectIdentifier, privilege sdk.AccountObjectPrivilege, shouldExist bool) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		grants, err := acc.TestClient().Grant.ShowGrantsToAccountRole(t, roleId)
		if err != nil {
			return err
		}

		if shouldExist {
			require.Equal(t, 1, len(grants))
			assert.Equal(t, privilege.String(), grants[0].Privilege)
			assert.Equal(t, sdk.ObjectTypeDatabase, grants[0].GrantedOn)
			assert.Equal(t, sdk.ObjectTypeRole, grants[0].GrantedTo)
			assert.Equal(t, roleId.FullyQualifiedName(), grants[0].GranteeName.FullyQualifiedName())
		} else {
			require.Equal(t, 0, len(grants))
		}

		// it does not matter what we return, because we have assertions above
		return nil
	}
}
