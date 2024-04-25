package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RoleOwnershipGrant_defaults(t *testing.T) {
	onRoleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	toRoleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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
