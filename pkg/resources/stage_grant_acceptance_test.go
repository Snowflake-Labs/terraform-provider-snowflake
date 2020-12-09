package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStageGrant_defaults(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: stageGrantConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.d", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.s", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.s", "database", name),
					resource.TestCheckResourceAttr("snowflake_role.r", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "schema_name", name),
					testRolesAndShares(t, "snowflake_stage_grant.g", []string{name}, []string{}),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_stage_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func stageGrantConfig(n string) string {
	return fmt.Sprintf(`
	resource snowflake_database d {
		name = "%v"
		comment = "Terraform acceptance test"
	}
	
	resource snowflake_schema s {
		name = "%v"
		database = snowflake_database.d.name
		comment = "Terraform acceptance test"
	}
	
	resource snowflake_stage s {
		name = "%v"
		database = snowflake_database.d.name
		schema = snowflake_schema.s.name
		comment = "Terraform acceptance test"
	}

	resource snowflake_role r {
		name = "%s"
	}

	resource snowflake_stage_grant g {
		database_name = snowflake_database.d.name
		schema_name = snowflake_schema.s.name
		stage_name = snowflake_stage.s.name

		privilege = "READ"

		roles = [
			snowflake_role.r.name
		]
	}
`, n, n, n, n)
}
