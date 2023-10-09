package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserOwnershipGrant_defaults(t *testing.T) {
	user := "tst-terraform" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	role := "tst-terraform" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
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
