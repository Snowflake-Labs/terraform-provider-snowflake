package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRowAccessPolicies(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	rowAccessPolicyName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: rowAccessPolicies(databaseName, schemaName, rowAccessPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_row_access_policies.v", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_row_access_policies.v", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_row_access_policies.v", "row_access_policies.#"),
					resource.TestCheckResourceAttr("data.snowflake_row_access_policies.v", "row_access_policies.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_row_access_policies.v", "row_access_policies.0.name", rowAccessPolicyName),
				),
			},
		},
	})
}

func rowAccessPolicies(databaseName string, schemaName string, rowAccessPolicyName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "test" {
		name = "%v"
	}

	resource snowflake_schema "test"{
		name 	 = "%v"
		database = snowflake_database.test.name
	}

	resource "snowflake_row_access_policy" "test" {
		name = "%v"
		database = snowflake_database.test.name
		schema = snowflake_schema.test.name
		signature = {
			N = "VARCHAR"
			V = "VARCHAR",
		}
		row_access_expression = "case when current_role() in ('ANALYST') then true else false end"
		comment = "Terraform acceptance test"
	}

	data snowflake_row_access_policies "v" {
		database = snowflake_row_access_policy.test.database
		schema = snowflake_row_access_policy.test.schema
		depends_on = [snowflake_row_access_policy.test]
	}
	`, databaseName, schemaName, rowAccessPolicyName)
}
