package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRoles(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	databaseRoleNamePrefix := acc.TestClient().Ids.Alpha()
	databaseRoleName := databaseRoleNamePrefix + "1" + acc.TestClient().Ids.Alpha()
	databaseRoleName2 := databaseRoleNamePrefix + "2" + acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: databaseRoles(databaseName, databaseRoleName, databaseRoleName2, databaseRoleNamePrefix+"%"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database_roles.db_roles", "database_roles.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.name", databaseRoleName2),
					resource.TestCheckResourceAttr("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.is_current"),
					resource.TestCheckResourceAttr("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.granted_to_database_roles", "0"),
					resource.TestCheckResourceAttr("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.granted_database_roles", "0"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.comment", "test"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_roles.db_roles", "database_roles.0.show_output.0.owner_role_type"),
				),
			},
		},
	})
}

func databaseRoles(databaseName, databaseRoleName, databaseRoleName2, databaseNamePrefix string) string {
	return fmt.Sprintf(`
resource snowflake_database "test_db" {
	name = "%v"
}

resource snowflake_database_role "test_role" {
	name = "%v"
	comment = "test"
	database = snowflake_database.test_db.name
}

resource snowflake_database_role "test_role2" {
	name = "%v"
	comment = "test"
	database = snowflake_database.test_db.name
}

data snowflake_database_roles "db_roles" {
	depends_on = [ snowflake_database_role.test_role, snowflake_database_role.test_role2 ]

	in_database = snowflake_database.test_db.name
	like = "%v"
	limit {
		rows = 1
		from = snowflake_database_role.test_role.name
	}
}
	`, databaseName, databaseRoleName, databaseRoleName2, databaseNamePrefix)
}
