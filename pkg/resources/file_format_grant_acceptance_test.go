package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_FileFormatGrant_defaults(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatGrantConfig(name, normal, "USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "file_format_name", name),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "privilege", "USAGE"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: fileFormatGrantConfig(name, normal, "ALL PRIVILEGES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", name),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatGrantConfig(name, onAll, "USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", name),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fileFormatGrantConfig(name, onFuture, "USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", name),
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

func fileFormatGrantConfig(name string, grantType grantType, privilege string) string {
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

resource snowflake_file_format test {
  name        = "%s"
  database    = snowflake_database.test.name
  schema      = snowflake_schema.test.name
  format_type = "PARQUET"

  compression = "AUTO"
}

resource snowflake_file_format_grant test {
    %s
	database_name = snowflake_database.test.name
	schema_name = snowflake_schema.test.name
	privilege = "%s"
	roles = [
		snowflake_role.test.name
	]
}

`, name, name, name, name, fileFormatNameConfig, privilege)
}
