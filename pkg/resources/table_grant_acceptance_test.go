package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTableGrant_defaults(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
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
					testRolesAndShares(t, "snowflake_table_grant.g", []string{name}, []string{}),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_table_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
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

	roles = [
		snowflake_role.r.name
	]
}

`, n, n, n, n)
}
