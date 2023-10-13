package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_SchemaGrant(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: schemaGrantConfig(name, normal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_all", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "privilege", "USAGE"),
				),
			},
			// FUTURE SHARES
			{
				Config: schemaGrantConfig(name, onFuture),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: schemaGrantConfig(name, onAll),
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

func schemaGrantConfig(name string, grantType grantType) string {
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
resource "snowflake_database" "test" {
  name = "%v"
}

resource "snowflake_schema" "test" {
  name      = "%v"
  database  = snowflake_database.test.name
  comment   = "Terraform acceptance test"
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_share" "test" {
  name     = "%v"
}

resource "snowflake_database_grant" "test" {
  database_name = snowflake_schema.test.database
  shares        = [snowflake_share.test.name]
}

resource "snowflake_schema_grant" "test" {
  database_name = snowflake_schema.test.database
  %v
  roles         = [snowflake_role.test.name]
}
`, name, name, name, name, schemaNameConfig)
}
