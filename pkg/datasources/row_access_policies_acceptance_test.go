package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RowAccessPolicies(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	rowAccessPolicyName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
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
