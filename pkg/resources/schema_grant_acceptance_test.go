package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_SchemaGrant(t *testing.T) {
	sName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: schemaGrantConfig(sName, roleName, shareName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "schema_name", sName),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "privilege", "USAGE"),
				),
			},
			// FUTURE SHARES
			{
				Config: schemaGrantConfig(sName, roleName, shareName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "schema_name", ""),
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

func TestAcc_SchemaFutureGrants(t *testing.T) {

	sName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleNameTable := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleNameView := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// TABLE AND VIEW FUTURE GRANTS
			{
				Config: futureTableAndViewGrantConfig(sName, roleNameTable, roleNameView),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.select_on_future_views", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table_grant.select_on_future_tables", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_view_grant.select_on_future_views", "privilege", "SELECT"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_view_grant.select_on_future_views",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func futureTableAndViewGrantConfig(n, roleTable, roleView string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%v"
}

resource "snowflake_schema" "test" {
  name      = "%v"
  database  = snowflake_database.test.name
  comment   = "Terraform acceptance test"
}

resource "snowflake_role" "table_reader" {
  name = "%v"
}

resource "snowflake_role" "view_reader" {
  name = "%v"
}

resource "snowflake_table_grant" "select_on_future_tables" {
  database_name = snowflake_database.test.name
  schema_name   = snowflake_schema.test.name
  privilege     = "SELECT"
  on_future     = true
  roles         = [snowflake_role.table_reader.name]
  depends_on    = [snowflake_schema.test, snowflake_role.table_reader]
}

resource "snowflake_view_grant" "select_on_future_views" {
  database_name = snowflake_database.test.name
  schema_name   = snowflake_schema.test.name
  privilege     = "SELECT"
  on_future     = true
  roles         = [snowflake_role.view_reader.name]
  depends_on    = [snowflake_schema.test, snowflake_role.view_reader]
}

`, n, n, roleTable, roleView)
}

func schemaGrantConfig(n, role, share string, future bool) string {
	schemaNameConfig := `schema_name   = snowflake_schema.test.name
  shares        = [snowflake_share.test.name]`

	if future {
		schemaNameConfig = "on_future     = true"
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

  depends_on = [snowflake_database_grant.test]
}
`, n, n, role, share, schemaNameConfig)
}
