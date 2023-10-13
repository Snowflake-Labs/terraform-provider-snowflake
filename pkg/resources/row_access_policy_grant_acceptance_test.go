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

func TestAcc_RowAccessPolicyGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_ROW_ACCESS_POLICY_TESTS"); ok {
		t.Skip("Skipping TestAccRowAccessPolicy")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: rowAccessPolicyGrantConfig(accName, "APPLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "row_access_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "privilege", "APPLY"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: rowAccessPolicyGrantConfig(accName, "ALL PRIVILEGES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "row_access_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			{
				ResourceName:      "snowflake_row_access_policy_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func rowAccessPolicyGrantConfig(n, privilege string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_role" "test" {
	name = "%v"
}

resource "snowflake_row_access_policy" "test" {
	name = "%v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	signature = {
		N = "VARCHAR"
		V = "VARCHAR",
	}
	row_access_expression = "case when current_role() in ('ANALYST') then true else false end"
	comment = "Terraform acceptance test"
}

resource "snowflake_row_access_policy_grant" "test" {
	row_access_policy_name = snowflake_row_access_policy.test.name
	database_name = snowflake_database.test.name
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	privilege = "%s"
}
`, n, n, n, n, privilege)
}
