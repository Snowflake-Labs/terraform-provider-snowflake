package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRoles(t *testing.T) {
	dbName := acc.TestClient().Ids.Alpha()
	dbRoleName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
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
