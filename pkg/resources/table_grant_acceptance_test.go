package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTableGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableGrantConfig(name, onAll),

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
			{
				ResourceName:      "snowflake_table_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAccTableGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableGrantConfig(name, onFuture),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "privilege", "SELECT"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "roles.0", name),

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
				},
			},
		},
	})
}

func TestAccTableGrant_defaults(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableGrantConfig(name, normal),

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
				},
			},
		},
	})
}

func tableGrantConfig(name string, grantType grantType) string {
	var tableNameConfig string
	switch grantType {
	case normal:
		tableNameConfig = "table_name = snowflake_table.t.name"
	case onFuture:
		tableNameConfig = "on_future = true"
	case onAll:
		tableNameConfig = "on_all = true"
	}

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
    %s
	database_name = snowflake_database.d.name
	schema_name = snowflake_schema.s.name
	privilege = "SELECT"
	roles = [
		snowflake_role.r.name
	]
}

`, name, name, name, name, tableNameConfig)
}
