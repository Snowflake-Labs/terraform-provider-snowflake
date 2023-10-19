package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_FunctionGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: functionGrantConfig(name, onFuture, "USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_function_grant.test", "function_name"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "privilege", "USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_function_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_FunctionGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: functionGrantConfig(name, onAll, "USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_function_grant.test", "function_name"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "privilege", "USAGE"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: functionGrantConfig(name, onAll, "ALL PRIVILEGES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_function_grant.test", "function_name"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_function_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func functionGrantConfig(name string, grantType grantType, privilege string) string {
	var functionNameConfig string
	switch grantType {
	case onFuture:
		functionNameConfig = "on_future = true"
	case onAll:
		functionNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource snowflake_role test {
  name = "%s"
}

resource "snowflake_function_grant" "test" {
    database_name = "terraform_test_database"
	roles         = [snowflake_role.test.name]
	schema_name   = "terraform_test_schema"
	%s
	privilege = "%s"
}
`, name, functionNameConfig, privilege)
}
