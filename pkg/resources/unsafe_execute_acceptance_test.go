package resources_test

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jmoiron/sqlx"
	"log"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	accTestDb  *sql.DB
	accTestDbx *sqlx.DB
)

func init() {
	db, err := provider.GetDatabaseHandleFromEnv()
	if err != nil {
		log.Fatalln(err)
	}
	dbx := sqlx.NewDb(db, "snowflake")
	accTestDb = db
	accTestDbx = dbx
}

func TestAcc_UnsafeExecute_basic(t *testing.T) {
	id := generateUnsafeExecuteTestDatabaseName(t)
	idLowerCase := strings.ToLower(generateUnsafeExecuteTestDatabaseName(t))
	createDatabaseStatement := func(raw bool, id string) string {
		if raw {
			return fmt.Sprintf(`create database \"%s\"`, id)
		}
		return fmt.Sprintf("create database \"%s\"", id)
	}
	dropDatabaseStatement := func(raw bool, id string) string {
		if raw {
			return fmt.Sprintf(`drop database \"%s\"`, id)
		}
		return fmt.Sprintf("drop database \"%s\"", id)
	}
	resourceName := "snowflake_unsafe_execute.test"

	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				Config: schemaExecResourceConfig(createDatabaseStatement(true, id), dropDatabaseStatement(true, id)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(false, id)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(false, id)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, id, true),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: testAccCheckDatabaseExistence(t, idLowerCase, false),
		Steps: []resource.TestStep{
			{
				Config: schemaExecResourceConfig(createDatabaseStatement(true, idLowerCase), dropDatabaseStatement(true, idLowerCase)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "execute", createDatabaseStatement(false, idLowerCase)),
					resource.TestCheckResourceAttr(resourceName, "revert", dropDatabaseStatement(false, idLowerCase)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testAccCheckDatabaseExistence(t, idLowerCase, true),
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

	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: testAccCheckDatabaseExistence(t, id, false),
		Steps: []resource.TestStep{
			{
				Config: schemaExecResourceConfig(execute, notMatchingRevert),
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
				Config: schemaExecResourceConfig(execute, revert),
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

	resource.Test(t, resource.TestCase{
		Providers: providers(),
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
				Config: schemaExecResourceConfig(execute, revert),
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
				Config: schemaExecResourceConfig(newExecute, newRevert),
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
	privilege := "CREATE SCHEMA"
	execute := fmt.Sprintf("GRANT %s ON DATABASE %s TO ROLE %s", privilege, id, roleId)
	revert := fmt.Sprintf("REVOKE %s ON DATABASE %s FROM ROLE %s", privilege, id, roleId)
	resourceName := "snowflake_unsafe_execute.test"

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		CheckDestroy: func(state *terraform.State) error {
			err := verifyGrantExists(t, roleId, privilege, false)(state)
			dropResourcesForUnsafeExecuteTestCaseForGrants(t, id, roleId)
			return err
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() { createResourcesForExecuteUnsafeTestCaseForGrants(t, id, roleId) },
				Config:    schemaExecResourceConfig(execute, revert),
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

func schemaExecResourceConfig(exec string, revert string) string {
	return fmt.Sprintf(`
resource "snowflake_unsafe_execute" "test" {
  execute = "%s"
  revert = "%s"
}
`, exec, revert)
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
		_, err := snowflake.ListDatabase(accTestDbx, id)

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

	createDatabaseSQL := snowflake.NewDatabaseBuilder(dbId).Create()
	err := snowflake.Exec(accTestDb, createDatabaseSQL)
	require.NoError(t, err)

	createRoleSQL := snowflake.NewRoleBuilder(roleId).Create().Statement()
	err = snowflake.Exec(accTestDb, createRoleSQL)
	require.NoError(t, err)
}

func dropResourcesForUnsafeExecuteTestCaseForGrants(t *testing.T, dbId string, roleId string) {
	t.Helper()

	dropDatabaseSQL := snowflake.NewDatabaseBuilder(dbId).Drop()
	err := snowflake.Exec(accTestDb, dropDatabaseSQL)
	require.NoError(t, err)

	dropRoleSQL := snowflake.NewRoleBuilder(roleId).Drop()
	err = snowflake.Exec(accTestDb, dropRoleSQL)
	require.NoError(t, err)
}

func verifyGrantExists(t *testing.T, roleId string, privilege string, shouldExist bool) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		grants, err := snowflake.ShowGrantsTo(accTestDb, "ROLE", roleId)
		require.NoError(t, err)

		if shouldExist {
			require.Equal(t, 1, len(grants))

			assert.True(t, grants[0].Privilege.Valid)
			assert.Equal(t, privilege, grants[0].Privilege.String)

			assert.True(t, grants[0].GrantedOn.Valid)
			assert.Equal(t, "DATABASE", grants[0].GrantedOn.String)

			assert.True(t, grants[0].GrantedTo.Valid)
			assert.Equal(t, "ROLE", grants[0].GrantedTo.String)

			assert.True(t, grants[0].GranteeName.Valid)
			assert.Equal(t, roleId, grants[0].GranteeName.String)
		} else {
			require.Equal(t, 0, len(grants))
		}

		// it does not matter what we return, because we have assertions above
		return nil
	}
}
