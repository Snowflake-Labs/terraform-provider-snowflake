package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FailoverGroupGrant(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: failoverGroupGrantConfig(name, accountName, "FAILOVER"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group_grant.g", "failover_group_name", name),
					resource.TestCheckResourceAttr("snowflake_failover_group_grant.g", "privilege", "FAILOVER"),
					resource.TestCheckResourceAttr("snowflake_failover_group_grant.g", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_failover_group_grant.g", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group_grant.g", "roles.0", name),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_failover_group_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func failoverGroupGrantConfig(name string, allowedAccount string, privilege string) string {
	return fmt.Sprintf(`

resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types =  ["ROLES"]
	allowed_accounts= ["%s"]
}

resource "snowflake_role" "test" {
	name = "%s"
}

resource "snowflake_failover_group_grant" "g" {
	failover_group_name = snowflake_failover_group.fg.name
	privilege = "%s"
	roles = [
		snowflake_role.test.name
	]
}
`, name, allowedAccount, name, privilege)
}
