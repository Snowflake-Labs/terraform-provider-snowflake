package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_MaskingPolicyGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_MASKING_POLICY_TESTS"); ok {
		t.Skip("Skipping TestAccMaskingPolicy")
	}
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyGrantConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "masking_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "privilege", "APPLY"),
				),
			},
		},
	})
}

func maskingPolicyGrantConfig(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_database" "test" {
		name = "%v"
		comment = "Terraform acceptance test"
	}

	resource "snowflake_schema" "test" {
		name = "%v"
		database = snowflake_database.test.name
		comment = "Terraform acceptance test"
	}

	resource "snowflake_role" "test" {
		name = "%v"
	}

	resource "snowflake_masking_policy" "test" {
		name = "%v"
		database = snowflake_database.test.name
		schema = snowflake_schema.test.name
		value_data_type = "VARCHAR"
		masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
		return_data_type = "VARCHAR(16777216)"
		comment = "Terraform acceptance test"
	}

	resource "snowflake_masking_policy_grant" "test" {
		masking_policy_name = snowflake_masking_policy.test.name
		database_name = snowflake_database.test.name
		roles         = [snowflake_role.test.name]
		schema_name   = snowflake_schema.test.name
		privilege = "APPLY"
	}
	`, name, name, name, name)
}
