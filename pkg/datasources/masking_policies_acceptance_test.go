package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_MaskingPolicies(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	maskingPolicyName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicies(databaseName, schemaName, maskingPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_masking_policies.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_masking_policies.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_masking_policies.t", "masking_policies.#"),
					resource.TestCheckResourceAttr("data.snowflake_masking_policies.t", "masking_policies.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_masking_policies.t", "masking_policies.0.name", maskingPolicyName),
				),
			},
		},
	})
}

func maskingPolicies(databaseName string, schemaName string, maskingPolicyName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "test" {
		name = "%v"
	}

	resource snowflake_schema "test"{
		name 	 = "%v"
		database = snowflake_database.test.name
	}

	resource "snowflake_masking_policy" "test" {
		name 	 		   = "%v"
		database 	       = snowflake_database.test.name
		schema   		   = snowflake_schema.test.name
		signature {
			column {
				name = "val"
				type = "VARCHAR"
			}
		}
		masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
		return_data_type   = "VARCHAR(16777216)"
		comment            = "Terraform acceptance test"
	}

	data snowflake_masking_policies "t" {
		database = snowflake_masking_policy.test.database
		schema   = snowflake_masking_policy.test.schema
		depends_on = [snowflake_masking_policy.test]
	}
	`, databaseName, schemaName, maskingPolicyName)
}
