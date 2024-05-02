package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SchemaGrant(t *testing.T) {
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
				Config: schemaGrantConfig(name, normal, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_all", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "privilege", "USAGE"),
				),
			},
			// FUTURE SHARES
			{
				Config: schemaGrantConfig(name, onFuture, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("snowflake_schema_grant.test", "schema_name"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_all", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "privilege", "USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_schema_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_SchemaGrantOnAll(t *testing.T) {
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
				Config: schemaGrantConfig(name, onAll, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("snowflake_schema_grant.test", "schema_name"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "privilege", "USAGE"),
				),
			},
		},
	})
}

func schemaGrantConfig(name string, grantType grantType, databaseName string) string {
	var schemaNameConfig string
	switch grantType {
	case normal:
		schemaNameConfig = "schema_name = snowflake_schema.test.name"
	case onFuture:
		schemaNameConfig = "on_future = true"
	case onAll:
		schemaNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
  name      = "%v"
  database  = "%s"
  comment   = "Terraform acceptance test"
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_share" "test" {
  name     = "%v"
}

resource "snowflake_database_grant" "test" {
  database_name = "%s"
  shares        = [snowflake_share.test.name]
}

resource "snowflake_schema_grant" "test" {
  database_name = "%s"
  %v
  roles         = [snowflake_role.test.name]
}
`, name, databaseName, name, name, databaseName, databaseName, schemaNameConfig)
}
