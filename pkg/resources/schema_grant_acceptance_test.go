package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_SchemaGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SHARE_TESTS"); ok {
		t.Skip("Skipping TestAccSchemaGrant")
	}

	sName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
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
			},
		},
	})
}

func TestAcc_SchemaFutureGrants(t *testing.T) {

	sName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleNameTable := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleNameView := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
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
			},
		},
	})
}

func futureTableAndViewGrantConfig(n, role_table, role_view string) string {
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

`, n, n, role_table, role_view)
}

func schemaGrantConfig(n, role, share string, future bool) string {
	schema_name_config := `schema_name   = snowflake_schema.test.name
  shares        = [snowflake_share.test.name]`

	if future {
		schema_name_config = "on_future     = true"
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
  accounts = ["PC37737"]
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
`, n, n, role, share, schema_name_config)
}
