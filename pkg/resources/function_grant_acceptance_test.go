package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFunctionGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: functionGrantConfig(name, onFuture),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "schema_name", name),
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

func TestAccFunctionGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: functionGrantConfig(name, onAll),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_function_grant.test", "function_name"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_function_grant.test", "on_all", "true"),
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

func functionGrantConfig(name string, grantType grantType) string {
	var functionNameConfig string
	switch grantType {
	case onFuture:
		functionNameConfig = "on_future = true"
	case onAll:
		functionNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource snowflake_database test {
  name = "%s"
}

resource snowflake_schema test {
	name = "%s"
	database = snowflake_database.test.name
}

resource snowflake_role test {
  name = "%s"
}

resource "snowflake_function_grant" "test" {
    database_name = snowflake_database.test.name
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	%s
	privilege = "USAGE"
}
`, name, name, name, functionNameConfig)
}
