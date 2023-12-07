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
		t.Skip("Skipping TestAcc_RowAccessPolicyGrant")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: rowAccessPolicyGrantConfig(accName, "APPLY", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "row_access_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "privilege", "APPLY"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: rowAccessPolicyGrantConfig(accName, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy_grant.test", "schema_name", acc.TestSchemaName),
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

func rowAccessPolicyGrantConfig(n string, privilege string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_role" "test" {
	name = "%v"
}

resource "snowflake_row_access_policy" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	signature = {
		N = "VARCHAR"
		V = "VARCHAR",
	}
	row_access_expression = "case when current_role() in ('ANALYST') then true else false end"
	comment = "Terraform acceptance test"
}

resource "snowflake_row_access_policy_grant" "test" {
	row_access_policy_name = snowflake_row_access_policy.test.name
	database_name = "%s"
	roles         = [snowflake_role.test.name]
	schema_name   = "%s"
	privilege = "%s"
}
`, n, n, databaseName, schemaName, databaseName, schemaName, privilege)
}
