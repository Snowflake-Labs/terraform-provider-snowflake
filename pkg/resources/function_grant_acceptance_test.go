package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FunctionGrant_onFuture(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: functionGrantConfig(name, onFuture, "USAGE", acc.TestDatabaseName, acc.TestSchemaName),
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
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: functionGrantConfig(name, onAll, "USAGE", acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: functionGrantConfig(name, onAll, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
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

func functionGrantConfig(name string, grantType grantType, privilege, databaseName, schemaName string) string {
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
    database_name = "%s"
	roles         = [snowflake_role.test.name]
	schema_name   = "%s"
	%s
	privilege = "%s"
}
`, name, databaseName, schemaName, functionNameConfig, privilege)
}
