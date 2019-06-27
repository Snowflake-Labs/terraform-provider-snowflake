package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDatabaseGrant(t *testing.T) {
	dbName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	roleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	shareName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: databaseGrantConfig(dbName, roleName, shareName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_grant.test", "database_name", dbName),
					resource.TestCheckResourceAttr("snowflake_database_grant.test", "privilege", "USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_database_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func databaseGrantConfig(db, role, share string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
}
resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_share" "test" {
  name = "%v"
}

resource "snowflake_database_grant" "test" {
  database_name = snowflake_database.test.name
  roles         = [snowflake_role.test.name]
  shares        = [snowflake_share.test.name]
}
`, db, role, share)
}
