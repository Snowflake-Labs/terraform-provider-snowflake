package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccViewGrantBasic(t *testing.T) {
	vName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	roleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfigFuture(vName, roleName, false),
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

func TestAccViewGrantShares(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SHARE_TESTS"); ok {
		t.Skip("Skipping TestAccViewGrantShares")
	}

	vName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	roleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	shareName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfigShares(vName, roleName, shareName),
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

func TestAccFutureViewGrantChange(t *testing.T) {
	vName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	roleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfigFuture(vName, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", vName),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
			// CHANGE FROM CURRENT TO FUTURE VIEWS
			{
				Config: viewGrantConfigFuture(vName, roleName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", ""),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "on_future", "true"),
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

func viewGrantConfigShares(n, role, share string) string {
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
  name     = "%v"
  accounts = ["PC37737"]
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

func viewGrantConfigFuture(n string, role string, future bool) string {
	view_name_config := "view_name = snowflake_view.test.name"
	if future {
		view_name_config = "on_future = true"
	}
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

resource "snowflake_view_grant" "test" {
  %v
  database_name = snowflake_view.test.database
  roles         = [snowflake_role.test.name]
}
`, n, n, role, view_name_config)
}
