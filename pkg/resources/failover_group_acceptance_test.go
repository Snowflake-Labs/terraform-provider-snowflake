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

func TestAcc_FailoverGroupBasic(t *testing.T) {
	randomCharacters := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	if _, ok := os.LookupEnv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"); !ok {
		t.Skip("Skipping TestAcc_FailoverGroup since not a business critical account")
	}
	accountName := os.Getenv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT")
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: failoverGroupBasic(randomCharacters, accountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", randomCharacters),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.expression", "0 0 10-20 * TUE,THU"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.time_zone", "UTC"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_failover_group.fg",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_edition_check"},
			},
		},
	})
}

func TestAcc_FailoverGroupRemoveObjectTypes(t *testing.T) {
	randomCharacters := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	if _, ok := os.LookupEnv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"); !ok {
		t.Skip("Skipping TestAcc_FailoverGroup since not a business critical account")
	}
	accountName := os.Getenv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT")
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithInterval(randomCharacters, accountName, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", randomCharacters),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "20"),
				),
			},
			{
				Config: failoverGroupWithNoWarehouse(randomCharacters, accountName, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", randomCharacters),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "3"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "20"),
				),
			},
		},
	})
}

func TestAcc_FailoverGroupInterval(t *testing.T) {
	randomCharacters := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	if _, ok := os.LookupEnv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"); !ok {
		t.Skip("Skipping TestAcc_FailoverGroup since not a business critical account")
	}
	accountName := os.Getenv("SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT")
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: failoverGroupWithInterval(randomCharacters, accountName, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", randomCharacters),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "10"),
				),
			},
			// Update Interval
			{
				Config: failoverGroupWithInterval(randomCharacters, accountName, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", randomCharacters),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "20"),
				),
			},
			// Change to Cron Expression
			{
				Config: failoverGroupWithCronExpression(randomCharacters, accountName, "0 0 10-20 * TUE,THU"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", randomCharacters),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.expression", "0 0 10-20 * TUE,THU"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.time_zone", "UTC"),
				),
			},
			// Update Cron Expression
			{
				Config: failoverGroupWithCronExpression(randomCharacters, accountName, "0 0 5-20 * TUE,THU"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", randomCharacters),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.expression", "0 0 5-20 * TUE,THU"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.0.time_zone", "UTC"),
				),
			},
			// Change to Interval
			{
				Config: failoverGroupWithInterval(randomCharacters, accountName, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "name", randomCharacters),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "object_types.#", "4"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_databases.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "allowed_integration_types.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.#", "1"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr("snowflake_failover_group.fg", "replication_schedule.0.interval", "10"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_failover_group.fg",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_edition_check"},
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
	name = "%s"
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
}
`, randomCharacters, randomCharacters, accountName)
}

func failoverGroupWithInterval(randomCharacters, accountName string, interval int) string {
	return fmt.Sprintf(`
resource "snowflake_database" "db" {
	name = "tst-terraform-%s"
}

resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_databases = [snowflake_database.db.name]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		interval = %d
	}
}
`, randomCharacters, randomCharacters, accountName, interval)
}

func failoverGroupWithNoWarehouse(randomCharacters, accountName string, interval int) string {
	return fmt.Sprintf(`
resource "snowflake_database" "db" {
	name = "tst-terraform-%s"
}

resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		interval = %d
	}
}
`, randomCharacters, randomCharacters, accountName, interval)
}

func failoverGroupWithCronExpression(randomCharacters, accountName, expression string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "db" {
	name = "tst-terraform-%s"
}

resource "snowflake_failover_group" "fg" {
	name = "%s"
	object_types = ["WAREHOUSES","DATABASES", "INTEGRATIONS", "ROLES"]
	allowed_accounts= ["%s"]
	allowed_databases = [snowflake_database.db.name]
	allowed_integration_types = ["SECURITY INTEGRATIONS"]
	replication_schedule {
		cron {
			expression = "%s"
			time_zone  = "UTC"
		}
	}
}
`, randomCharacters, randomCharacters, accountName, expression)
}
