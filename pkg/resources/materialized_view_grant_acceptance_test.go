package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_MaterializedViewFutureGrant(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: materializedViewGrantConfigFuture(name, onFuture, "SELECT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_materialized_view_grant.test", "materialized_view_name"),
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "privilege", "SELECT"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_materialized_view_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_MaterializedViewAllGrant(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: materializedViewGrantConfigFuture(name, onAll, "SELECT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "schema_name", name),
					resource.TestCheckNoResourceAttr("snowflake_materialized_view_grant.test", "materialized_view_name"),
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_materialized_view_grant.test", "privilege", "SELECT"),
				),
			},
			{
				ResourceName:      "snowflake_materialized_view_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func materializedViewGrantConfigFuture(name string, grantType grantType, privilege string) string {
	var materializedViewNameConfig string
	switch grantType {
	case onFuture:
		materializedViewNameConfig = "on_future = true"
	case onAll:
		materializedViewNameConfig = "on_all = true"
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

resource "snowflake_materialized_view_grant" "test" {
    database_name = snowflake_database.test.name	
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	%s
	privilege = "%s"
}
`, name, name, name, materializedViewNameConfig, privilege)
}
