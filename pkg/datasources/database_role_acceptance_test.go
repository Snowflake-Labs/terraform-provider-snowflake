package datasources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRole(t *testing.T) {
	dbName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dbRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: databaseRole(dbName, dbRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_database_role.db_role", "name"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_role.db_role", "comment"),
					resource.TestCheckResourceAttrSet("data.snowflake_database_role.db_role", "owner"),
				),
			},
			{
				Config:      databaseRoleEmpty(dbName),
				ExpectError: regexp.MustCompile("Error: object does not exist"),
			},
		},
	})
}

func databaseRole(dbName, dbRoleName string) string {
	return fmt.Sprintf(`
		resource snowflake_database "test_db" {
			name = "%v"
		}

		resource snowflake_database_role "test_role" {
			name = "%v"
            comment = "test"
			database = snowflake_database.test_db.name
		}

		data snowflake_database_role "db_role" {
            database = snowflake_database.test_db.name
			name = snowflake_database_role.test_role.name
			depends_on = [
				snowflake_database_role.test_role,
			]
		}
	`, dbName, dbRoleName)
}

func databaseRoleEmpty(dbName string) string {
	return fmt.Sprintf(`
		resource snowflake_database "test_db" {
			name = "%v"
		}

		data snowflake_database_role "db_role" {
            database = snowflake_database.test_db.name
			name = "dummy_missing"
			depends_on = [
				snowflake_database.test_db,
			]
		}
	`, dbName)
}
