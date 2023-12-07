package resources_test

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
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
	id := generateUnsafeExecuteTestDatabaseName(t)
	idLowerCase := strings.ToLower(generateUnsafeExecuteTestDatabaseName(t))
	idLowerCaseEscaped := fmt.Sprintf(`"%s"`, idLowerCase)
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
				ConfigVariables: createConfigVariables(id),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(id)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(id)),
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
		CheckDestroy: testAccCheckDatabaseExistence(t, idLowerCase, false),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(idLowerCaseEscaped),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(idLowerCaseEscaped)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(idLowerCaseEscaped)),
					resource.TestCheckNoResourceAttr(resourceName, "query"),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, idLowerCase, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_withRead(t *testing.T) {
	id := generateUnsafeExecuteTestDatabaseName(t)
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
				ConfigVariables: createConfigVariables(id),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(id)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(id)),
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(id)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.name", id),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.0.created_on"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.budget", ""),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.comment", ""),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_readRemoved(t *testing.T) {
	id := generateUnsafeExecuteTestDatabaseName(t)
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
					"execute": config.StringVariable(createDatabaseStatement(id)),
					"revert":  config.StringVariable(dropDatabaseStatement(id)),
					"query":   config.StringVariable(showDatabaseStatement(id)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(id)),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(id)),
					"revert":  config.StringVariable(dropDatabaseStatement(id)),
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
	id := generateUnsafeExecuteTestDatabaseName(t)
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
					"execute": config.StringVariable(createDatabaseStatement(id)),
					"revert":  config.StringVariable(dropDatabaseStatement(id)),
					"query":   config.StringVariable("bad query"),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(id)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(id)),
					resource.TestCheckResourceAttr(resourceName, "query", "bad query"),
					resource.TestCheckNoResourceAttr(resourceName, "query_results.#"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_withRead"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(id)),
					"revert":  config.StringVariable(dropDatabaseStatement(id)),
					"query":   config.StringVariable(showDatabaseStatement(id)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query", showDatabaseStatement(id)),
					resource.TestCheckResourceAttrSet(resourceName, "query_results.#"),
					resource.TestCheckResourceAttr(resourceName, "query_results.0.name", id),
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
	id := generateUnsafeExecuteTestDatabaseName(t)
	updatedId := generateUnsafeExecuteTestDatabaseName(t)
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
					"execute": config.StringVariable(createDatabaseStatement(id)),
					"revert":  config.StringVariable(invalidDropStatement),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(id)),
					resource.TestCheckResourceAttr(resourceName, "revert", invalidDropStatement),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(updatedId)),
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
					"execute": config.StringVariable(createDatabaseStatement(id)),
					"revert":  config.StringVariable(dropDatabaseStatement(id)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(id)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(id)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
					testAccCheckDatabaseExistence(t, updatedId, false),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: map[string]config.Variable{
					"execute": config.StringVariable(createDatabaseStatement(updatedId)),
					"revert":  config.StringVariable(dropDatabaseStatement(updatedId)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(updatedId)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(updatedId)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, false),
					testAccCheckDatabaseExistence(t, updatedId, true),
				),
			},
		},
	})
}

func TestAcc_UnsafeExecute_revertUpdated(t *testing.T) {
	id := generateUnsafeExecuteTestDatabaseName(t)
	execute := fmt.Sprintf("create database %s", id)
	revert := fmt.Sprintf("drop database %s", id)
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
	id := generateUnsafeExecuteTestDatabaseName(t)
	execute := fmt.Sprintf("create database %s", id)
	revert := fmt.Sprintf("drop database %s", id)

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
	id := generateUnsafeExecuteTestDatabaseName(t)
	roleId := generateUnsafeExecuteTestRoleName(t)
	privilege := sdk.AccountObjectPrivilegeCreateSchema
	execute := fmt.Sprintf("GRANT %s ON DATABASE %s TO ROLE %s", privilege, id, roleId)
	revert := fmt.Sprintf("REVOKE %s ON DATABASE %s FROM ROLE %s", privilege, id, roleId)

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
			err := verifyGrantExists(t, roleId, privilege, false)(state)
			dropResourcesForUnsafeExecuteTestCaseForGrants(t, id, roleId)
			return err
		},
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createResourcesForExecuteUnsafeTestCaseForGrants(t, id, roleId) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_UnsafeExecute_commonSetup"),
				ConfigVariables: createConfigVariables(execute, revert),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", execute),
					resource.TestCheckResourceAttr(resourceName, "revert", revert),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					verifyGrantExists(t, roleId, privilege, true),
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

	dbId1 := generateUnsafeExecuteTestDatabaseName(t)
	dbId2 := generateUnsafeExecuteTestDatabaseName(t)
	roleId1 := generateUnsafeExecuteTestRoleName(t)
	roleId2 := generateUnsafeExecuteTestRoleName(t)
	privilege1 := sdk.AccountObjectPrivilegeCreateSchema
	privilege2 := sdk.AccountObjectPrivilegeModify
	privilege3 := sdk.AccountObjectPrivilegeUsage

	// resourceName1 := "snowflake_unsafe_execute.test.0"
	// resourceName2 := "snowflake_unsafe_execute.test.1"
	createConfigVariables := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database_grants": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
				"database_name": config.StringVariable(dbId1),
				"role_id":       config.StringVariable(roleId1),
				"privileges":    config.ListVariable(config.StringVariable(privilege1.String()), config.StringVariable(privilege2.String())),
			}), config.ObjectVariable(map[string]config.Variable{
				"database_name": config.StringVariable(dbId2),
				"role_id":       config.StringVariable(roleId2),
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
			dropResourcesForUnsafeExecuteTestCaseForGrants(t, dbId1, roleId1)
			dropResourcesForUnsafeExecuteTestCaseForGrants(t, dbId2, roleId2)
			return err
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					createResourcesForExecuteUnsafeTestCaseForGrants(t, dbId1, roleId1)
					createResourcesForExecuteUnsafeTestCaseForGrants(t, dbId2, roleId2)
				},
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

// generateUnsafeExecuteTestDatabaseName returns capitalized name on purpose.
// Using small caps without escaping creates problem with later using sdk client which uses identifier that is escaped by default.
func generateUnsafeExecuteTestDatabaseName(t *testing.T) string {
	t.Helper()
	id, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		t.Fatalf("Failed to generate database id: %v", err)
	}
	return fmt.Sprintf("UNSAFE_EXECUTE_TEST_DATABASE_%d", id)
}

// generateUnsafeExecuteTestRoleName returns capitalized name on purpose.
// Using small caps without escaping creates problem with later using sdk client which uses identifier that is escaped by default.
func generateUnsafeExecuteTestRoleName(t *testing.T) string {
	t.Helper()
	id, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		t.Fatalf("Failed to generate role id: %v", err)
	}
	return fmt.Sprintf("UNSAFE_EXECUTE_TEST_ROLE_%d", id)
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

func createResourcesForExecuteUnsafeTestCaseForGrants(t *testing.T, dbId string, roleId string) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()

	err = client.Databases.Create(ctx, sdk.NewAccountObjectIdentifier(dbId), &sdk.CreateDatabaseOptions{})
	require.NoError(t, err)

	err = client.Roles.Create(ctx, sdk.NewCreateRoleRequest(sdk.NewAccountObjectIdentifier(roleId)))
	require.NoError(t, err)
}

func dropResourcesForUnsafeExecuteTestCaseForGrants(t *testing.T, dbId string, roleId string) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()

	err = client.Databases.Drop(ctx, sdk.NewAccountObjectIdentifier(dbId), &sdk.DropDatabaseOptions{})
	assert.NoError(t, err)

	err = client.Roles.Drop(ctx, sdk.NewDropRoleRequest(sdk.NewAccountObjectIdentifier(roleId)))
	assert.NoError(t, err)
}

func verifyGrantExists(t *testing.T, roleId string, privilege sdk.AccountObjectPrivilege, shouldExist bool) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		client, err := sdk.NewDefaultClient()
		require.NoError(t, err)
		ctx := context.Background()

		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: sdk.NewAccountObjectIdentifier(roleId),
			},
		})
		require.NoError(t, err)

		if shouldExist {
			require.Equal(t, 1, len(grants))
			assert.Equal(t, privilege.String(), grants[0].Privilege)
			assert.Equal(t, sdk.ObjectTypeDatabase, grants[0].GrantedOn)
			assert.Equal(t, sdk.ObjectTypeRole, grants[0].GrantedTo)
			assert.Equal(t, sdk.NewAccountObjectIdentifier(roleId).FullyQualifiedName(), grants[0].GranteeName.FullyQualifiedName())
		} else {
			require.Equal(t, 0, len(grants))
		}

		// it does not matter what we return, because we have assertions above
		return nil
	}
}
