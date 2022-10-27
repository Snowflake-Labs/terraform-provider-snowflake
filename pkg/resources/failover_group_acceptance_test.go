package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_FailoverGroup(t *testing.T) {
	randomCharacters := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	if _, ok := os.LookupEnv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"); !ok {
		t.Skip("Skipping TestAcc_FailoverGroup since not a business critical account")
	}
	accountName := os.Getenv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT")
	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: failoverGroupBasic(randomCharacters, accountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", fmt.Sprintf("tst-terraform-%s", randomCharacters)),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.expression", "0 0 10-20 * TUE,THU"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.time_zone", "UTC"),
				),
			},
		},
	})
}

func failoverGroupBasic(randomCharacters, accountName string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "db" {
	name = "tst-terraform-%s"
}

resource "snowflake_failover_group" "fg" {
	name = "tst-terraform-%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_databases = [snowflake_database.db.name]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
	cron {
		expression = "0 0 10-20 * TUE,THU"
		time_zone = "UTC"
	}
}
`, randomCharacters, randomCharacters, accountName)
}
