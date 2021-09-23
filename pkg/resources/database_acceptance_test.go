package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestCheckReplication(t *testing.T, path string, replicas, failovers []string) func(*terraform.State) error {
	return func(state *terraform.State) error {
		is := state.RootModule().Resources[path].Primary

		if c, ok := is.Attributes["replication_accounts.#"]; !ok || MustParseInt(t, c) != int64(len(replicas)) {
			return fmt.Errorf("expected replication_accounts.# to equal %d but got %s", len(replicas), c)
		}
		r, err := extractList(is.Attributes, "replication_accounts")
		if err != nil {
			return err
		}

		if !listSetEqual(replicas, r) {
			return fmt.Errorf("expected replication_accounts %#v but got %#v", replicas, r)
		}

		if c, ok := is.Attributes["replication_failover_accounts.#"]; !ok || MustParseInt(t, c) != int64(len(failovers)) {
			return fmt.Errorf("expected replication_failover_accounts.# to equal %d but got %s", len(failovers), c)
		}
		r, err = extractList(is.Attributes, "replication_failover_accounts")
		if err != nil {
			return err
		}

		if !listSetEqual(replicas, r) {
			return fmt.Errorf("expected replication_failover_accounts %#v but got %#v", failovers, r)
		}

		return nil
	}
}

func TestAcc_Database(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_DATABASE_TESTS"); ok {
		t.Skip("Skipping TestAccDatabase")
	}

	prefix := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	acc1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	acc2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: dbConfig(prefix, acc1, acc2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_database.db", "data_retention_time_in_days"),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_is_primary", "true"),
					testCheckReplication(t, "snowflake_database.db", []string{acc1, acc2}, []string{acc1}),
				),
			},
			// RENAME
			{
				Config: dbConfig(prefix2, acc1, acc2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_database.db", "data_retention_time_in_days"),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_is_primary", "true"),
					testCheckReplication(t, "snowflake_database.db", []string{acc1, acc2}, []string{acc1}),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: dbConfig2(prefix2, acc2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment 2"),
					resource.TestCheckResourceAttr("snowflake_database.db", "data_retention_time_in_days", "3"),
					resource.TestCheckResourceAttr("snowflake_database.db", "replication_is_primary", "false"),
					testCheckReplication(t, "snowflake_database.db", []string{acc2}, []string{}),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_database.db",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func dbConfig(prefix, acc1, acc2 string) string {
	s := `
resource "snowflake_database" "db" {
	name = "%s"
	comment = "test comment"
        replication_accounts = ["%s", "%s"]
        replication_failover_accounts = ["%s"]
        replication_is_primary = true
}
`
	return fmt.Sprintf(s, prefix, acc1, acc2, acc1)
}

func dbConfig2(prefix, acc2 string) string {
	s := `
resource "snowflake_database" "db" {
	name = "%s"
	comment = "test comment 2"
	data_retention_time_in_days = 3
        replication_accounts = ["%s"]
        replication_failover_accounts = []
}
`
	return fmt.Sprintf(s, prefix, acc2)
}
