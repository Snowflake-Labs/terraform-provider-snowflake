package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserOwnershipGrant_defaults(t *testing.T) {
	user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	role := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: userOwnershipGrantConfig(user, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_ownership_grant.grant", "on_user_name", user),
					resource.TestCheckResourceAttr("snowflake_user_ownership_grant.grant", "to_role_name", role),
					resource.TestCheckResourceAttr("snowflake_user_ownership_grant.grant", "current_grants", "COPY"),
				),
			},
		},
	})
}

func userOwnershipGrantConfig(user, role string) string {
	return fmt.Sprintf(`

resource "snowflake_user" "user" {
	name = "%v"
}

resource "snowflake_role" "role" {
	name = "%v"
}

resource "snowflake_role_grants" "grants" {
	role_name = snowflake_role.role.name

	roles = [
		"ACCOUNTADMIN",
	]
}

resource "snowflake_user_ownership_grant" "grant" {
	on_user_name = snowflake_user.user.name

	to_role_name = snowflake_role.role.name

	current_grants = "COPY"
}
`, user, role)
}
