package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProcedureGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: procedureGrantConfig(name, onAll, "USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_procedure_grant.test", "procedure_name"),
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "privilege", "USAGE"),
				),
			},
			{
				ResourceName:      "snowflake_procedure_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAccProcedureGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: procedureGrantConfig(name, onFuture, "USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_procedure_grant.test", "procedure_name"),
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_procedure_grant.test", "privilege", "USAGE"),
				),
			},
			{
				ResourceName:      "snowflake_procedure_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func procedureGrantConfig(name string, grantType grantType, privilege string) string {
	var procedureNameConfig string
	switch grantType {
	case onFuture:
		procedureNameConfig = "on_future = true"
	case onAll:
		procedureNameConfig = "on_all = true"
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

resource snowflake_procedure_grant test {
    database_name = snowflake_database.test.name	
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	%s
	privilege = "%s"
}
`, name, name, name, procedureNameConfig, privilege)
}
