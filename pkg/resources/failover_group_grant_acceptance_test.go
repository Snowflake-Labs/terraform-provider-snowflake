package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_FailoverGroupGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"); !ok {
		t.Skip("Skipping TestAcc_FailoverGroup since not a business critical account")
	}
	accountName := os.Getenv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT")
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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
