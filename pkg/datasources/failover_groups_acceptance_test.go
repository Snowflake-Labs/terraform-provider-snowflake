package datasources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_FailoverGroups(t *testing.T) {
	if _, ok := os.LookupEnv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"); !ok {
		t.Skip("Skipping TestAcc_FailoverGroup since not a business critical account")
	}
	accountName := os.Getenv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT")

	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: failoverGroupsConfig(name, accountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_groups.d", "failover_groups.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_groups.d", "failover_groups.0.name", name),
					resource.TestCheckResourceAttr("snowflake_failover_groups.d", "failover_groups.0.object_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_groups.d", "failover_groups.0.object_types.0", "ROLES"),
					resource.TestCheckResourceAttr("snowflake_failover_groups.d", "failover_groups.0.allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_groups.d", "failover_groups.0.allowed_accounts.0", accountName),
				),
			},
		},
	})
}

func failoverGroupsConfig(failoverGroupName string, allowedAccount string) string {
	return fmt.Sprintf(`
	resource "snowflake_failover_group" "source_failover_group" {
		name                      = "%s"
		object_types              = ["ROLES"]
		allowed_accounts          = ["%s"]
	}

	data "snowflake_failover_groups" "d" {
		depends_on = [snowflake_failover_group.source_failover_group]
	}
	`, failoverGroupName, allowedAccount)
}
