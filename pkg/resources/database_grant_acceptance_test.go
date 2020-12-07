package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testRolesAndShares(t *testing.T, path string, roles, shares []string) func(*terraform.State) error {
	return func(state *terraform.State) error {
		is := state.RootModule().Resources[path].Primary

		if c, ok := is.Attributes["roles.#"]; !ok || MustParseInt(t, c) != int64(len(roles)) {
			return fmt.Errorf("expected roles.# to equal %d but got %s", len(roles), c)
		}
		r, err := extractList(is.Attributes, "roles")
		if err != nil {
			return err
		}

		// TODO case sensitive?
		if !listSetEqual(roles, r) {
			return fmt.Errorf("expected roles %#v but got %#v", roles, r)
		}

		return nil
	}
}

func TestAcc_DatabaseGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SHARE_TESTS"); ok {
		t.Skip("Skipping TestAccDatabaseGrant")
	}

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
					resource.TestCheckResourceAttr("snowflake_database_grant.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database_grant.test", "shares.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database_grant.test", "shares.#", "1"),
					testRolesAndShares(t, "snowflake_database_grant.test", []string{roleName}, []string{shareName}),
				),
			},
			// IMPORT
			{
				PreConfig:         func() { fmt.Println("[DEBUG] IMPORT") },
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
