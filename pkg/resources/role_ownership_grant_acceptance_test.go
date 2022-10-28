package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_RoleOwnershipGrant_defaults(t *testing.T) {
	onRoleName := "tst-terraform" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	toRoleName := "tst-terraform" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: roleOwnershipGrantConfig(onRoleName, toRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role_ownership_grant.grant", "on_role_name", onRoleName),
					resource.TestCheckResourceAttr("snowflake_role_ownership_grant.grant", "to_role_name", toRoleName),
					resource.TestCheckResourceAttr("snowflake_role_ownership_grant.grant", "current_grants", "COPY"),
				),
			},
		},
	})
}

func roleOwnershipGrantConfig(onRoleName, toRoleName string) string {
	return fmt.Sprintf(`

resource "snowflake_role" "role" {
	name = "%v"
}

resource "snowflake_role" "other_role" {
	name = "%v"
}

resource "snowflake_role_grants" "grants" {
	role_name = snowflake_role.role.name

	roles = [
		"ACCOUNTADMIN",
	]
}

resource "snowflake_role_ownership_grant" "grant" {
	on_role_name = snowflake_role.role.name

	to_role_name = snowflake_role.other_role.name

	current_grants = "COPY"
}
`, onRoleName, toRoleName)
}
