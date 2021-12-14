package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_RowAccessPolicyGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_ROW_ACCESS_POLICY_TESTS"); ok {
		t.Skip("Skipping TestAccRowAccessPolicy")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: rowAccessPolicyGrantConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "row_access_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "privilege", "APPLY"),
				),
			},
		},
	})
}

func rowAccessPolicyGrantConfig(n string) string {
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
	privilege = "APPLY"
}
`, n, n, n, n)
}
