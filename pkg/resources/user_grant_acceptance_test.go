package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_UserGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_USER_GRANT_TESTS"); ok {
		t.Skip("Skipping TestAccUserGrant")
	}
	wName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: userGrantConfig(wName, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_grant.test", "user_name", wName),
					resource.TestCheckResourceAttr("snowflake_user_grant.test", "privilege", "MONITOR"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_user_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func userGrantConfig(n, role string) string {
	return fmt.Sprintf(`

resource "snowflake_user" "test" {
  name      = "%v"
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_user_grant" "test" {
  user_name = snowflake_user.test.name
  roles     = [snowflake_role.test.name]
  privilege = "MONITOR"
}
`, n, role)
}
