package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_DatabaseRoles(t *testing.T) {
	dbName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dbRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: databaseRoles(dbName, dbRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.name"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.comment"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.owner"),
				),
			},
			{
				Config: databaseRolesEmpty(dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.name"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.comment"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.owner"),
				),
			},
		},
	})
}

func databaseRoles(dbName, dbRoleName string) string {
	return fmt.Sprintf(`
		resource snowflake_database "test_db" {
			name = "%v"
		}

		resource snowflake_database_role "test_role" {
			name = "%v"
            comment = "test"
			database = snowflake_database.test_db.name
		}

		data snowflake_database_roles "db_roles" {
            database = snowflake_database.test_db.name
			depends_on = [
				snowflake_database_role.test_role,
			]
		}
	`, dbName, dbRoleName)
}

func databaseRolesEmpty(dbName string) string {
	return fmt.Sprintf(`
		resource snowflake_database "test_db" {
			name = "%v"
		}

		data snowflake_database_roles "db_roles" {
            database = snowflake_database.test_db.name
			depends_on = [
				snowflake_database.test_db,
			]
		}
	`, dbName)
}
