package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccViewGrant(t *testing.T) {
	vName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	roleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	shareName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfig(vName, roleName, shareName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", vName),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_view_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func viewGrantConfig(n, role, share string) string {
	return fmt.Sprintf(`

resource "snowflake_database" "test" {
  name = "%v"
}

resource "snowflake_view" "test" {
  name      = "%v"
  database  = snowflake_database.test.name
  statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
  is_secure = true
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_share" "test" {
  name = "%v"
}

resource "snowflake_database_grant" "test" {
  database_name = snowflake_view.test.database
  shares        = [snowflake_share.test.name]
}

resource "snowflake_view_grant" "test" {
  view_name     = snowflake_view.test.name
  database_name = snowflake_view.test.database
  roles         = [snowflake_role.test.name]
  shares        = [snowflake_share.test.name]

  depends_on = [snowflake_database_grant.test]
}
`, n, n, role, share)
}
