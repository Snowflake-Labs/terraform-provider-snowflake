package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTableGrant_onAll(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableGrantConfigOnAll(name),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "privilege", "SELECT"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "roles.0", name),

					testRolesAndShares(t, "snowflake_table_grant.g", []string{name}),
				),
			},
		},
	})
}

func TestAccTableGrant_defaults(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableGrantConfig(name),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.d", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.s", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.s", "database", name),
					resource.TestCheckResourceAttr("snowflake_role.r", "name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "table_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "privilege", "SELECT"),
					testRolesAndShares(t, "snowflake_table_grant.g", []string{name}),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_table_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
					"with_grant_option",
					"on_future",
					"on_all",
				},
			},
		},
	})
}

func tableGrantConfig(n string) string {
	return fmt.Sprintf(`

resource snowflake_database d {
	name = "%s"
}

resource snowflake_schema s {
	name = "%s"
	database = snowflake_database.d.name
}

resource snowflake_role r {
  name = "%s"
}

resource snowflake_table t {
	database = snowflake_database.d.name
	schema   = snowflake_schema.s.name
	name     = "%s"

	column {
		name = "id"
		type = "NUMBER(38,0)"
	}
}

resource snowflake_table_grant g {
	database_name = snowflake_database.d.name
	schema_name = snowflake_schema.s.name
	table_name = snowflake_table.t.name
	privilege = "SELECT"
	roles = [
		snowflake_role.r.name
	]
}

`, n, n, n, n)
}

func tableGrantConfigOnAll(n string) string {
	return fmt.Sprintf(`

resource snowflake_database d {
	name = "%s"
}

resource snowflake_schema s {
	name = "%s"
	database = snowflake_database.d.name
}

resource snowflake_role r {
  name = "%s"
}

resource snowflake_table t {
	database = snowflake_database.d.name
	schema   = snowflake_schema.s.name
	name     = "%s"

	column {
		name = "id"
		type = "NUMBER(38,0)"
	}
}

resource snowflake_table_grant g {
	database_name = snowflake_database.d.name
	schema_name = snowflake_schema.s.name

	roles = [
		snowflake_role.r.name
	]
	privilege = "SELECT"
	on_all = true
}

`, n, n, n, n)
}
