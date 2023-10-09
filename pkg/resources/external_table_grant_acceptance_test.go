package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ExternalTableGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalTableGrantConfig(name, onAll, "SELECT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_external_table_grant.test", "external_table_name"),
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "privilege", "SELECT"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_external_table_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_ExternalTableGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalTableGrantConfig(name, onFuture, "SELECT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_external_table_grant.test", "external_table_name"),
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_external_table_grant.test", "privilege", "SELECT"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_external_table_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func externalTableGrantConfig(name string, grantType grantType, privilege string) string {
	var externalTableNameConfig string
	switch grantType {
	case onFuture:
		externalTableNameConfig = "on_future = true"
	case onAll:
		externalTableNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%s"
}

resource "snowflake_schema" "test" {
	name = "%s"
	database = snowflake_database.test.name
}

resource "snowflake_role" "test" {
  name = "%s"
}

resource "snowflake_external_table_grant" "test" {
    database_name = snowflake_database.test.name	
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	%s
	privilege = "%s"
}
`, name, name, name, externalTableNameConfig, privilege)
}
