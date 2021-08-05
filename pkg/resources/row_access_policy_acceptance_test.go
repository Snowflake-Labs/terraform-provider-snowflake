package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_RowAccessPolicy(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_ROW_ACCESS_POLICY_TESTS"); ok {
		t.Skip("Skipping TestAccRowAccessPolicy")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: rowAccessPolicyConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "row_access_expression", "case when current_role() in ('ANALYST') then true else false end"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.N", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.V", "VARCHAR"),
				),
			},
		},
	})
}

func rowAccessPolicyConfig(n string) string {
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
`, n, n, n)
}
