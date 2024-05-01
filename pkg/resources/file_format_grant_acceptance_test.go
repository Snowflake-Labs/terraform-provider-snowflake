package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FileFormatGrant_defaults(t *testing.T) {
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
				Config: fileFormatGrantConfig(name, normal, "USAGE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "file_format_name", name),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "privilege", "USAGE"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: fileFormatGrantConfig(name, normal, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "file_format_name", name),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_file_format_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_FileFormatGrant_onAll(t *testing.T) {
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
				Config: fileFormatGrantConfig(name, onAll, "USAGE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_file_format_grant.test", "file_format_name"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "privilege", "USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_file_format_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_FileFormatGrant_onFuture(t *testing.T) {
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
				Config: fileFormatGrantConfig(name, onFuture, "USAGE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_file_format_grant.test", "file_format_name"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "privilege", "USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_file_format_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func fileFormatGrantConfig(name string, grantType grantType, privilege string, databaseName string, schemaName string) string {
	var fileFormatNameConfig string
	switch grantType {
	case normal:
		fileFormatNameConfig = "file_format_name = snowflake_file_format.test.name"
	case onFuture:
		fileFormatNameConfig = "on_future = true"
	case onAll:
		fileFormatNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource snowflake_role test {
  name = "%s"
}

resource snowflake_file_format test {
  name        = "%s"
  database    = "%s"
  schema      = "%s"
  format_type = "PARQUET"

  compression = "AUTO"
}

resource snowflake_file_format_grant test {
    %s
	database_name = "%s"
	schema_name = "%s"
	privilege = "%s"
	roles = [
		snowflake_role.test.name
	]
}
`, name, name, databaseName, schemaName, fileFormatNameConfig, databaseName, schemaName, privilege)
}
